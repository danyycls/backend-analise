package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/esferas-brasileiras/tse/usecase"

	"github.com/gin-gonic/gin"
)

type ListarCargosHandler struct {
	useCase *usecase.BuscarCandidatosUseCase
}

func NovoListarCargosHandler(useCase *usecase.BuscarCandidatosUseCase) *ListarCargosHandler {
	return &ListarCargosHandler{useCase: useCase}
}

func (h *ListarCargosHandler) ListarCargos(c *gin.Context) {
	result, err := h.useCase.ExecutarListarCargos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
