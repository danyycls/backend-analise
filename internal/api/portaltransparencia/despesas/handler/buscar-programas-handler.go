package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarProgramasHandler struct {
	useCase *usecase.BuscarProgramasUseCase
}

func NovoBuscarProgramasHandler(useCase *usecase.BuscarProgramasUseCase) *BuscarProgramasHandler {
	return &BuscarProgramasHandler{useCase: useCase}
}

func (h *BuscarProgramasHandler) BuscarProgramas(c *gin.Context) {
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
		resultado = []portaltransparencia.CodigoDescricao{}
	}
	c.JSON(http.StatusOK, resultado)
}
