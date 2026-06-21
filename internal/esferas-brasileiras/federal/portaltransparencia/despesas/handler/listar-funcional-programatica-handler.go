package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
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
