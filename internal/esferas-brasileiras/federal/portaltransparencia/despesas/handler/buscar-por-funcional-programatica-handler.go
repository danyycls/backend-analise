package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarPorFuncionalProgramaticaHandler struct {
	useCase *usecase.BuscarDespesasPorFuncionalProgramaticaUseCase
}

func NovoBuscarPorFuncionalProgramaticaHandler(useCase *usecase.BuscarDespesasPorFuncionalProgramaticaUseCase) *BuscarPorFuncionalProgramaticaHandler {
	return &BuscarPorFuncionalProgramaticaHandler{useCase: useCase}
}

func (h *BuscarPorFuncionalProgramaticaHandler) BuscarPorFuncionalProgramatica(c *gin.Context) {
	filtro := portaltransparencia.DespesaFuncionalProgramaticaQueryParams{
		Ano:       c.Query("ano"),
		Pagina:    func() int { p, _ := strconv.Atoi(c.DefaultQuery("pagina", "1")); return p }(),
		Funcao:    c.Query("funcao"),
		Subfuncao: c.Query("subfuncao"),
		Programa:  c.Query("programa"),
		Acao:      c.Query("acao"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.DespesaAnualPorFuncaoESubfuncao{}
	}
	c.JSON(http.StatusOK, resultado)
}
