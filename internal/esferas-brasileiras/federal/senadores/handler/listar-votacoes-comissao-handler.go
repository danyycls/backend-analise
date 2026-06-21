package handler

import (
	"net/http"

	senadores "github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores"
	"github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarVotacoesComissaoHandler struct {
	useCase *usecase.ListarVotacoesComissaoUseCase
}

func NovoListarVotacoesComissaoHandler(useCase *usecase.ListarVotacoesComissaoUseCase) *ListarVotacoesComissaoHandler {
	return &ListarVotacoesComissaoHandler{useCase: useCase}
}

func (h *ListarVotacoesComissaoHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	sigla := c.Param("sigla")
	if sigla == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "sigla invalida"})
		return
	}

	votacoes, err := h.useCase.Listar(c.Request.Context(), sigla, senadores.QueryParams(c))
	if err != nil {
		log.Error("erro ao listar votacoes comissao", "sigla", sigla, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar votacoes comissao: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": votacoes})
}
