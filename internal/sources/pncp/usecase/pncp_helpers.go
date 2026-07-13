package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/repositorios"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/shared/utils"
	"github.com/danyele/podp/internal/sources/opencnpj/client"
	"github.com/danyele/podp/internal/sources/pncp/client"
)

var maxConcorrencia int

func init() {
	maxConcorrencia = getEnvInt("PNCP_MAX_CONCORRENCIA", 1)
}

func getEnvInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}

type PncpUseCaseBase struct {
	pncpClient     *pncp.PNCPClient
	opencnpjClient *opencnpj.OpenCNPJClient
	redis          *redis.RedisCache
	repo           repositorios.PNCPRepository
}

func NewPncpUseCaseBase(pncp *pncp.PNCPClient, opencnpj *opencnpj.OpenCNPJClient, redis *redis.RedisCache, repo repositorios.PNCPRepository) PncpUseCaseBase {
	return PncpUseCaseBase{
		pncpClient:     pncp,
		opencnpjClient: opencnpj,
		redis:          redis,
		repo:           repo,
	}
}

type FetchPaginaContratos func(ctx context.Context, valor, dataInicial, dataFinal string, pagina, tamanho int) (*pncp.ContratoResponse, error)

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
			defer time.Sleep(3 * time.Second)

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

func checkCache(ctx context.Context, redisCli *redis.RedisCache, tipo, valor, dataInicial, dataFinal string) []*pncp.AnaliseResultado {
	cacheParams := map[string]interface{}{
		"tipo": tipo, "valor": valor,
		"dataInicial": dataInicial, "dataFinal": dataFinal,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache(redis.ChaveLicitacoesTrimestre, raw)

	var cached []*pncp.AnaliseResultado
	if ok, err := redisCli.Get(ctx, cacheKey, &cached); err != nil {
		logger.New("PNCP: UseCase").Warn("cache indisponivel", "erro", err)
	} else if ok && len(cached) > 0 {
		logger.New("PNCP: UseCase").Info("cache hit", "tipo", tipo, "valor", valor)
		return cached
	}
	return nil
}

func writeCache(ctx context.Context, redisCli *redis.RedisCache, tipo, valor, dataInicial, dataFinal string, resultados []*pncp.AnaliseResultado) {
	cacheParams := map[string]interface{}{
		"tipo": tipo, "valor": valor,
		"dataInicial": dataInicial, "dataFinal": dataFinal,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache(redis.ChaveLicitacoesTrimestre, raw)

	if err := redisCli.Set(ctx, cacheKey, resultados); err != nil {
		logger.New("PNCP: UseCase").Warn("cache indisponivel", "erro", err)
	}
}

func consultarContratosExistentes(ctx context.Context, repo repositorios.PNCPRepository, tipo, valor string, meses []utils.AnoMes) ([]pncp.Contrato, []utils.AnoMes) {
	var contratos []pncp.Contrato
	var mesesPendentes []utils.AnoMes

	for _, am := range meses {
		jaRealizada, err := repo.BuscaJaRealizada(ctx, tipo, valor, am.Ano, am.Mes)
		if err != nil {
			logger.New("PNCP: UseCase: consultarContratosExistentes").Warn("erro ao verificar busca no PG", "ano", am.Ano, "mes", am.Mes, "erro", err)
			mesesPendentes = append(mesesPendentes, am)
			continue
		}
		if jaRealizada {
			persistidos, err := repo.BuscarContratosPorFiltro(ctx, tipo, valor, am.Ano, am.Mes)
			if err != nil {
				logger.New("PNCP: UseCase: consultarContratosExistentes").Warn("erro ao buscar contratos do PG", "ano", am.Ano, "mes", am.Mes, "erro", err)
				mesesPendentes = append(mesesPendentes, am)
				continue
			}
			for i := range persistidos {
				contratos = append(contratos, repositorios.PersistidoParaContrato(persistidos[i]))
			}
			continue
		}
		mesesPendentes = append(mesesPendentes, am)
	}

	return contratos, mesesPendentes
}

func buscarContratosPaginado(
	ctx context.Context,
	fetchPagina FetchPaginaContratos,
	valor string, ano, mes int,
	tamanhoPagina int,
	delay time.Duration,
) ([]pncp.Contrato, error) {
	log := logger.New("PNCP: UseCase: buscarContratosPaginado")

	dataInicial, dataFinal := utils.FormatarPeriodoMes(ano, mes)
	pagina := 1
	var contratos []pncp.Contrato
	seen := make(map[string]struct{})

	for {
		resp, err := fetchPagina(ctx, valor, dataInicial, dataFinal, pagina, tamanhoPagina)
		if err != nil {
			return nil, fmt.Errorf("pagina %d: %w", pagina, err)
		}
		if resp == nil || resp.Empty || len(resp.Data) == 0 {
			break
		}

		for _, c := range resp.Data {
			key := ""
			if c.NumeroControlePNCP != nil {
				key = *c.NumeroControlePNCP
			}
			if key == "" {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			contratos = append(contratos, c)
		}

		if pagina >= resp.TotalPaginas {
			break
		}
		pagina++

		if delay > 0 {
			time.Sleep(delay)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	log.Info("contratos obtidos do PNCP", "valor", valor, "ano", ano, "mes", mes, "total", len(contratos))
	return contratos, nil
}

func montarResultadoUnico(
	ctx context.Context,
	repo repositorios.PNCPRepository,
	opencnpjClient *opencnpj.OpenCNPJClient,
	contratos []pncp.Contrato,
	cnpjOrgao, dataInicial, dataFinal string,
) *pncp.AnaliseResultado {
	log := logger.New("PNCP: UseCase: montarResultadoUnico")

	fornecedoresMap := extrairFornecedores(contratos)
	log.Info("fornecedores unicos extraidos", "total", len(fornecedoresMap))

	enrichedFornecedores := enriquecerFornecedoresComRepo(ctx, repo, opencnpjClient, fornecedoresMap)
	log.Info("empresas enriquecidas", "total", len(enrichedFornecedores))

	contratosDTO := montarContratosDTO(contratos, enrichedFornecedores)

	var orgao *pncp.OrgaoInfo
	if len(contratosDTO) > 0 && contratosDTO[0].OrgaoEntidade != nil {
		orgao = &pncp.OrgaoInfo{
			CNPJ:        contratosDTO[0].OrgaoEntidade.CNPJ,
			RazaoSocial: contratosDTO[0].OrgaoEntidade.RazaoSocial,
		}
	} else {
		orgao = &pncp.OrgaoInfo{CNPJ: pncp.StrPtr(cnpjOrgao)}
	}

	totalContratos := len(contratosDTO)
	totalEmpresas := len(enrichedFornecedores)
	var valorTotal float64
	for _, c := range contratosDTO {
		if c.ValorGlobal != nil {
			valorTotal += *c.ValorGlobal
		}
	}

	return &pncp.AnaliseResultado{
		Orgao: orgao,
		Periodo: &pncp.Periodo{
			DataInicial: pncp.StrPtr(dataInicial),
			DataFinal:   pncp.StrPtr(dataFinal),
		},
		Resumo: &pncp.Resumo{
			TotalContratos:      &totalContratos,
			TotalEmpresas:       &totalEmpresas,
			ValorTotalContratos: &valorTotal,
		},
		Contratos: contratosDTO,
	}
}

func montarResultadosAgrupados(
	ctx context.Context,
	repo repositorios.PNCPRepository,
	opencnpjClient *opencnpj.OpenCNPJClient,
	contratos []pncp.Contrato,
	tipo, valor, dataInicial, dataFinal string,
) ([]*pncp.AnaliseResultado, error) {
	grupos := agruparContratosPorOrgao(contratos)
	if len(grupos) == 0 {
		return nil, fmt.Errorf("nenhum contrato encontrado para %s %s", tipo, valor)
	}

	resultados := make([]*pncp.AnaliseResultado, 0, len(grupos))
	for _, g := range grupos {
		cnpj := ""
		if g.CNPJ != nil {
			cnpj = *g.CNPJ
		}
		r := montarResultadoUnico(ctx, repo, opencnpjClient, g.Contratos, cnpj, dataInicial, dataFinal)
		resultados = append(resultados, r)
	}
	return resultados, nil
}

type grupoOrgao struct {
	CNPJ        *string
	RazaoSocial *string
	Contratos   []pncp.Contrato
}

func agruparContratosPorOrgao(contratos []pncp.Contrato) []*grupoOrgao {
	grupos := make(map[string]*grupoOrgao)
	ordem := make([]string, 0)

	for _, c := range contratos {
		cnpj := ""
		razao := ""
		if c.OrgaoEntidade != nil && c.OrgaoEntidade.CNPJ != nil {
			cnpj = *c.OrgaoEntidade.CNPJ
		}
		if c.OrgaoEntidade != nil && c.OrgaoEntidade.RazaoSocial != nil {
			razao = *c.OrgaoEntidade.RazaoSocial
		}

		if cnpj == "" {
			cnpj = "sem_cnpj"
		}

		if _, ok := grupos[cnpj]; !ok {
			grupos[cnpj] = &grupoOrgao{
				CNPJ:        pncp.StrPtr(cnpj),
				RazaoSocial: pncp.StrPtr(razao),
				Contratos:   make([]pncp.Contrato, 0),
			}
			ordem = append(ordem, cnpj)
		}
		grupos[cnpj].Contratos = append(grupos[cnpj].Contratos, c)
	}

	resultado := make([]*grupoOrgao, 0, len(grupos))
	for _, key := range ordem {
		resultado = append(resultado, grupos[key])
	}
	return resultado
}
