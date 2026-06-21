package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarDocumentosRelacionadosHandler struct {
	useCase *usecase.BuscarDocumentosRelacionadosUseCase
}

func NovoBuscarDocumentosRelacionadosHandler(useCase *usecase.BuscarDocumentosRelacionadosUseCase) *BuscarDocumentosRelacionadosHandler {
	return &BuscarDocumentosRelacionadosHandler{useCase: useCase}
}

func (h *BuscarDocumentosRelacionadosHandler) BuscarDocumentosRelacionados(c *gin.Context) {
	resultado, err := h.useCase.Buscar(c.Request.Context(), c.Query("codigoDocumento"), c.Query("fase"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.DocumentoRelacionado{}
	}
	c.JSON(http.StatusOK, resultado)
}
