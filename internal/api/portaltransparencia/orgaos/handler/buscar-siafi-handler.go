package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/orgaos/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarSIAFIHandler struct {
	useCase *usecase.BuscarOrgaosSIAFIUseCase
}

func NovoBuscarSIAFIHandler(useCase *usecase.BuscarOrgaosSIAFIUseCase) *BuscarSIAFIHandler {
	return &BuscarSIAFIHandler{useCase: useCase}
}

func (h *BuscarSIAFIHandler) BuscarSIAFI(c *gin.Context) {
	filtro := portaltransparencia.OrgaoQueryParams{
		Pagina:    func() int { p, _ := strconv.Atoi(c.DefaultQuery("pagina", "1")); return p }(),
		Codigo:    c.Query("codigoOrgao"),
		Descricao: c.Query("nomeOrgao"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.OrgaoSIAFI{}
	}
	c.JSON(http.StatusOK, resultado)
}
