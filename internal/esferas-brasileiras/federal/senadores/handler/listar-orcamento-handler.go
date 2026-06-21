package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarOrcamentoHandler struct {
	useCase *usecase.ListarOrcamentoUseCase
}

func NovoListarOrcamentoHandler(useCase *usecase.ListarOrcamentoUseCase) *ListarOrcamentoHandler {
	return &ListarOrcamentoHandler{useCase: useCase}
}

func (h *ListarOrcamentoHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	itens, err := h.useCase.Listar(c.Request.Context())
	if err != nil {
		log.Error("erro ao listar orcamento", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar orcamento: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": itens})
}
