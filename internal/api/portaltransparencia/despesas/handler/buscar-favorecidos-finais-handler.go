package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarFavorecidosFinaisPorDocumentoHandler struct {
	useCase *usecase.BuscarFavorecidosFinaisPorDocumentoUseCase
}

func NovoBuscarFavorecidosFinaisPorDocumentoHandler(useCase *usecase.BuscarFavorecidosFinaisPorDocumentoUseCase) *BuscarFavorecidosFinaisPorDocumentoHandler {
	return &BuscarFavorecidosFinaisPorDocumentoHandler{useCase: useCase}
}

func (h *BuscarFavorecidosFinaisPorDocumentoHandler) BuscarFavorecidosFinaisPorDocumento(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	resultado, err := h.useCase.Buscar(c.Request.Context(), c.Query("codigoDocumento"), pagina)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.ConsultaFavorecidosFinaisPorDocumento{}
	}
	c.JSON(http.StatusOK, resultado)
}
