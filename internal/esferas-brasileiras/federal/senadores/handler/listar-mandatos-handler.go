package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/logger"
)

type ListarMandatosHandler struct {
	useCase *usecase.ListarMandatosUseCase
}

func NovoListarMandatosHandler(useCase *usecase.ListarMandatosUseCase) *ListarMandatosHandler {
	return &ListarMandatosHandler{useCase: useCase}
}

func (h *ListarMandatosHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo invalido"})
		return
	}

	mandatos, err := h.useCase.Listar(c.Request.Context(), codigo)
	if err != nil {
		log.Error("erro ao listar mandatos senador", "codigo", codigo, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar mandatos: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": mandatos})
}
