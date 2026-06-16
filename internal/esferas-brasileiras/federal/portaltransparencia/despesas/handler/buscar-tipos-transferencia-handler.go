package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarTiposTransferenciaHandler struct {
	useCase *usecase.BuscarTiposTransferenciaUseCase
}

func NovoBuscarTiposTransferenciaHandler(useCase *usecase.BuscarTiposTransferenciaUseCase) *BuscarTiposTransferenciaHandler {
	return &BuscarTiposTransferenciaHandler{useCase: useCase}
}

func (h *BuscarTiposTransferenciaHandler) BuscarTiposTransferencia(c *gin.Context) {
	resultado, err := h.useCase.Buscar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.CodigoDescricao{}
	}
	c.JSON(http.StatusOK, resultado)
}
