package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/logger"
)

type ListarProcessoAssuntosHandler struct {
	useCase *usecase.ListarProcessoAssuntosUseCase
}

func NovoListarProcessoAssuntosHandler(useCase *usecase.ListarProcessoAssuntosUseCase) *ListarProcessoAssuntosHandler {
	return &ListarProcessoAssuntosHandler{useCase: useCase}
}

func (h *ListarProcessoAssuntosHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	assuntos, err := h.useCase.Listar(c.Request.Context())
	if err != nil {
		log.Error("erro ao listar assuntos processo", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar assuntos: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": assuntos})
}
