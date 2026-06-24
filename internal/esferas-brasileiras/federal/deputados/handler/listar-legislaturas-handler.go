package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarLegislaturasHandler struct {
	useCase *usecase.EsferaFederalListarLegislaturasUseCase
}

func NovoEsferaFederalListarLegislaturasHandler(useCase *usecase.EsferaFederalListarLegislaturasUseCase) *EsferaFederalListarLegislaturasHandler {
	return &EsferaFederalListarLegislaturasHandler{useCase: useCase}
}

func (h *EsferaFederalListarLegislaturasHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarLegislaturas")

	resp, err := h.useCase.Executar(c.Request.Context())
	if err != nil {
		log.Error("erro ao listar legislaturas", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar legislaturas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
