package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/esferas-brasileiras/estadual/usecase"
	"github.com/danyele/laceu/internal/shared/logger"
)

type EsferaEstadualBuscarCandidatosHandler struct {
	useCase *usecase.EsferaEstadualBuscarCandidatosUseCase
}

func NovoEsferaEstadualBuscarCandidatosHandler(useCase *usecase.EsferaEstadualBuscarCandidatosUseCase) *EsferaEstadualBuscarCandidatosHandler {
	return &EsferaEstadualBuscarCandidatosHandler{useCase: useCase}
}

func (h *EsferaEstadualBuscarCandidatosHandler) BuscarCandidatosEstado(c *gin.Context) {
	log := logger.New("Estadual: Handler: BuscarCandidatosEstado")
	uf := c.Param("uf")
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaEstadualBuscarCandidatosRequest{UF: uf})
	if err != nil {
		log.Error("erro ao buscar candidatos UF", "uf", uf, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar candidatos"})
		return
	}

	c.JSON(http.StatusOK, resp.Dados)
}
