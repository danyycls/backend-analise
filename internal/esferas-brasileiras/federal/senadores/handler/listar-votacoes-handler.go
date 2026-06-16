package handler

import (
	"net/http"

	senadores "github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores"
	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/logger"
)

type ListarVotacoesHandler struct {
	useCase *usecase.ListarVotacoesUseCase
}

func NovoListarVotacoesHandler(useCase *usecase.ListarVotacoesUseCase) *ListarVotacoesHandler {
	return &ListarVotacoesHandler{useCase: useCase}
}

func (h *ListarVotacoesHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	votacoes, err := h.useCase.Listar(c.Request.Context(), senadores.QueryParams(c))
	if err != nil {
		log.Error("erro ao listar votacoes", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar votacoes: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": votacoes})
}
