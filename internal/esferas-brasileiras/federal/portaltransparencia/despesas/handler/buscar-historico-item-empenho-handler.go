package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarHistoricoItemEmpenhoHandler struct {
	useCase *usecase.BuscarHistoricoItemEmpenhoUseCase
}

func NovoBuscarHistoricoItemEmpenhoHandler(useCase *usecase.BuscarHistoricoItemEmpenhoUseCase) *BuscarHistoricoItemEmpenhoHandler {
	return &BuscarHistoricoItemEmpenhoHandler{useCase: useCase}
}

func (h *BuscarHistoricoItemEmpenhoHandler) BuscarHistoricoEmpenho(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	sequencial, _ := strconv.Atoi(c.Query("sequencial"))
	resultado, err := h.useCase.Buscar(c.Request.Context(), c.Query("codigoDocumento"), sequencial, pagina)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.HistoricoSubItemEmpenho{}
	}
	c.JSON(http.StatusOK, resultado)
}
