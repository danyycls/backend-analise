package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarTodasComissoesHandler struct {
	useCase *usecase.ListarTodasComissoesUseCase
}

func NovoListarTodasComissoesHandler(useCase *usecase.ListarTodasComissoesUseCase) *ListarTodasComissoesHandler {
	return &ListarTodasComissoesHandler{useCase: useCase}
}

func (h *ListarTodasComissoesHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	comissoes, err := h.useCase.Listar(c.Request.Context())
	if err != nil {
		log.Error("erro ao listar todas comissoes", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar comissoes: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": comissoes})
}
