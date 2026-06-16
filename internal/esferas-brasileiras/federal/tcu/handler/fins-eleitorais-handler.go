package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/tcu/usecase"

	"github.com/gin-gonic/gin"

	client "github.com/danyele/laceu/internal/shared/clients/tcu"
)

type FinsEleitoraisHandler struct {
	useCase *usecase.FinsEleitoraisUseCase
}

func NovoFinsEleitoraisHandler(useCase *usecase.FinsEleitoraisUseCase) *FinsEleitoraisHandler {
	return &FinsEleitoraisHandler{useCase: useCase}
}

func (h *FinsEleitoraisHandler) Buscar(c *gin.Context) {
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
		resultado = []client.FinsEleitorais{}
	}

	c.JSON(http.StatusOK, resultado)
}
