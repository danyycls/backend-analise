package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarDocumentosPorFavorecidoHandler struct {
	useCase *usecase.BuscarDocumentosPorFavorecidoUseCase
}

func NovoBuscarDocumentosPorFavorecidoHandler(useCase *usecase.BuscarDocumentosPorFavorecidoUseCase) *BuscarDocumentosPorFavorecidoHandler {
	return &BuscarDocumentosPorFavorecidoHandler{useCase: useCase}
}

func (h *BuscarDocumentosPorFavorecidoHandler) BuscarDocumentosPorFavorecido(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.DespesaDocumentosPorFavorecidoQueryParams{
		CodigoPessoa:       c.Query("codigoPessoa"),
		Fase:               c.Query("fase"),
		Ano:                c.Query("ano"),
		Pagina:             pagina,
		UG:                 c.Query("ug"),
		Gestao:             c.Query("gestao"),
		OrdenacaoResultado: c.Query("ordenacaoResultado"),
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
