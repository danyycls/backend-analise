package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalBuscarDespesasDeputadoHandler struct {
	useCase *usecase.EsferaFederalBuscarDespesasDeputadoUseCase
}

func NovoEsferaFederalBuscarDespesasDeputadoHandler(useCase *usecase.EsferaFederalBuscarDespesasDeputadoUseCase) *EsferaFederalBuscarDespesasDeputadoHandler {
	return &EsferaFederalBuscarDespesasDeputadoHandler{useCase: useCase}
}

func (h *EsferaFederalBuscarDespesasDeputadoHandler) BuscarDespesasDeputado(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarDespesasDeputado")
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

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaFederalBuscarDespesasDeputadoRequest{
		ID:                id,
		Params:            params,
		TipoDespesa:       params["tipoDespesa"],
		CNPJCPFFornecedor: params["cnpjCpfFornecedor"],
	})
	if err != nil {
		log.Error("erro ao buscar despesas deputado", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar despesas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp.Despesas})
}
