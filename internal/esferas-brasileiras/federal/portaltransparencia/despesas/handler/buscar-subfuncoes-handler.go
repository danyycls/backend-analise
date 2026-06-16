package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarSubfuncoesHandler struct {
	useCase *usecase.BuscarSubfuncoesUseCase
}

func NovoBuscarSubfuncoesHandler(useCase *usecase.BuscarSubfuncoesUseCase) *BuscarSubfuncoesHandler {
	return &BuscarSubfuncoesHandler{useCase: useCase}
}

func (h *BuscarSubfuncoesHandler) BuscarSubfuncoes(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	anoInicio, _ := strconv.Atoi(c.DefaultQuery("anoInicio", "0"))
	filtro := portaltransparencia.ListarFuncionalProgramaticaQueryParams{
		AnoInicio: anoInicio,
		Pagina:    pagina,
		Codigo:    c.Query("codigo"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.Subfuncao{}
	}
	c.JSON(http.StatusOK, resultado)
}
