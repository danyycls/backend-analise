package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/danyele/podp/internal/shared/clients/opencnpj"
	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/repositorios"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/shared/utils"
)

func extrairFornecedores(contratos []pncp.Contrato) map[string]string {
	fornecedores := make(map[string]string)
	for _, c := range contratos {
		if c.NIFornecedor == nil || *c.NIFornecedor == "" {
			continue
		}
		ni := utils.NormalizarCNPJ(*c.NIFornecedor)
		if _, ok := fornecedores[ni]; !ok {
			nome := ""
			if c.NomeRazaoSocialFornecedor != nil {
				nome = *c.NomeRazaoSocialFornecedor
			}
			fornecedores[ni] = nome
		}
	}
	return fornecedores
}

func montarContratosDTO(contratos []pncp.Contrato, enrichedFornecedores map[string]*types.FornecedorOpenCNPJ) []pncp.Contrato {
	contratosDTO := make([]pncp.Contrato, len(contratos))
	copy(contratosDTO, contratos)
	for i, c := range contratosDTO {
		if c.NIFornecedor == nil {
			continue
		}
		ni := utils.NormalizarCNPJ(*c.NIFornecedor)
		if enriched, ok := enrichedFornecedores[ni]; ok {
			contratosDTO[i].Fornecedor = enriched
		}
	}
	return contratosDTO
}

func enriquecerFornecedoresComRepo(ctx context.Context, repo repositorios.PNCPRepository, opencnpjClient *opencnpj.OpenCNPJClient, fornecedoresMap map[string]string) map[string]*types.FornecedorOpenCNPJ {
	log := logger.New("PNCP: UseCase: enriquecerFornecedoresComRepo")
	enriched := make(map[string]*types.FornecedorOpenCNPJ, len(fornecedoresMap))

	for cnpjF, nome := range fornecedoresMap {
		fp, err := repo.BuscarFornecedor(ctx, cnpjF)
		if err == nil && fp != nil {
			enriched[cnpjF] = repositorios.PersistidoParaFornecedor(*fp)
			continue
		}

		data, err := opencnpjClient.Buscar(ctx, cnpjF)
		if err != nil {
			log.Warn("erro ao consultar OpenCNPJ", "cnpj", cnpjF, "erro", err)
			enriched[cnpjF] = &types.FornecedorOpenCNPJ{CNPJ: pncp.StrPtr(cnpjF), RazaoSocial: pncp.StrPtr(nome)}
			continue
		}

		dto := utils.BuildFornecedorDTO(data)
		enriched[cnpjF] = dto
		cp := repositorios.FornecedorParaPersistido(*dto)
		if err := repo.SalvarFornecedores(ctx, []repositorios.FornecedorPersistido{cp}); err != nil {
			log.Warn("erro ao persistir fornecedor", "cnpj", cnpjF, "erro", err)
		}

		if dto.Socios != nil {
			socios := make([]repositorios.FornecedorSocioPersistido, 0, len(dto.Socios))
			for _, s := range dto.Socios {
				sp := repositorios.SocioParaPersistido(s)
				socioID, err := repo.SalvarSocio(ctx, sp)
				if err != nil {
					continue
				}
				vs := repositorios.SocioParaFornecedorSocio(cnpjF, socioID, s)
				socios = append(socios, vs)
			}
			if len(socios) > 0 {
				_ = repo.SalvarFornecedorSocios(ctx, socios)
			}
		}
	}

	return enriched
}

func persistirContratos(ctx context.Context, repo repositorios.PNCPRepository, tipo, valor string, ano, mes int, contratos []pncp.Contrato) error {
	if len(contratos) == 0 {
		controle := repositorios.BuscaControlePersistido{
			TipoBusca:                tipo,
			ValorBusca:               valor,
			Ano:                      ano,
			Mes:                      mes,
			DataInicial:              time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC),
			DataFinal:                time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1),
			TotalContratosEncontrados: 0,
		}
		return repo.RegistrarBusca(ctx, controle)
	}

	persistidos := make([]repositorios.ContratoPersistido, len(contratos))
	amparosMap := make(map[int]bool)
	fornecedoresMap := make(map[string]*types.FornecedorOpenCNPJ)

	for i, c := range contratos {
		cp := repositorios.ContratoParaPersistido(c)
		persistidos[i] = cp

		if c.AmparoLegal != nil && c.AmparoLegal.Codigo != nil {
			amparosMap[*c.AmparoLegal.Codigo] = true
		}
		if c.Fornecedor != nil && c.Fornecedor.CNPJ != nil {
			cnpj := *c.Fornecedor.CNPJ
			if _, ok := fornecedoresMap[cnpj]; !ok {
				fornecedoresMap[cnpj] = c.Fornecedor
			}
		}
	}

	if len(amparosMap) > 0 {
		for codigo := range amparosMap {
			for _, c := range contratos {
				if c.AmparoLegal != nil && c.AmparoLegal.Codigo != nil && *c.AmparoLegal.Codigo == codigo {
					ap := repositorios.ExtrairAmparoLegal(c.AmparoLegal)
					if ap != nil {
						_ = repo.SalvarAmparosLegais(ctx, []repositorios.AmparoLegalPersistido{*ap})
					}
					break
				}
			}
		}
	}

	if len(fornecedoresMap) > 0 {
		fornecedores := make([]repositorios.FornecedorPersistido, 0, len(fornecedoresMap))
		socios := make([]repositorios.FornecedorSocioPersistido, 0)
		for cnpjF, f := range fornecedoresMap {
			fp := repositorios.FornecedorParaPersistido(*f)
			fornecedores = append(fornecedores, fp)

			if f.Socios != nil {
				for _, s := range f.Socios {
					sp := repositorios.SocioParaPersistido(s)
					socioID, err := repo.SalvarSocio(ctx, sp)
					if err != nil {
						continue
					}
					vs := repositorios.SocioParaFornecedorSocio(cnpjF, socioID, s)
					socios = append(socios, vs)
				}
			}
		}
		if len(fornecedores) > 0 {
			_ = repo.SalvarFornecedores(ctx, fornecedores)
		}
		if len(socios) > 0 {
			_ = repo.SalvarFornecedorSocios(ctx, socios)
		}
	}

	if err := repo.SalvarContratos(ctx, persistidos); err != nil {
		return err
	}

	controle := repositorios.BuscaControlePersistido{
		TipoBusca:                tipo,
		ValorBusca:               valor,
		Ano:                      ano,
		Mes:                      mes,
		DataInicial:              time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC),
		DataFinal:                time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1),
		TotalContratosEncontrados: len(contratos),
	}
	return repo.RegistrarBusca(ctx, controle)
}

func buscarMesesParalelo(
	ctx context.Context,
	tipo, valor string,
	meses []utils.AnoMes,
	buscarMes func(context.Context, string, string, int, int) []pncp.Contrato,
) []pncp.Contrato {
	total := make([]pncp.Contrato, 0)
	sem := make(chan struct{}, maxConcorrencia)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, am := range meses {
		wg.Add(1)
		sem <- struct{}{}

		go func(am utils.AnoMes) {
			defer wg.Done()
			defer func() { <-sem }()

			contratosMes := buscarMes(ctx, tipo, valor, am.Ano, am.Mes)
			if len(contratosMes) == 0 {
				return
			}

			mu.Lock()
			total = append(total, contratosMes...)
			mu.Unlock()
		}(am)
	}

	wg.Wait()
	return total
}
