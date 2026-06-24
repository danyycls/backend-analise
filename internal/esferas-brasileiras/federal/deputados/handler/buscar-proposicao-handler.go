package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type BuscarProposicaoHandler struct {
	useCase *usecase.BuscarProposicaoUseCase
}

func NovoBuscarProposicaoHandler(useCase *usecase.BuscarProposicaoUseCase) *BuscarProposicaoHandler {
	return &BuscarProposicaoHandler{useCase: useCase}
}

func (h *BuscarProposicaoHandler) Buscar(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarProposicao")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao buscar proposicao", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar proposicao: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
