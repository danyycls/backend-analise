package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/cartoes/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarCartoesHandler struct {
	useCase *usecase.BuscarCartoesUseCase
}

func NovoBuscarCartoesHandler(useCase *usecase.BuscarCartoesUseCase) *BuscarCartoesHandler {
	return &BuscarCartoesHandler{useCase: useCase}
}

func (h *BuscarCartoesHandler) Buscar(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.CartaoQueryParams{
		Pagina:              pagina,
		MesExtratoInicio:    c.Query("mesExtratoInicio"),
		MesExtratoFim:       c.Query("mesExtratoFim"),
		DataTransacaoInicio: c.Query("dataTransacaoInicio"),
		DataTransacaoFim:    c.Query("dataTransacaoFim"),
		TipoCartao:          c.Query("tipoCartao"),
		CodigoOrgao:         c.Query("codigoOrgao"),
		CPFPortador:         c.Query("cpfPortador"),
		CPFCNPJFavorecido:   c.Query("cpfCnpjFavorecido"),
		ValorDe:             c.Query("valorDe"),
		ValorAte:            c.Query("valorAte"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.Cartao{}
	}
	c.JSON(http.StatusOK, resultado)
}
