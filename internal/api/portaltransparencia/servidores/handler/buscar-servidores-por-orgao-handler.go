package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/servidores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarServidoresPorOrgaoHandler struct {
	useCase *usecase.BuscarServidoresPorOrgaoUseCase
}

func NovoBuscarServidoresPorOrgaoHandler(useCase *usecase.BuscarServidoresPorOrgaoUseCase) *BuscarServidoresPorOrgaoHandler {
	return &BuscarServidoresPorOrgaoHandler{useCase: useCase}
}

func (h *BuscarServidoresPorOrgaoHandler) BuscarPorOrgao(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.ServidorPorOrgaoQueryParams{
		Pagina:         pagina,
		OrgaoLotacao:   c.Query("orgaoLotacao"),
		OrgaoExercicio: c.Query("orgaoExercicio"),
		TipoServidor:   c.Query("tipoServidor"),
		TipoVinculo:    c.Query("tipoVinculo"),
		Licenca:        c.Query("licenca"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.ServidorPorOrgao{}
	}
	c.JSON(http.StatusOK, resultado)
}
