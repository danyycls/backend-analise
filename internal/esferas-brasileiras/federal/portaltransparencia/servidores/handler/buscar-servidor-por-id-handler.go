package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/servidores/usecase"

	"github.com/gin-gonic/gin"
)

type BuscarServidorPorIDHandler struct {
	useCase *usecase.BuscarServidorPorIDUseCase
}

func NovoBuscarServidorPorIDHandler(useCase *usecase.BuscarServidorPorIDUseCase) *BuscarServidorPorIDHandler {
	return &BuscarServidorPorIDHandler{useCase: useCase}
}

func (h *BuscarServidorPorIDHandler) BuscarPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}
