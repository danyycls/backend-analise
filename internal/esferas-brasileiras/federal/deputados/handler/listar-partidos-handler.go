package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarPartidosHandler struct {
	useCase *usecase.EsferaFederalListarPartidosUseCase
}

func NovoEsferaFederalListarPartidosHandler(useCase *usecase.EsferaFederalListarPartidosUseCase) *EsferaFederalListarPartidosHandler {
	return &EsferaFederalListarPartidosHandler{useCase: useCase}
}

func (h *EsferaFederalListarPartidosHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarPartidos")
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resp, err := h.useCase.Executar(c.Request.Context(), params)
	if err != nil {
		log.Error("erro ao listar partidos", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar partidos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
