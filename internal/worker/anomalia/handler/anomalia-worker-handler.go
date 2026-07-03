package handler

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/logger"

	anomalia "github.com/danyele/podp/internal/worker/anomalia"
	"github.com/danyele/podp/internal/worker/anomalia/usecase"
)

type workerJobState struct {
	cancel    context.CancelFunc
	eventChan chan anomalia.WorkerEvento
	result    anomalia.WorkerProgressoResponse
	mu        sync.RWMutex
}

type AnomaliaWorkerHandler struct {
	useCase *usecase.AnaliseAnomaliaWorkerUseCase
	jobsMu  sync.Mutex
	jobs    map[string]*workerJobState
}

func NovoAnomaliaWorkerHandler(useCase *usecase.AnaliseAnomaliaWorkerUseCase) *AnomaliaWorkerHandler {
	return &AnomaliaWorkerHandler{
		useCase: useCase,
		jobs:    make(map[string]*workerJobState),
	}
}

func (h *AnomaliaWorkerHandler) Iniciar(c *gin.Context) {
	var req anomalia.IniciarWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "corpo inválido: " + err.Error()})
		return
	}

	if len(req.Licitacoes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "lista de licitacoes nao pode ser vazia"})
		return
	}

	jobID := uuid.New().String()
	ctx, cancel := context.WithCancel(context.Background())
	eventChan := make(chan anomalia.WorkerEvento, 500)

	h.jobsMu.Lock()
	h.jobs[jobID] = &workerJobState{
		cancel:    cancel,
		eventChan: eventChan,
		result: anomalia.WorkerProgressoResponse{
			JobID:      jobID,
			Status:     "processing",
			EtapaAtual: "analisando_vinculos",
		},
	}
	h.jobsMu.Unlock()

	go func() {
		defer close(eventChan)

		h.useCase.Executar(ctx, req, eventChan)

		h.jobsMu.Lock()
		if job, ok := h.jobs[jobID]; ok {
			job.mu.Lock()
			if job.result.Status != "cancelled" {
				job.result.Status = "completed"
			}
			job.mu.Unlock()
		}
		h.jobsMu.Unlock()
	}()

	go func() {
		for evento := range eventChan {
			h.jobsMu.Lock()
			if job, ok := h.jobs[jobID]; ok {
				job.mu.Lock()
				if evento.EtapaAtual != "" {
					job.result.EtapaAtual = evento.EtapaAtual
				}
				switch evento.Type {
				case "started":
					job.result.Total = evento.Total
				case "progress":
					job.result.Processed = evento.Processed
					job.result.Total = evento.Total
					job.result.Success = evento.Success
					job.result.Errors = evento.Errors
					job.result.AnomaliasEncontradas = evento.AnomaliasEncontradas
				case "error":
					job.result.Message = evento.Message
					job.result.Errors++
				case "completed":
					job.result.Status = "completed"
					job.result.Processed = evento.Processed
					job.result.Total = evento.Total
					job.result.Success = evento.Success
					job.result.Errors = evento.Errors
					job.result.AnomaliasEncontradas = evento.AnomaliasEncontradas
					job.result.EtapaAtual = "concluido"
				}
				job.mu.Unlock()
			}
			h.jobsMu.Unlock()
		}
	}()

	logger.Info("job criado", "job_id", jobID, "licitacoes", len(req.Licitacoes))
	c.JSON(http.StatusAccepted, gin.H{
		"job_id":  jobID,
		"status":  "processing",
		"message": "worker iniciado",
	})
}

func (h *AnomaliaWorkerHandler) Parar(c *gin.Context) {
	jobID := c.Param("jobId")

	h.jobsMu.Lock()
	job, exists := h.jobs[jobID]
	h.jobsMu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"erro": "job nao encontrado"})
		return
	}

	job.mu.Lock()
	if job.result.Status == "completed" {
		job.mu.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{"erro": "job ja concluido"})
		return
	}
	job.result.Status = "cancelled"
	job.mu.Unlock()

	job.cancel()

	logger.Info("job cancelado", "job_id", jobID)
	c.JSON(http.StatusOK, gin.H{
		"job_id":  jobID,
		"status":  "cancelled",
		"message": "worker parado",
	})
}

func (h *AnomaliaWorkerHandler) Progression(c *gin.Context) {
	jobID := c.Param("jobId")

	h.jobsMu.Lock()
	job, exists := h.jobs[jobID]
	h.jobsMu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"erro": "job nao encontrado"})
		return
	}

	job.mu.RLock()
	result := job.result
	job.mu.RUnlock()

	c.JSON(http.StatusOK, result)
}

func (h *AnomaliaWorkerHandler) GetJobChan(jobID string) (<-chan anomalia.WorkerEvento, bool) {
	h.jobsMu.Lock()
	defer h.jobsMu.Unlock()
	job, exists := h.jobs[jobID]
	if !exists {
		return nil, false
	}
	return job.eventChan, true
}
