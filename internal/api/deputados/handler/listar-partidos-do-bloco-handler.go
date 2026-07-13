package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarPartidosDoBlocoHandler struct {
	useCase *usecase.EsferaFederalListarPartidosDoBlocoUseCase
}

func NovoEsferaFederalListarPartidosDoBlocoHandler(useCase *usecase.EsferaFederalListarPartidosDoBlocoUseCase) *EsferaFederalListarPartidosDoBlocoHandler {
	return &EsferaFederalListarPartidosDoBlocoHandler{useCase: useCase}
}

func (h *EsferaFederalListarPartidosDoBlocoHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarPartidosDoBloco")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao listar partidos do bloco", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar partidos do bloco: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
