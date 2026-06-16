package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/logger"
)

type BuscarSenadorHandler struct {
	useCase *usecase.BuscarSenadorUseCase
}

func NovoBuscarSenadorHandler(useCase *usecase.BuscarSenadorUseCase) *BuscarSenadorHandler {
	return &BuscarSenadorHandler{useCase: useCase}
}

func (h *BuscarSenadorHandler) Buscar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Buscar")
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo invalido"})
		return
	}

	result, err := h.useCase.Buscar(c.Request.Context(), codigo)
	if err != nil {
		log.Error("erro ao buscar senador", "codigo", codigo, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar senador: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
