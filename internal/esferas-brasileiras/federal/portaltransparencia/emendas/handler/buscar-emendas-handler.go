package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/emendas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarEmendasHandler struct {
	useCase *usecase.BuscarEmendasUseCase
}

func NovoBuscarEmendasHandler(useCase *usecase.BuscarEmendasUseCase) *BuscarEmendasHandler {
	return &BuscarEmendasHandler{useCase: useCase}
}

func (h *BuscarEmendasHandler) Buscar(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.EmendaQueryParams{
		Pagina:          pagina,
		CodigoEmenda:    c.Query("codigoEmenda"),
		NumeroEmenda:    c.Query("numeroEmenda"),
		NomeAutor:       c.Query("nomeAutor"),
		TipoEmenda:      c.Query("tipoEmenda"),
		Ano:             c.Query("ano"),
		CodigoFuncao:    c.Query("codigoFuncao"),
		CodigoSubfuncao: c.Query("codigoSubfuncao"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.ConsultaEmendas{}
	}
	c.JSON(http.StatusOK, resultado)
}
