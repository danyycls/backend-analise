package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/pessoas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarJuridicaHandler struct {
	useCase *usecase.BuscarPessoasJuridicasUseCase
}

func NovoBuscarJuridicaHandler(useCase *usecase.BuscarPessoasJuridicasUseCase) *BuscarJuridicaHandler {
	return &BuscarJuridicaHandler{useCase: useCase}
}

func (h *BuscarJuridicaHandler) BuscarJuridica(c *gin.Context) {
	filtro := portaltransparencia.PessoaJuridicaQueryParams{
		CNPJ:                      c.Query("cnpj"),
		RazaoSocial:               c.Query("razaoSocial"),
		NomeFantasia:              c.Query("nomeFantasia"),
		FavorecidoDespesas:        boolQuery(c, "favorecidoDespesas"),
		PossuiContratacao:         boolQuery(c, "possuiContratacao"),
		Convenios:                 boolQuery(c, "convenios"),
		FavorecidoTransferencias:  boolQuery(c, "favorecidoTransferencias"),
		SancionadoCEPIM:           boolQuery(c, "sancionadoCEPIM"),
		SancionadoCEIS:            boolQuery(c, "sancionadoCEIS"),
		SancionadoCNEP:            boolQuery(c, "sancionadoCNEP"),
		SancionadoCEAF:            boolQuery(c, "sancionadoCEAF"),
		SancionadoAcordoLeniencia: boolQuery(c, "sancionadoAcordoLeniencia"),
		Ordenacao:                 c.Query("ordenacao"),
		OrdenacaoDirecao:          c.Query("ordenacaoDirecao"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}
