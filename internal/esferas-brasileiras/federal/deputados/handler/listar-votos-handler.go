package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarVotosHandler struct {
	useCase *usecase.EsferaFederalListarVotosUseCase
}

func NovoEsferaFederalListarVotosHandler(useCase *usecase.EsferaFederalListarVotosUseCase) *EsferaFederalListarVotosHandler {
	return &EsferaFederalListarVotosHandler{useCase: useCase}
}

func (h *EsferaFederalListarVotosHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarVotos")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao listar votos", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar votos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
