package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalBuscarOrgaoAssociadoDeputadoHandler struct {
	useCase *usecase.EsferaFederalBuscarOrgaoAssociadoDeputadoUseCase
}

func NovoEsferaFederalBuscarOrgaoAssociadoDeputadoHandler(useCase *usecase.EsferaFederalBuscarOrgaoAssociadoDeputadoUseCase) *EsferaFederalBuscarOrgaoAssociadoDeputadoHandler {
	return &EsferaFederalBuscarOrgaoAssociadoDeputadoHandler{useCase: useCase}
}

func (h *EsferaFederalBuscarOrgaoAssociadoDeputadoHandler) BuscarOrgaoAssociadoDeputado(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarOrgaoAssociadoDeputado")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaFederalBuscarOrgaoAssociadoDeputadoRequest{
		ID:     id,
		Params: params,
	})
	if err != nil {
		log.Error("erro ao buscar orgaos deputado", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar orgaos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp.Orgaos})
}
