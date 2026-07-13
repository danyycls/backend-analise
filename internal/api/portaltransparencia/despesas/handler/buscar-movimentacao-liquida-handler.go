package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarMovimentacaoLiquidaHandler struct {
	useCase *usecase.BuscarMovimentacaoLiquidaUseCase
}

func NovoBuscarMovimentacaoLiquidaHandler(useCase *usecase.BuscarMovimentacaoLiquidaUseCase) *BuscarMovimentacaoLiquidaHandler {
	return &BuscarMovimentacaoLiquidaHandler{useCase: useCase}
}

func (h *BuscarMovimentacaoLiquidaHandler) BuscarMovimentacaoLiquida(c *gin.Context) {
	filtro := portaltransparencia.DespesaMovimentacaoLiquidaQueryParams{
		Ano:                 c.Query("ano"),
		Pagina:              func() int { p, _ := strconv.Atoi(c.DefaultQuery("pagina", "1")); return p }(),
		Funcao:              c.Query("funcao"),
		Subfuncao:           c.Query("subfuncao"),
		Programa:            c.Query("programa"),
		Acao:                c.Query("acao"),
		GrupoDespesa:        c.Query("grupoDespesa"),
		ElementoDespesa:     c.Query("elementoDespesa"),
		ModalidadeAplicacao: c.Query("modalidadeAplicacao"),
		IDPlanoOrcamentario: c.Query("idPlanoOrcamentario"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.DespesaLiquidaAnualPorFuncaoESubfuncao{}
	}
	c.JSON(http.StatusOK, resultado)
}
