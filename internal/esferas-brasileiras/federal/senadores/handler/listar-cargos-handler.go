package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarCargosHandler struct {
	useCase *usecase.ListarCargosUseCase
}

func NovoListarCargosHandler(useCase *usecase.ListarCargosUseCase) *ListarCargosHandler {
	return &ListarCargosHandler{useCase: useCase}
}

func (h *ListarCargosHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo invalido"})
		return
	}

	cargos, err := h.useCase.Listar(c.Request.Context(), codigo)
	if err != nil {
		log.Error("erro ao listar cargos senador", "codigo", codigo, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar cargos: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": cargos})
}
