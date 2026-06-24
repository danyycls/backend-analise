package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalBuscarFrenteHandler struct {
	useCase *usecase.EsferaFederalBuscarFrenteUseCase
}

func NovoEsferaFederalBuscarFrenteHandler(useCase *usecase.EsferaFederalBuscarFrenteUseCase) *EsferaFederalBuscarFrenteHandler {
	return &EsferaFederalBuscarFrenteHandler{useCase: useCase}
}

func (h *EsferaFederalBuscarFrenteHandler) Buscar(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarFrente")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao buscar frente", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar frente: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
