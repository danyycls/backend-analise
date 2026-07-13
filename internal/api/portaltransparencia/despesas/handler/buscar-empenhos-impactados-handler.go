package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarEmpenhosImpactadosHandler struct {
	useCase *usecase.BuscarEmpenhosImpactadosUseCase
}

func NovoBuscarEmpenhosImpactadosHandler(useCase *usecase.BuscarEmpenhosImpactadosUseCase) *BuscarEmpenhosImpactadosHandler {
	return &BuscarEmpenhosImpactadosHandler{useCase: useCase}
}

func (h *BuscarEmpenhosImpactadosHandler) BuscarEmpenhosImpactados(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	resultado, err := h.useCase.Buscar(c.Request.Context(), c.Query("codigoDocumento"), c.Query("fase"), pagina)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.EmpenhoImpactadoBasico{}
	}
	c.JSON(http.StatusOK, resultado)
}
