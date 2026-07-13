package handler

import (
	"net/http"

	senado "github.com/danyele/podp/internal/api/senado"
	"github.com/danyele/podp/internal/api/senado/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarVotacoesComissaoParlamentarHandler struct {
	useCase *usecase.ListarVotacoesComissaoParlamentarUseCase
}

func NovoListarVotacoesComissaoParlamentarHandler(useCase *usecase.ListarVotacoesComissaoParlamentarUseCase) *ListarVotacoesComissaoParlamentarHandler {
	return &ListarVotacoesComissaoParlamentarHandler{useCase: useCase}
}

func (h *ListarVotacoesComissaoParlamentarHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo invalido"})
		return
	}

	votacoes, err := h.useCase.Listar(c.Request.Context(), codigo, senado.QueryParams(c))
	if err != nil {
		log.Error("erro ao listar votacoes comissao parlamentar", "codigo", codigo, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar votacoes parlamentar: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": votacoes})
}
