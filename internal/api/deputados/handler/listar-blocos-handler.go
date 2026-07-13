package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarBlocosHandler struct {
	useCase *usecase.EsferaFederalListarBlocosUseCase
}

func NovoEsferaFederalListarBlocosHandler(useCase *usecase.EsferaFederalListarBlocosUseCase) *EsferaFederalListarBlocosHandler {
	return &EsferaFederalListarBlocosHandler{useCase: useCase}
}

func (h *EsferaFederalListarBlocosHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarBlocos")
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resp, err := h.useCase.Executar(c.Request.Context(), params)
	if err != nil {
		log.Error("erro ao listar blocos", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar blocos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
