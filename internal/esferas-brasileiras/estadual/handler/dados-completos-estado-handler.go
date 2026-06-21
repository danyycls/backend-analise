package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaEstadualBuscarDadosCompletosEstadoHandler struct {
	useCase *usecase.EsferaEstadualBuscarDadosCompletosEstadoUseCase
}

func NovoEsferaEstadualBuscarDadosCompletosEstadoHandler(useCase *usecase.EsferaEstadualBuscarDadosCompletosEstadoUseCase) *EsferaEstadualBuscarDadosCompletosEstadoHandler {
	return &EsferaEstadualBuscarDadosCompletosEstadoHandler{useCase: useCase}
}

func (h *EsferaEstadualBuscarDadosCompletosEstadoHandler) BuscarDadosEstado(c *gin.Context) {
	log := logger.New("Estadual: Handler: BuscarDadosEstado")
	uf := c.Param("uf")
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaEstadualBuscarDadosCompletosEstadoRequest{UF: uf})
	if err != nil {
		log.Error("erro ao buscar dados estado", "uf", uf, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar dados do estado"})
		return
	}

	c.JSON(http.StatusOK, resp.Dados)
}
