package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarItensEmpenhoHandler struct {
	useCase *usecase.BuscarItensEmpenhoUseCase
}

func NovoBuscarItensEmpenhoHandler(useCase *usecase.BuscarItensEmpenhoUseCase) *BuscarItensEmpenhoHandler {
	return &BuscarItensEmpenhoHandler{useCase: useCase}
}

func (h *BuscarItensEmpenhoHandler) BuscarItensEmpenho(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	resultado, err := h.useCase.Buscar(c.Request.Context(), c.Query("codigoDocumento"), pagina)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.DetalhamentoDoGasto{}
	}
	c.JSON(http.StatusOK, resultado)
}
