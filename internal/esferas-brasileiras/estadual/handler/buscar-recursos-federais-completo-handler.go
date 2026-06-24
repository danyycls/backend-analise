package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	dadosfinanceiros "github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase/dadosfinanceiros"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaEstadualBuscarRecursosFederaisCompletoHandler struct {
	useCase *dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisCompletoUseCase
}

func NovoEsferaEstadualBuscarRecursosFederaisCompletoHandler(useCase *dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisCompletoUseCase) *EsferaEstadualBuscarRecursosFederaisCompletoHandler {
	return &EsferaEstadualBuscarRecursosFederaisCompletoHandler{useCase: useCase}
}

func (h *EsferaEstadualBuscarRecursosFederaisCompletoHandler) Buscar(c *gin.Context) {
	log := logger.New("Estadual: Handler: BuscarRecursosFederaisCompleto")
	uf := c.Param("uf")
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	anoStr := c.Query("ano")
	if anoStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "parametro 'ano' obrigatorio"})
		return
	}

	ano, err := strconv.ParseInt(anoStr, 10, 64)
	if err != nil || ano < 2000 || ano > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ano invalido"})
		return
	}

	codigoIBGE := c.Query("codigoIBGE")

	req := &dadosfinanceiros.EsferaEstadualBuscarRecursosFederaisCompletoRequest{
		UF:         uf,
		Exercicio:  ano,
		CodigoIBGE: codigoIBGE,
	}

	resp, err := h.useCase.Executar(c.Request.Context(), req)
	if err != nil {
		log.Error("erro ao buscar recursos federais completos", "uf", uf, "ano", ano, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar recursos federais"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
