package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalBuscarVotacaoHandler struct {
	useCase *usecase.EsferaFederalBuscarVotacaoUseCase
}

func NovoEsferaFederalBuscarVotacaoHandler(useCase *usecase.EsferaFederalBuscarVotacaoUseCase) *EsferaFederalBuscarVotacaoHandler {
	return &EsferaFederalBuscarVotacaoHandler{useCase: useCase}
}

func (h *EsferaFederalBuscarVotacaoHandler) Buscar(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarVotacao")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao buscar votacao", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar votacao: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
