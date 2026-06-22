package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/pncp/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"

	pncp "github.com/danyele/podp/internal/shared/clients/pncp"
)

type pubJobState struct {
	eventChan   chan pncp.EventoAnalise
	results     []*pncp.AnaliseResultado
	paginasErro []int
	mu          sync.RWMutex
}

type AnalisePublicacaoHandler struct {
	useCase *usecase.ConsultaPublicacaoPNCPUseCase
	redis   *redis.RedisCache
	jobsMu  sync.Mutex
	jobs    map[string]*pubJobState
}

func NovoAnalisePublicacaoHandler(useCase *usecase.ConsultaPublicacaoPNCPUseCase, redis *redis.RedisCache) *AnalisePublicacaoHandler {
	return &AnalisePublicacaoHandler{
		useCase: useCase,
		redis:   redis,
		jobs:    make(map[string]*pubJobState),
	}
}

func (h *AnalisePublicacaoHandler) AnalisePublicacao(c *gin.Context) {
	log := logger.New("PNCP: Handler: AnalisePublicacao")
	var req pncp.AnalisePublicacaoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if req.Tipo != "uf" && req.Tipo != "municipio" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "tipo deve ser 'uf' ou 'municipio'"})
		return
	}

	var valor string
	if req.Tipo == "municipio" {
		valor = req.CodigoMunicipioIbge
	} else {
		valor = req.UF
	}
	params := map[string]interface{}{
		"tipo":                        req.Tipo,
		"valor":                       valor,
		"dataInicial":                 req.DataInicial,
		"dataFinal":                   req.DataFinal,
		"codigoModalidadeContratacao": req.CodigoModalidadeContratacao,
	}
	raw, _ := json.Marshal(params)
	chave := redis.ChaveCache("publicacao-analise", raw)

	var cached []*pncp.AnaliseResultado
	cacheHit := false
	if ok, err := h.redis.Get(c.Request.Context(), chave, &cached); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	} else if ok {
		cacheHit = true
	}

	jobID := uuid.New().String()
	eventChan := make(chan pncp.EventoAnalise, 200)

	h.jobsMu.Lock()
	js := &pubJobState{eventChan: eventChan}
	if cacheHit {
		js.results = cached
	}
	h.jobs[jobID] = js
	h.jobsMu.Unlock()

	go func() {
		log := logger.New("PNCP: Handler: AnalisePublicacao")
		eventos := eventChan
		defer close(eventos)

		if cacheHit {
			eventos <- pncp.EventoAnalise{Type: "started", CNPJ: req.Tipo, Total: 1}
			eventos <- pncp.EventoAnalise{Type: "success", CNPJ: req.UF, Orgao: fmt.Sprintf("Busca por %s (cache)", req.Tipo)}
			eventos <- pncp.EventoAnalise{Type: "progress", Processed: 1, Total: 1, Success: 1}
			eventos <- pncp.EventoAnalise{Type: "completed", Total: 1}
			return
		}

		eventos <- pncp.EventoAnalise{Type: "started", CNPJ: req.Tipo, Total: 1}

		var results []*pncp.AnaliseResultado
		var err error
		var paginasErro []int

		if req.Tipo == "municipio" {
			results, err = h.useCase.BuscarPorMunicipio(context.Background(), req.CodigoMunicipioIbge, req.DataInicial, req.DataFinal, req.CodigoModalidadeContratacao, &paginasErro)
		} else {
			results, err = h.useCase.BuscarPorUF(context.Background(), req.UF, req.DataInicial, req.DataFinal, req.CodigoModalidadeContratacao, &paginasErro)
		}

		if err != nil {
			eventos <- pncp.EventoAnalise{Type: "error", Message: err.Error()}
			results = []*pncp.AnaliseResultado{}
		} else {
			totalContratos := 0
			totalEmpresas := 0
			var valorTotal float64
			for _, r := range results {
				if r.Resumo != nil {
					if r.Resumo.TotalContratos != nil {
						totalContratos += *r.Resumo.TotalContratos
					}
					if r.Resumo.TotalEmpresas != nil {
						totalEmpresas += *r.Resumo.TotalEmpresas
					}
					if r.Resumo.ValorTotalContratos != nil {
						valorTotal += *r.Resumo.ValorTotalContratos
					}
				}
			}

			eventos <- pncp.EventoAnalise{
				Type:                "success",
				TotalContratos:      totalContratos,
				ValorTotalContratos: valorTotal,
				CNPJ:                req.UF,
				Orgao:               fmt.Sprintf("Busca por %s", req.Tipo),
			}
		}

		eventos <- pncp.EventoAnalise{
			Type:      "progress",
			Processed: 1,
			Total:     1,
			Success:   1,
			Errors:    0,
		}

		eventos <- pncp.EventoAnalise{Type: "completed", Total: 1}

		if results != nil && err == nil {
			if err := h.redis.Set(context.Background(), chave, results); err != nil {
				log.Warn("cache indisponivel", "erro", err)
			}
		}

		h.jobsMu.Lock()
		if job, ok := h.jobs[jobID]; ok {
			job.mu.Lock()
			job.results = results
			job.paginasErro = paginasErro
			job.mu.Unlock()
		}
		h.jobsMu.Unlock()
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"jobId":  jobID,
		"status": "processing",
		"total":  1,
	})
}

func (h *AnalisePublicacaoHandler) BuscarResultadosBatch(c *gin.Context) {
	jobID := c.Param("jobId")

	h.jobsMu.Lock()
	job, exists := h.jobs[jobID]
	h.jobsMu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"erro": "job nao encontrado"})
		return
	}

	job.mu.RLock()
	results := job.results
	paginasErro := job.paginasErro
	job.mu.RUnlock()

	status := "processing"
	if results != nil {
		status = "completed"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      status,
		"results":     results,
		"paginasErro": paginasErro,
	})
}

func (h *AnalisePublicacaoHandler) PubJobs() map[string]*pubJobState {
	return h.jobs
}
