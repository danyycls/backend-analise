package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaEstadualBuscarMunicipiosPopulacaoHandler struct {
	useCase *usecase.EsferaEstadualBuscarMunicipiosPopulacaoUseCase
}

func NovoEsferaEstadualBuscarMunicipiosPopulacaoHandler(useCase *usecase.EsferaEstadualBuscarMunicipiosPopulacaoUseCase) *EsferaEstadualBuscarMunicipiosPopulacaoHandler {
	return &EsferaEstadualBuscarMunicipiosPopulacaoHandler{useCase: useCase}
}

func (h *EsferaEstadualBuscarMunicipiosPopulacaoHandler) BuscarMunicipiosPopulacao(c *gin.Context) {
	log := logger.New("Estadual: Handler: BuscarMunicipiosPopulacao")
	uf := c.Param("uf")
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaEstadualBuscarMunicipiosPopulacaoRequest{UF: uf})
	if err != nil {
		log.Error("erro ao buscar municipios/populacao", "uf", uf, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar municipios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp.Municipios})
}
