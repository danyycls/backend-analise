package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalBuscarBlocoHandler struct {
	useCase *usecase.EsferaFederalBuscarBlocoUseCase
}

func NovoEsferaFederalBuscarBlocoHandler(useCase *usecase.EsferaFederalBuscarBlocoUseCase) *EsferaFederalBuscarBlocoHandler {
	return &EsferaFederalBuscarBlocoHandler{useCase: useCase}
}

func (h *EsferaFederalBuscarBlocoHandler) Buscar(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarBloco")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao buscar bloco", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar bloco: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
