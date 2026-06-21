package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarComissoesHandler struct {
	useCase *usecase.ListarComissoesUseCase
}

func NovoListarComissoesHandler(useCase *usecase.ListarComissoesUseCase) *ListarComissoesHandler {
	return &ListarComissoesHandler{useCase: useCase}
}

func (h *ListarComissoesHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo invalido"})
		return
	}

	comissoes, err := h.useCase.Listar(c.Request.Context(), codigo)
	if err != nil {
		log.Error("erro ao listar comissoes senador", "codigo", codigo, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar comissoes: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": comissoes})
}
