package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarOrgaosHandler struct {
	useCase *usecase.EsferaFederalListarOrgaosUseCase
}

func NovoEsferaFederalListarOrgaosHandler(useCase *usecase.EsferaFederalListarOrgaosUseCase) *EsferaFederalListarOrgaosHandler {
	return &EsferaFederalListarOrgaosHandler{useCase: useCase}
}

func (h *EsferaFederalListarOrgaosHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarOrgaos")
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resp, err := h.useCase.Executar(c.Request.Context(), params)
	if err != nil {
		log.Error("erro ao listar orgaos", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar orgaos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
