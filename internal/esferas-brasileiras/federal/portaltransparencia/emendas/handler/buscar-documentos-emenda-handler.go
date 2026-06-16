package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/emendas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarDocumentosEmendaHandler struct {
	useCase *usecase.BuscarDocumentosEmendaUseCase
}

func NovoBuscarDocumentosEmendaHandler(useCase *usecase.BuscarDocumentosEmendaUseCase) *BuscarDocumentosEmendaHandler {
	return &BuscarDocumentosEmendaHandler{useCase: useCase}
}

func (h *BuscarDocumentosEmendaHandler) BuscarDocumentos(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	resultado, err := h.useCase.Buscar(c.Request.Context(), c.Param("codigo"), pagina)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.DocumentoRelacionadoEmenda{}
	}
	c.JSON(http.StatusOK, resultado)
}
