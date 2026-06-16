package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarDespesasPorOrgaoHandler struct {
	useCase *usecase.BuscarDespesasPorOrgaoUseCase
}

func NovoBuscarDespesasPorOrgaoHandler(useCase *usecase.BuscarDespesasPorOrgaoUseCase) *BuscarDespesasPorOrgaoHandler {
	return &BuscarDespesasPorOrgaoHandler{useCase: useCase}
}

func (h *BuscarDespesasPorOrgaoHandler) BuscarPorOrgao(c *gin.Context) {
	filtro := portaltransparencia.DespesaPorOrgaoQueryParams{
		Ano:           c.Query("ano"),
		Pagina:        func() int { p, _ := strconv.Atoi(c.DefaultQuery("pagina", "1")); return p }(),
		OrgaoSuperior: c.Query("orgaoSuperior"),
		Orgao:         c.Query("orgao"),
	}
	if filtro.Ano == "" || (filtro.OrgaoSuperior == "" && filtro.Orgao == "") {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "É necessário informar 'ano' e ao menos um dos filtros: 'orgaoSuperior' ou 'orgao'"})
		return
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.DespesaAnualPorOrgao{}
	}
	c.JSON(http.StatusOK, resultado)
}
