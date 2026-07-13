package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalBuscarOrgaoHandler struct {
	useCase *usecase.EsferaFederalBuscarOrgaoUseCase
}

func NovoEsferaFederalBuscarOrgaoHandler(useCase *usecase.EsferaFederalBuscarOrgaoUseCase) *EsferaFederalBuscarOrgaoHandler {
	return &EsferaFederalBuscarOrgaoHandler{useCase: useCase}
}

func (h *EsferaFederalBuscarOrgaoHandler) Buscar(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarOrgao")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao buscar orgao", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar orgao: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
