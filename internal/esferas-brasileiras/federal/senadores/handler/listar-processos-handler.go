package handler

import (
	"net/http"

	senadores "github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores"
	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/logger"
)

type ListarProcessosHandler struct {
	useCase *usecase.ListarProcessosUseCase
}

func NovoListarProcessosHandler(useCase *usecase.ListarProcessosUseCase) *ListarProcessosHandler {
	return &ListarProcessosHandler{useCase: useCase}
}

func (h *ListarProcessosHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	processors, err := h.useCase.Listar(c.Request.Context(), senadores.QueryParams(c))
	if err != nil {
		log.Error("erro ao listar processors", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar processors: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": processors})
}
