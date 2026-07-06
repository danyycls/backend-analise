package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/danyele/podp/internal/shared/clients/opencnpj"
	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/pncpbusca"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/repositorios"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/shared/utils"
)

const maxConcorrencia = 5

type ConsultaCNPJOrgaoPNCPUseCase struct {
	pncpClient     *pncp.PNCPClient
	opencnpjClient *opencnpj.OpenCNPJClient
	redis          *redis.RedisCache
	repo           repositorios.PNCPRepository
}

func NovoConsultaCNPJOrgaoPNCPUseCase(pncp *pncp.PNCPClient, opencnpj *opencnpj.OpenCNPJClient, redis *redis.RedisCache, repo repositorios.PNCPRepository) *ConsultaCNPJOrgaoPNCPUseCase {
	return &ConsultaCNPJOrgaoPNCPUseCase{
		pncpClient:     pncp,
		opencnpjClient: opencnpj,
		redis:          redis,
		repo:           repo,
	}
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) AnaliseMultiplos(ctx context.Context, req pncp.AnaliseOrgaoPNCPRequest, eventos chan<- pncp.EventoAnalise) []*pncp.AnaliseResultado {
	log := logger.New("PNCP: UseCase: AnaliseMultiplos")

	total := len(req.CNPJs)
	if total == 0 {
		eventos <- pncp.EventoAnalise{Type: "completed", Total: 0}
		return nil
	}

	sem := make(chan struct{}, maxConcorrencia)
	var wg sync.WaitGroup
	var mu sync.Mutex
	processed := 0
	success := 0
	errors := 0
	results := make([]*pncp.AnaliseResultado, 0, total)

	for _, cnpj := range req.CNPJs {
		wg.Add(1)
		sem <- struct{}{}

		go func(cnpj string) {
			defer wg.Done()
			defer func() { <-sem }()

			eventos <- pncp.EventoAnalise{Type: "started", CNPJ: cnpj}

			result, err := u.executar(ctx, cnpj, req.DataInicial, req.DataFinal)

			mu.Lock()
			processed++
			if err != nil {
				errors++
				log.Error("erro ao processar CNPJ", "cnpj", cnpj, "erro", err)
				eventos <- pncp.EventoAnalise{Type: "error", CNPJ: cnpj, Message: err.Error()}
			} else {
				success++
				results = append(results, result)
				eventos <- pncp.EventoAnalise{
					Type:  "success",
					CNPJ:  cnpj,
					Orgao: nomeOrgaoFromResult(result),
				}
			}
			eventos <- pncp.EventoAnalise{
				Type: "progress", Processed: processed, Total: total,
				Success: success, Errors: errors,
			}
			mu.Unlock()
		}(cnpj)
	}

	wg.Wait()
	log.Info("analise multipla concluida", "total", total, "success", success, "errors", errors)
	eventos <- pncp.EventoAnalise{Type: "completed", Total: total}
	return results
}

func nomeOrgaoFromResult(result *pncp.AnaliseResultado) string {
	if result != nil && result.Orgao != nil && result.Orgao.RazaoSocial != nil {
		return *result.Orgao.RazaoSocial
	}
	return ""
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) executar(ctx context.Context, cnpjOrgao, dataInicial, dataFinal string) (*pncp.AnaliseResultado, error) {
	log := logger.New("PNCP: UseCase: executar")
	cnpj := utils.NormalizarCNPJ(cnpjOrgao)

	if result := u.checkCache(ctx, cnpj, dataInicial, dataFinal); result != nil {
		return result, nil
	}

	log.Info("iniciando analise orgao", "cnpj", cnpj, "data_inicial", dataInicial, "data_final", dataFinal)

	contratos := u.coletarContratos(ctx, "orgao", cnpj, dataInicial, dataFinal)
	if len(contratos) == 0 {
		return nil, fmt.Errorf("erro ao consultar contratos para CNPJ %s", cnpj)
	}

	log.Info("contratos encontrados", "total", len(contratos))
	result := u.montarResultado(ctx, contratos, cnpj, dataInicial, dataFinal)
	u.writeCache(ctx, cnpj, dataInicial, dataFinal, result)

	return result, nil
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) checkCache(ctx context.Context, cnpj, dataInicial, dataFinal string) *pncp.AnaliseResultado {
	cacheParams := map[string]interface{}{
		"tipo": "orgao", "valor": cnpj,
		"dataInicial": dataInicial, "dataFinal": dataFinal,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache(redis.ChaveLicitacoesTrimestre, raw)

	var cached []*pncp.AnaliseResultado
	if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
		logger.New("PNCP: UseCase: executar").Warn("cache indisponivel", "erro", err)
	} else if ok && len(cached) > 0 {
		logger.New("PNCP: UseCase: executar").Info("cache hit", "cnpj", cnpj)
		return cached[0]
	}
	return nil
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) writeCache(ctx context.Context, cnpj, dataInicial, dataFinal string, result *pncp.AnaliseResultado) {
	cacheParams := map[string]interface{}{
		"tipo": "orgao", "valor": cnpj,
		"dataInicial": dataInicial, "dataFinal": dataFinal,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache(redis.ChaveLicitacoesTrimestre, raw)

	if err := u.redis.Set(ctx, cacheKey, []*pncp.AnaliseResultado{result}); err != nil {
		logger.New("PNCP: UseCase: executar").Warn("cache indisponivel", "erro", err)
	}
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) coletarContratos(ctx context.Context, tipo, valor, dataInicial, dataFinal string) []pncp.Contrato {
	meses := utils.ExtrairMeses(dataInicial, dataFinal)
	if len(meses) == 0 {
		return nil
	}

	var contratos []pncp.Contrato
	var mesesPendentes []utils.AnoMes

	for _, am := range meses {
		jaRealizada, err := u.repo.BuscaJaRealizada(ctx, tipo, valor, am.Ano, am.Mes)
		if err != nil {
			logger.New("PNCP: UseCase: coletarContratos").Warn("erro ao verificar busca no PG", "ano", am.Ano, "mes", am.Mes, "erro", err)
			mesesPendentes = append(mesesPendentes, am)
			continue
		}
		if jaRealizada {
			persistidos, err := u.repo.BuscarContratosPorFiltro(ctx, tipo, valor, am.Ano, am.Mes)
			if err != nil {
				logger.New("PNCP: UseCase: coletarContratos").Warn("erro ao buscar contratos do PG", "ano", am.Ano, "mes", am.Mes, "erro", err)
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

	if len(mesesPendentes) > 0 {
		pendentes := buscarMesesParalelo(ctx, tipo, valor, mesesPendentes, u.buscarMesComLock)
		contratos = append(contratos, pendentes...)
	}

	return contratos
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) buscarMesComLock(ctx context.Context, tipo, valor string, ano, mes int) []pncp.Contrato {
	return pncpbusca.BuscarMesComLock(ctx, u.redis, u.repo,
		u.buscarContratosDoPNCP,
		func(c context.Context, t string, v string, a, m int, contratos []pncp.Contrato) error {
			return persistirContratos(c, u.repo, t, v, a, m, contratos)
		},
		tipo, valor, ano, mes,
	)
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) buscarContratosDoPNCP(ctx context.Context, _ string, valor string, ano, mes int) ([]pncp.Contrato, error) {
	log := logger.New("PNCP: UseCase: buscarContratosDoPNCP")

	dataInicial, dataFinal := utils.FormatarPeriodoMes(ano, mes)
	pagina := 1
	tamanho := 500
	contratos := make([]pncp.Contrato, 0)
	seen := make(map[string]struct{})

	for {
		resp, err := u.pncpClient.BuscarContratos(ctx, valor, dataInicial, dataFinal, pagina, tamanho)
		if err != nil {
			return nil, fmt.Errorf("pagina %d: %w", pagina, err)
		}
		if len(resp) == 0 {
			break
		}

		for _, c := range resp {
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

		if len(resp) < tamanho {
			break
		}
		pagina++

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	log.Info("contratos obtidos do PNCP", "valor", valor, "ano", ano, "mes", mes, "total", len(contratos))
	return contratos, nil
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) montarResultado(ctx context.Context, contratos []pncp.Contrato, cnpj, dataInicial, dataFinal string) *pncp.AnaliseResultado {
	log := logger.New("PNCP: UseCase: montarResultado")

	fornecedoresMap := extrairFornecedores(contratos)
	log.Info("fornecedores unicos extraidos", "total", len(fornecedoresMap))

	enrichedFornecedores := enriquecerFornecedoresComRepo(ctx, u.repo, u.opencnpjClient, fornecedoresMap)
	log.Info("empresas enriquecidas", "total", len(enrichedFornecedores))

	contratosDTO := montarContratosDTO(contratos, enrichedFornecedores)
	totalContratos, totalEmpresas, valorTotal := u.calcularResumo(contratosDTO, enrichedFornecedores)
	orgao := u.montarOrgaoInfo(cnpj, contratosDTO)

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

func (u *ConsultaCNPJOrgaoPNCPUseCase) calcularResumo(contratosDTO []pncp.Contrato, enrichedFornecedores map[string]*types.FornecedorOpenCNPJ) (int, int, float64) {
	totalContratos := len(contratosDTO)
	totalEmpresas := len(enrichedFornecedores)
	var valorTotal float64
	for _, c := range contratosDTO {
		if c.ValorGlobal != nil {
			valorTotal += *c.ValorGlobal
		}
	}
	return totalContratos, totalEmpresas, valorTotal
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) montarOrgaoInfo(cnpj string, contratosDTO []pncp.Contrato) *pncp.OrgaoInfo {
	orgao := &pncp.OrgaoInfo{CNPJ: pncp.StrPtr(cnpj)}
	if len(contratosDTO) > 0 && contratosDTO[0].OrgaoEntidade != nil {
		orgao = &pncp.OrgaoInfo{
			CNPJ:        contratosDTO[0].OrgaoEntidade.CNPJ,
			RazaoSocial: contratosDTO[0].OrgaoEntidade.RazaoSocial,
		}
	}
	return orgao
}
