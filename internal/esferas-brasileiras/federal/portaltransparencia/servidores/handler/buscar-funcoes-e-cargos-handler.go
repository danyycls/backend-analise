package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/servidores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarFuncoesECargosHandler struct {
	useCase *usecase.BuscarFuncoesECargosUseCase
}

func NovoBuscarFuncoesECargosHandler(useCase *usecase.BuscarFuncoesECargosUseCase) *BuscarFuncoesECargosHandler {
	return &BuscarFuncoesECargosHandler{useCase: useCase}
}

func (h *BuscarFuncoesECargosHandler) BuscarFuncoesECargos(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.FuncaoCargoQueryParams{
		Pagina:               pagina,
		CodigoFuncaoCargo:    c.Query("codigoFuncaoCargo"),
		DescricaoFuncaoCargo: c.Query("descricaoFuncaoCargo"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.FuncaoServidor{}
	}
	c.JSON(http.StatusOK, resultado)
}
