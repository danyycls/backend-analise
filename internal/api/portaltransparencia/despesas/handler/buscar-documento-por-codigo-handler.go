package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"
)

type BuscarDocumentoPorCodigoHandler struct {
	useCase *usecase.BuscarDocumentoPorCodigoUseCase
}

func NovoBuscarDocumentoPorCodigoHandler(useCase *usecase.BuscarDocumentoPorCodigoUseCase) *BuscarDocumentoPorCodigoHandler {
	return &BuscarDocumentoPorCodigoHandler{useCase: useCase}
}

func (h *BuscarDocumentoPorCodigoHandler) BuscarDocumentoPorCodigo(c *gin.Context) {
	resultado, err := h.useCase.Buscar(c.Request.Context(), c.Param("codigo"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}
