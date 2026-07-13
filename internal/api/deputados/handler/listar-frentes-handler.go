package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarFrentesHandler struct {
	useCase *usecase.EsferaFederalListarFrentesUseCase
}

func NovoEsferaFederalListarFrentesHandler(useCase *usecase.EsferaFederalListarFrentesUseCase) *EsferaFederalListarFrentesHandler {
	return &EsferaFederalListarFrentesHandler{useCase: useCase}
}

func (h *EsferaFederalListarFrentesHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarFrentes")
	idLegislaturaStr := c.Query("idLegislatura")
	idLegislatura, err := strconv.Atoi(idLegislaturaStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "idLegislatura invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), idLegislatura)
	if err != nil {
		log.Error("erro ao listar frentes", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar frentes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
