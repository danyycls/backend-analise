package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/danyele/podp/internal/shared/clients/opencnpj"
	"github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/shared/utils"
)

const maxConcorrencia = 5

type ConsultaCNPJOrgaoPNCPUseCase struct {
	pncpClient     *pncp.PNCPClient
	opencnpjClient *opencnpj.OpenCNPJClient
	redis          *redis.RedisCache
	httpClient     *http.Client
}

func NovoConsultaCNPJOrgaoPNCPUseCase(pncp *pncp.PNCPClient, opencnpj *opencnpj.OpenCNPJClient, redis *redis.RedisCache) *ConsultaCNPJOrgaoPNCPUseCase {
	return &ConsultaCNPJOrgaoPNCPUseCase{
		pncpClient:     pncp,
		opencnpjClient: opencnpj,
		redis:          redis,
		httpClient:     &http.Client{Timeout: 15 * time.Second},
	}
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) AnaliseMultiplos(ctx context.Context, req pncp.AnaliseOrgaoPNCPRequest, eventos chan<- pncp.EventoAnalise) []*pncp.AnaliseResultado {
	log := logger.New("PNCP: UseCase: AnaliseMultiplos")
	defer close(eventos)

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
			log := logger.New("PNCP: UseCase: AnaliseMultiplos")
			defer wg.Done()
			defer func() { <-sem }()

			log.Info("iniciando analise", "cnpj", cnpj)
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
				nomeOrgao := ""
				if result != nil && result.Orgao != nil && result.Orgao.RazaoSocial != nil {
					nomeOrgao = *result.Orgao.RazaoSocial
				}
				totalContratos := 0
				var valorTotal float64
				if result != nil && result.Resumo != nil {
					if result.Resumo.TotalContratos != nil {
						totalContratos = *result.Resumo.TotalContratos
					}
					if result.Resumo.ValorTotalContratos != nil {
						valorTotal = *result.Resumo.ValorTotalContratos
					}
				}
				log.Info("analise concluida com sucesso", "cnpj", cnpj, "orgao", nomeOrgao, "contratos", totalContratos)
				eventos <- pncp.EventoAnalise{
					Type:                "success",
					CNPJ:                cnpj,
					Orgao:               nomeOrgao,
					TotalContratos:      totalContratos,
					ValorTotalContratos: valorTotal,
				}
			}
			eventos <- pncp.EventoAnalise{
				Type:      "progress",
				Processed: processed,
				Total:     total,
				Success:   success,
				Errors:    errors,
			}
			mu.Unlock()
		}(cnpj)
	}

	wg.Wait()
	log.Info("analise multipla concluida", "total", total, "success", success, "errors", errors)
	eventos <- pncp.EventoAnalise{Type: "completed", Total: total}
	return results
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) executar(ctx context.Context, cnpjOrgao, dataInicial, dataFinal string) (*pncp.AnaliseResultado, error) {
	log := logger.New("PNCP: UseCase: executar")
	cnpj := utils.NormalizarCNPJ(cnpjOrgao)

	cacheParams := map[string]interface{}{
		"cnpj":        cnpj,
		"dataInicial": dataInicial,
		"dataFinal":   dataFinal,
	}
	raw, _ := json.Marshal(cacheParams)
	cacheKey := redis.ChaveCache("pncp-orgao-executar", raw)

	var cached pncp.AnaliseResultado
	if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	} else if ok {
		log.Info("cache hit", "cnpj", cnpj)
		return &cached, nil
	}

	log.Info("iniciando analise orgao", "cnpj", cnpj, "data_inicial", dataInicial, "data_final", dataFinal)

	contratosRaw := u.buscarContratos(ctx, cnpj, dataInicial, dataFinal)
	if contratosRaw == nil {
		return nil, fmt.Errorf("erro ao consultar contratos para CNPJ %s", cnpj)
	}

	log.Info("contratos encontrados", "total", len(contratosRaw))
	fornecedoresMap := u.extrairFornecedores(contratosRaw)
	log.Info("fornecedores unicos extraidos", "total", len(fornecedoresMap))

	enrichedFornecedores := u.enriquecerFornecedores(ctx, fornecedoresMap)
	log.Info("empresas enriquecidas", "total", len(enrichedFornecedores))

	contratosDTO := u.montarContratosDTO(contratosRaw, enrichedFornecedores)
	totalContratos, totalEmpresas, valorTotal := u.calcularResumo(contratosDTO, enrichedFornecedores)

	orgao := u.montarOrgaoInfo(cnpj, contratosDTO)

	result := &pncp.AnaliseResultado{
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

	if err := u.redis.Set(ctx, cacheKey, result); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}

	return result, nil
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) buscarContratos(ctx context.Context, cnpj, dataInicial, dataFinal string) []pncp.Contrato {
	log := logger.New("PNCP: UseCase: buscarContratos")
	pagina := 1
	tamanho := 500
	contratos := make([]pncp.Contrato, 0)
	seenContratos := make(map[string]struct{})

	for {
		cacheParams := map[string]interface{}{
			"cnpj":        cnpj,
			"dataInicial": dataInicial,
			"dataFinal":   dataFinal,
			"pagina":      pagina,
			"tamanho":     tamanho,
		}
		raw, _ := json.Marshal(cacheParams)
		cacheKey := redis.ChaveCache("pncp-buscar-contratos", raw)

		var cached []pncp.Contrato
		cacheHit := false
		if ok, err := u.redis.Get(ctx, cacheKey, &cached); err != nil {
			log.Warn("cache indisponivel", "erro", err)
		} else if ok {
			cacheHit = true
			log.Info("cache hit pagina", "pagina", pagina)
		}

		var resp []pncp.Contrato
		var err error
		if cacheHit {
			resp = cached
		} else {
			resp, err = u.pncpClient.BuscarContratos(ctx, cnpj, dataInicial, dataFinal, pagina, tamanho)
			if err != nil {
				log.Error("erro ao consultar PNCP", "pagina", pagina, "erro", err)
				return nil
			}
			if setErr := u.redis.Set(ctx, cacheKey, resp); setErr != nil {
				log.Warn("cache indisponivel", "erro", setErr)
			}
		}

		log.Info("contratos por pagina", "pagina", pagina, "total", len(resp))
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
			if _, ok := seenContratos[key]; ok {
				continue
			}
			seenContratos[key] = struct{}{}
			contratos = append(contratos, c)
		}
		if len(resp) < tamanho {
			break
		}
		pagina++
	}
	return contratos
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) extrairFornecedores(contratos []pncp.Contrato) map[string]string {
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

func (u *ConsultaCNPJOrgaoPNCPUseCase) enriquecerFornecedores(ctx context.Context, fornecedoresMap map[string]string) map[string]*types.FornecedorOpenCNPJ {
	log := logger.New("PNCP: UseCase: enriquecerFornecedores")
	enriched := make(map[string]*types.FornecedorOpenCNPJ, len(fornecedoresMap))
	for cnpjF, nome := range fornecedoresMap {
		data, err := u.opencnpjClient.Buscar(ctx, cnpjF)
		if err != nil {
			log.Error("erro ao consultar OpenCNPJ", "cnpj", cnpjF, "erro", err)
			enriched[cnpjF] = &types.FornecedorOpenCNPJ{CNPJ: pncp.StrPtr(cnpjF), RazaoSocial: pncp.StrPtr(nome)}
			continue
		}
		enriched[cnpjF] = utils.BuildFornecedorDTO(data)
	}
	return enriched
}

func (u *ConsultaCNPJOrgaoPNCPUseCase) montarContratosDTO(contratos []pncp.Contrato, enrichedFornecedores map[string]*types.FornecedorOpenCNPJ) []pncp.Contrato {
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
