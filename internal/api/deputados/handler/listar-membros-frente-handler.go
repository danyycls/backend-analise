package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarMembrosFrenteHandler struct {
	useCase *usecase.EsferaFederalListarMembrosFrenteUseCase
}

func NovoEsferaFederalListarMembrosFrenteHandler(useCase *usecase.EsferaFederalListarMembrosFrenteUseCase) *EsferaFederalListarMembrosFrenteHandler {
	return &EsferaFederalListarMembrosFrenteHandler{useCase: useCase}
}

func (h *EsferaFederalListarMembrosFrenteHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarMembrosFrente")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao listar membros da frente", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar membros da frente: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
