package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/tcu/usecase"

	"github.com/gin-gonic/gin"

	client "github.com/danyele/podp/internal/shared/clients/tcu"
)

type InabilitadosHandler struct {
	useCase *usecase.InabilitadosUseCase
}

func NovoInabilitadosHandler(useCase *usecase.InabilitadosUseCase) *InabilitadosHandler {
	return &InabilitadosHandler{useCase: useCase}
}

func (h *InabilitadosHandler) Buscar(c *gin.Context) {
	var filter client.TCUQueryParams
	if err := c.ShouldBindJSON(&filter); err != nil {
		filter = client.TCUQueryParams{}
	}

	resultado, err := h.useCase.Buscar(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao consultar TCU: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []client.Sancoes{}
	}

	c.JSON(http.StatusOK, resultado)
}
