package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaEstadualListarEstadosHandler struct {
	useCase *usecase.EsferaEstadualListarEstadosUseCase
}

func NovoEsferaEstadualListarEstadosHandler(useCase *usecase.EsferaEstadualListarEstadosUseCase) *EsferaEstadualListarEstadosHandler {
	return &EsferaEstadualListarEstadosHandler{useCase: useCase}
}

func (h *EsferaEstadualListarEstadosHandler) ListarEstados(c *gin.Context) {
	log := logger.New("Estadual: Handler: ListarEstados")
	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaEstadualListarEstadosRequest{})
	if err != nil {
		log.Error("erro ao listar estados", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar estados"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": resp.Estados})
}
