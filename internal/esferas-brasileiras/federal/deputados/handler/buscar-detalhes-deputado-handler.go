package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalBuscarDetalhesDeputadoHandler struct {
	useCase *usecase.EsferaFederalBuscarDetalhesDeputadoUseCase
}

func NovoEsferaFederalBuscarDetalhesDeputadoHandler(useCase *usecase.EsferaFederalBuscarDetalhesDeputadoUseCase) *EsferaFederalBuscarDetalhesDeputadoHandler {
	return &EsferaFederalBuscarDetalhesDeputadoHandler{useCase: useCase}
}

func (h *EsferaFederalBuscarDetalhesDeputadoHandler) BuscarDetalhesDeputado(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarDetalhesDeputado")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaFederalBuscarDetalhesDeputadoRequest{ID: id})
	if err != nil {
		log.Error("erro ao buscar deputado", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar deputado: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Deputado)
}
