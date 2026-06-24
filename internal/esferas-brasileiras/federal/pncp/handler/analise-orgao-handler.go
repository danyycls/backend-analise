package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"sync"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/pncp/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"

	pncp "github.com/danyele/podp/internal/shared/clients/pncp"
)

type AnaliseOrgaoPNCPHandler struct {
	useCase *usecase.ConsultaCNPJOrgaoPNCPUseCase
	redis   *redis.RedisCache
	jobsMu  sync.Mutex
	jobs    map[string]*jobState
}

func NovoAnaliseOrgaoPNCPHandler(useCase *usecase.ConsultaCNPJOrgaoPNCPUseCase, redis *redis.RedisCache) *AnaliseOrgaoPNCPHandler {
	return &AnaliseOrgaoPNCPHandler{
		useCase: useCase,
		redis:   redis,
		jobs:    make(map[string]*jobState),
	}
}

type jobState struct {
	eventChan chan pncp.EventoAnalise
	results   []*pncp.AnaliseResultado
	mu        sync.RWMutex
	done      chan struct{}
}

func (h *AnaliseOrgaoPNCPHandler) AnaliseOrgaoPNCP(c *gin.Context) {
	log := logger.New("PNCP: Handler: AnaliseOrgaoPNCP")
	var req pncp.AnaliseOrgaoPNCPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if len(req.CNPJs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "lista de CNPJs nao pode ser vazia"})
		return
	}

	cnpjs := make([]string, len(req.CNPJs))
	copy(cnpjs, req.CNPJs)
	sort.Strings(cnpjs)
	params := map[string]interface{}{
		"cnpjs":       cnpjs,
		"dataInicial": req.DataInicial,
		"dataFinal":   req.DataFinal,
	}
	raw, _ := json.Marshal(params)
	chave := redis.ChaveCache("orgao-analise", raw)

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
	js := &jobState{eventChan: eventChan, done: make(chan struct{})}
	if cacheHit {
		js.results = cached
		close(js.done)
	}
	h.jobs[jobID] = js
	h.jobsMu.Unlock()

	go func() {
		log := logger.New("PNCP: Handler: AnaliseOrgaoPNCP")
		ctx := context.Background()

		if cacheHit {
			for i, r := range cached {
				nomeOrgao := ""
				if r.Orgao != nil && r.Orgao.RazaoSocial != nil {
					nomeOrgao = *r.Orgao.RazaoSocial
				}
				eventChan <- pncp.EventoAnalise{Type: "started", CNPJ: req.CNPJs[i], Total: len(req.CNPJs)}
				eventChan <- pncp.EventoAnalise{Type: "success", CNPJ: req.CNPJs[i], Orgao: nomeOrgao}
			}
			eventChan <- pncp.EventoAnalise{
				Type:      "progress",
				Processed: len(req.CNPJs),
				Total:     len(req.CNPJs),
				Success:   len(req.CNPJs),
			}
			eventChan <- pncp.EventoAnalise{Type: "completed", Total: len(req.CNPJs)}
			close(eventChan)
			return
		}

		results := h.useCase.AnaliseMultiplos(ctx, req, eventChan)

		h.jobsMu.Lock()
		if job, ok := h.jobs[jobID]; ok {
			job.mu.Lock()
			job.results = results
			job.mu.Unlock()
			close(job.done)
		}
		h.jobsMu.Unlock()
		close(eventChan)

		if results != nil {
			if err := h.redis.Set(ctx, chave, results); err != nil {
				log.Warn("cache indisponivel", "erro", err)
			}
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"jobId":  jobID,
		"status": "processing",
		"total":  len(req.CNPJs),
	})
}

func (h *AnaliseOrgaoPNCPHandler) BuscarResultadosBatch(c *gin.Context) {
	jobID := c.Param("jobId")

	h.jobsMu.Lock()
	job, exists := h.jobs[jobID]
	h.jobsMu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"erro": "job nao encontrado"})
		return
	}

	<-job.done

	job.mu.RLock()
	results := job.results
	job.mu.RUnlock()

	status := "processing"
	if results != nil {
		status = "completed"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"results": results,
	})
}

func (h *AnaliseOrgaoPNCPHandler) Jobs() map[string]*jobState {
	return h.jobs
}

func (h *AnaliseOrgaoPNCPHandler) GetJobChan(jobID string) (<-chan pncp.EventoAnalise, bool) {
	h.jobsMu.Lock()
	defer h.jobsMu.Unlock()
	job, exists := h.jobs[jobID]
	if !exists {
		return nil, false
	}
	return job.eventChan, true
}
