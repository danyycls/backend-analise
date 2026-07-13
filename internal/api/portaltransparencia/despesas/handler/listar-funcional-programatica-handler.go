package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type ListarFuncionalProgramaticaHandler struct {
	useCase *usecase.ListarFuncionalProgramaticaUseCase
}

func NovoListarFuncionalProgramaticaHandler(useCase *usecase.ListarFuncionalProgramaticaUseCase) *ListarFuncionalProgramaticaHandler {
	return &ListarFuncionalProgramaticaHandler{useCase: useCase}
}

func (h *ListarFuncionalProgramaticaHandler) ListarFuncionalProgramatica(c *gin.Context) {
	ano, _ := strconv.Atoi(c.Query("ano"))
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	resultado, err := h.useCase.Buscar(c.Request.Context(), ano, pagina)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.FuncionalProgramatica{}
	}
	c.JSON(http.StatusOK, resultado)
}
