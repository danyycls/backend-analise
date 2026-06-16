package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/servidores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarRemuneracaoHandler struct {
	useCase *usecase.BuscarRemuneracaoServidoresUseCase
}

func NovoBuscarRemuneracaoHandler(useCase *usecase.BuscarRemuneracaoServidoresUseCase) *BuscarRemuneracaoHandler {
	return &BuscarRemuneracaoHandler{useCase: useCase}
}

func (h *BuscarRemuneracaoHandler) BuscarRemuneracao(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.ServidorRemuneracaoQueryParams{
		Pagina: pagina,
		CPF:    c.Query("cpf"),
		MesAno: c.Query("mesAno"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.ServidorRemuneracao{}
	}
	c.JSON(http.StatusOK, resultado)
}
