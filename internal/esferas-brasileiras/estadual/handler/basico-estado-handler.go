package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaEstadualBuscarDadosBasicosEstadoHandler struct {
	useCase *usecase.EsferaEstadualBuscarDadosBasicosEstadoUseCase
}

func NovoEsferaEstadualBuscarDadosBasicosEstadoHandler(useCase *usecase.EsferaEstadualBuscarDadosBasicosEstadoUseCase) *EsferaEstadualBuscarDadosBasicosEstadoHandler {
	return &EsferaEstadualBuscarDadosBasicosEstadoHandler{useCase: useCase}
}

func (h *EsferaEstadualBuscarDadosBasicosEstadoHandler) BuscarBasicoEstado(c *gin.Context) {
	log := logger.New("Estadual: Handler: BuscarBasicoEstado")
	uf := c.Param("uf")
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaEstadualBuscarDadosBasicosEstadoRequest{UF: uf})
	if err != nil {
		log.Error("erro ao buscar basico estado", "uf", uf, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar dados basicos"})
		return
	}

	c.JSON(http.StatusOK, resp.Dados)
}
