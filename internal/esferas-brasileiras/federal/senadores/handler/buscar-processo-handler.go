package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/logger"
)

type BuscarProcessoHandler struct {
	useCase *usecase.BuscarProcessoUseCase
}

func NovoBuscarProcessoHandler(useCase *usecase.BuscarProcessoUseCase) *BuscarProcessoHandler {
	return &BuscarProcessoHandler{useCase: useCase}
}

func (h *BuscarProcessoHandler) Buscar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Buscar")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	processo, err := h.useCase.Buscar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao buscar processo", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar processo: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, processo)
}
