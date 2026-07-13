package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/repositorios"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/sources/pncp/client"
)

func filtrarContratosPorMes(contratos []pncp.Contrato, ano, mes int) []pncp.Contrato {
	prefixo := fmt.Sprintf("%d-%02d", ano, mes)
	var filtrados []pncp.Contrato
	for _, c := range contratos {
		if c.DataPublicacao != nil && strings.HasPrefix(*c.DataPublicacao, prefixo) {
			filtrados = append(filtrados, c)
		}
	}
	return filtrados
}

func persistirContratos(
	ctx context.Context,
	repo repositorios.PNCPRepository,
	tipo, valor string,
	ano, mes int,
	contratos []pncp.Contrato,
) error {
	log := logger.New("PNCP: UseCase: persistirContratos")

	if len(contratos) == 0 {
		log.Info("nenhum contrato para persistir", "tipo", tipo, "valor", valor, "ano", ano, "mes", mes)
		return nil
	}

	persistidos := make([]repositorios.ContratoPersistido, len(contratos))
	amparosMap := make(map[int]*pncp.AmparoLegal)
	fornecedoresMap := make(map[string]*types.FornecedorOpenCNPJ)

	for i, c := range contratos {
		cp := repositorios.ContratoParaPersistido(c)
		persistidos[i] = cp

		if c.AmparoLegal != nil && c.AmparoLegal.Codigo != nil {
			if _, ok := amparosMap[*c.AmparoLegal.Codigo]; !ok {
				amparosMap[*c.AmparoLegal.Codigo] = c.AmparoLegal
			}
		}

		cnpjF := ""
		if c.Fornecedor != nil && c.Fornecedor.CNPJ != nil {
			cnpjF = *c.Fornecedor.CNPJ
		} else if c.NIFornecedor != nil {
			cnpjF = *c.NIFornecedor
		}
		if cnpjF != "" {
			if _, ok := fornecedoresMap[cnpjF]; !ok && c.Fornecedor != nil {
				fornecedoresMap[cnpjF] = c.Fornecedor
			}
		}
	}

	return repo.ComTransaction(ctx, func(txRepo repositorios.PNCPRepository) error {
		if len(amparosMap) > 0 {
			for _, a := range amparosMap {
				ap := repositorios.ExtrairAmparoLegal(a)
				if ap != nil {
					if err := txRepo.SalvarAmparosLegais(ctx, []repositorios.AmparoLegalPersistido{*ap}); err != nil {
						log.Warn("erro ao salvar amparo legal", "codigo", a.Codigo, "erro", err)
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
						socioID, err := txRepo.SalvarSocio(ctx, sp)
						if err != nil {
							log.Warn("erro ao salvar socio", "cnpj", cnpjF, "erro", err)
							continue
						}
						vs := repositorios.SocioParaFornecedorSocio(cnpjF, socioID, s)
						socios = append(socios, vs)
					}
				}
			}
			if len(fornecedores) > 0 {
				if err := txRepo.SalvarFornecedores(ctx, fornecedores); err != nil {
					return err
				}
			}
			if len(socios) > 0 {
				if err := txRepo.SalvarFornecedorSocios(ctx, socios); err != nil {
					return err
				}
			}
		}

		if err := txRepo.SalvarContratos(ctx, persistidos); err != nil {
			return err
		}

		controle := repositorios.BuscaControlePersistido{
			TipoBusca:                 tipo,
			ValorBusca:                valor,
			Ano:                       ano,
			Mes:                       mes,
			DataInicial:               time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC),
			DataFinal:                 time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1),
			TotalContratosEncontrados: len(contratos),
		}
		return txRepo.RegistrarBusca(ctx, controle)
	})
}
