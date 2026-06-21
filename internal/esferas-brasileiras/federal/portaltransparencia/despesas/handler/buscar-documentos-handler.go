package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarDocumentosHandler struct {
	useCase *usecase.BuscarDocumentosUseCase
}

func NovoBuscarDocumentosHandler(useCase *usecase.BuscarDocumentosUseCase) *BuscarDocumentosHandler {
	return &BuscarDocumentosHandler{useCase: useCase}
}

func (h *BuscarDocumentosHandler) BuscarDocumentos(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.DespesaDocumentosQueryParams{
		DataEmissao:    c.Query("dataEmissao"),
		Fase:           c.Query("fase"),
		Pagina:         pagina,
		UnidadeGestora: c.Query("unidadeGestora"),
		Gestao:         c.Query("gestao"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []interface{}{}
	}
	c.JSON(http.StatusOK, resultado)
}
