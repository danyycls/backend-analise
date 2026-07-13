package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/ibge/usecase"
)

type ListarEstadosHandler struct {
	useCase *usecase.ListarEstadosUseCase
}

func NovoListarEstadosHandler(useCase *usecase.ListarEstadosUseCase) *ListarEstadosHandler {
	return &ListarEstadosHandler{useCase: useCase}
}

func (h *ListarEstadosHandler) ListarEstados(c *gin.Context) {
	estados, err := h.useCase.Executar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, estados)
}
