package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/pncp/usecase"
	pncp "github.com/danyele/podp/internal/shared/clients/pncp"
)

type AnaliseOrgaoPNCPHandler struct {
	*JobManager
	useCase *usecase.ConsultaContratoOrgaoPNCPUseCase
}

func NovoAnaliseOrgaoPNCPHandler(useCase *usecase.ConsultaContratoOrgaoPNCPUseCase) *AnaliseOrgaoPNCPHandler {
	return &AnaliseOrgaoPNCPHandler{
		JobManager: NovoJobManager(),
		useCase:    useCase,
	}
}

func (h *AnaliseOrgaoPNCPHandler) AnaliseOrgaoPNCP(c *gin.Context) {
	var req pncp.AnaliseContratoOrgaoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if len(req.CNPJs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "lista de CNPJs nao pode ser vazia"})
		return
	}

	jobID, eventChan, _ := h.CriarJob()

	go func() {
		ctx := context.Background()
		results := h.useCase.Executar(ctx, req, eventChan)

		eventChan <- pncp.EventoAnalise{
			Type:    "results",
			Results: results,
		}

		h.FinalizarJob(jobID, results, nil)
		close(eventChan)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"jobId":  jobID,
		"status": "processing",
		"total":  len(req.CNPJs),
	})
}
