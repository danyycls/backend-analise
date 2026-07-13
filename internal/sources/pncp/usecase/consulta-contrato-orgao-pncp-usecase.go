package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/repositorios"
	"github.com/danyele/podp/internal/shared/utils"
	"github.com/danyele/podp/internal/sources/opencnpj/client"
	"github.com/danyele/podp/internal/sources/pncp/client"
)

type ConsultaContratoOrgaoPNCPUseCase struct {
	PncpUseCaseBase
}

func NovoConsultaContratoOrgaoPNCPUseCase(pncp *pncp.PNCPClient, opencnpj *opencnpj.OpenCNPJClient, redis *redis.RedisCache, repo repositorios.PNCPRepository) *ConsultaContratoOrgaoPNCPUseCase {
	return &ConsultaContratoOrgaoPNCPUseCase{
		PncpUseCaseBase: NewPncpUseCaseBase(pncp, opencnpj, redis, repo),
	}
}

func (u *ConsultaContratoOrgaoPNCPUseCase) Executar(ctx context.Context, req pncp.AnaliseContratoOrgaoRequest, eventos chan<- pncp.EventoAnalise) []*pncp.AnaliseResultado {
	log := logger.New("PNCP: UseCase: ConsultaContratoOrgao")

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

			result, err := u.buscarContratosPorOrgao(ctx, cnpj, req.DataInicial, req.DataFinal)

			mu.Lock()
			processed++
			if err != nil {
				errors++
				log.Error("erro ao processar CNPJ", "cnpj", cnpj, "erro", err)
				eventos <- pncp.EventoAnalise{Type: "error", CNPJ: cnpj, Message: err.Error()}
			} else {
				success++
				results = append(results, result)
				orgaoNome := ""
				if result != nil && result.Orgao != nil && result.Orgao.RazaoSocial != nil {
					orgaoNome = *result.Orgao.RazaoSocial
				}
				eventos <- pncp.EventoAnalise{
					Type:  "success",
					CNPJ:  cnpj,
					Orgao: orgaoNome,
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
	log.Info("analise orgao concluida", "total", total, "success", success, "errors", errors)
	eventos <- pncp.EventoAnalise{Type: "completed", Total: total}
	return results
}

func (u *ConsultaContratoOrgaoPNCPUseCase) buscarContratosPorOrgao(ctx context.Context, cnpjOrgao, dataInicial, dataFinal string) (*pncp.AnaliseResultado, error) {
	log := logger.New("PNCP: UseCase: buscarContratosPorOrgao")
	cnpj := utils.NormalizarCNPJ(cnpjOrgao)

	if result := checkCache(ctx, u.redis, "orgao", cnpj, dataInicial, dataFinal); result != nil {
		return result[0], nil
	}

	log.Info("iniciando analise orgao", "cnpj", cnpj, "data_inicial", dataInicial, "data_final", dataFinal)

	meses := utils.ExtrairMeses(dataInicial, dataFinal)
	if len(meses) == 0 {
		return nil, fmt.Errorf("periodo invalido: %s a %s", dataInicial, dataFinal)
	}

	contratos, mesesPendentes := consultarContratosExistentes(ctx, u.repo, "orgao", cnpj, meses)
	if len(contratos) == 0 && len(mesesPendentes) == 0 {
		return nil, fmt.Errorf("erro ao consultar contratos para CNPJ %s", cnpj)
	}

	if len(mesesPendentes) > 0 {
		pendentes := buscarMesesParalelo(ctx, "orgao", cnpj, mesesPendentes, u.fetchContratosPorMes)
		contratos = append(contratos, pendentes...)
	}

	log.Info("contratos encontrados", "total", len(contratos), "meses_pendentes", len(mesesPendentes))

	result := montarResultadoUnico(ctx, u.repo, u.opencnpjClient, contratos, cnpj, dataInicial, dataFinal)

	for _, am := range mesesPendentes {
		contratosDoMes := filtrarContratosPorMes(result.Contratos, am.Ano, am.Mes)
		if len(contratosDoMes) > 0 {
			if err := persistirContratos(ctx, u.repo, "orgao", cnpj, am.Ano, am.Mes, contratosDoMes); err != nil {
				log.Warn("erro ao persistir contratos", "ano", am.Ano, "mes", am.Mes, "erro", err)
			}
		}
	}

	cacheResult := []*pncp.AnaliseResultado{result}
	writeCache(ctx, u.redis, "orgao", cnpj, dataInicial, dataFinal, cacheResult)

	return result, nil
}

func (u *ConsultaContratoOrgaoPNCPUseCase) fetchContratosPorMes(ctx context.Context, _ string, valor string, ano, mes int) []pncp.Contrato {
	log := logger.New("PNCP: UseCase: fetchContratosPorMes")

	fetchPagina := func(ctx context.Context, valor, dataInicial, dataFinal string, pagina, tamanho int) (*pncp.ContratoResponse, error) {
		return u.pncpClient.BuscarContratos(ctx, valor, dataInicial, dataFinal, pagina, tamanho)
	}

	contratos, err := buscarContratosPaginado(ctx, fetchPagina, valor, ano, mes, 200, 0)
	if err != nil {
		log.Error("erro ao buscar contratos do PNCP", "valor", valor, "ano", ano, "mes", mes, "erro", err)
		return nil
	}

	return contratos
}
