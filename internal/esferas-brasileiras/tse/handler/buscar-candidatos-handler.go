package handler

import (
	"net/http"

	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/esferas-brasileiras/tse/usecase"

	"github.com/gin-gonic/gin"
)

type BuscarCandidatosHandler struct {
	useCase *usecase.BuscarCandidatosUseCase
}

func NovoBuscarCandidatosHandler(useCase *usecase.BuscarCandidatosUseCase) *BuscarCandidatosHandler {
	return &BuscarCandidatosHandler{useCase: useCase}
}

func (h *BuscarCandidatosHandler) BuscarCandidatos(c *gin.Context) {
	var req tsetypes.BuscaCandidatosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	result, err := h.useCase.Executar(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
