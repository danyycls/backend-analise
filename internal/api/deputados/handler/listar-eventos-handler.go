package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarEventosHandler struct {
	useCase *usecase.EsferaFederalListarEventosUseCase
}

func NovoEsferaFederalListarEventosHandler(useCase *usecase.EsferaFederalListarEventosUseCase) *EsferaFederalListarEventosHandler {
	return &EsferaFederalListarEventosHandler{useCase: useCase}
}

func (h *EsferaFederalListarEventosHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarEventos")
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resp, err := h.useCase.Executar(c.Request.Context(), params)
	if err != nil {
		log.Error("erro ao listar eventos", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar eventos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
