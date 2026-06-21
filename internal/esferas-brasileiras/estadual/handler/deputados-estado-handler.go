package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaEstadualBuscarDeputadosHandler struct {
	useCase *usecase.EsferaEstadualBuscarDeputadosUseCase
}

func NovoEsferaEstadualBuscarDeputadosHandler(useCase *usecase.EsferaEstadualBuscarDeputadosUseCase) *EsferaEstadualBuscarDeputadosHandler {
	return &EsferaEstadualBuscarDeputadosHandler{useCase: useCase}
}

func (h *EsferaEstadualBuscarDeputadosHandler) BuscarDeputadosEstado(c *gin.Context) {
	log := logger.New("Estadual: Handler: BuscarDeputadosEstado")
	uf := c.Param("uf")
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaEstadualBuscarDeputadosRequest{UF: uf})
	if err != nil {
		log.Error("erro ao buscar deputados UF", "uf", uf, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar deputados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp.Deputados})
}
