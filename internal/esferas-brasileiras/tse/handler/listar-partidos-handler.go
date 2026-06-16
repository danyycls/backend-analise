package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/tse/usecase"

	"github.com/gin-gonic/gin"
)

type ListarPartidosHandler struct {
	useCase *usecase.BuscarCandidatosUseCase
}

func NovoListarPartidosHandler(useCase *usecase.BuscarCandidatosUseCase) *ListarPartidosHandler {
	return &ListarPartidosHandler{useCase: useCase}
}

func (h *ListarPartidosHandler) ListarPartidos(c *gin.Context) {
	result, err := h.useCase.ExecutarListarPartidos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
