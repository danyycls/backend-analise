package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarRecursosRecebidosHandler struct {
	useCase *usecase.BuscarRecursosRecebidosUseCase
}

func NovoBuscarRecursosRecebidosHandler(useCase *usecase.BuscarRecursosRecebidosUseCase) *BuscarRecursosRecebidosHandler {
	return &BuscarRecursosRecebidosHandler{useCase: useCase}
}

func (h *BuscarRecursosRecebidosHandler) BuscarRecursosRecebidos(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.DespesaRecursosRecebidosQueryParams{
		Pagina:           pagina,
		MesAnoInicio:     c.Query("mesAnoInicio"),
		MesAnoFim:        c.Query("mesAnoFim"),
		NomeFavorecido:   c.Query("nomeFavorecido"),
		CodigoFavorecido: c.Query("codigoFavorecido"),
		TipoFavorecido:   c.Query("tipoFavorecido"),
		UF:               c.Query("uf"),
		CodigoIBGE:       c.Query("codigoIBGE"),
		OrgaoSuperior:    c.Query("orgaoSuperior"),
		Orgao:            c.Query("orgao"),
		UnidadeGestora:   c.Query("unidadeGestora"),
	}
	if filtro.MesAnoInicio == "" || filtro.MesAnoFim == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "É necessário informar 'mesAnoInicio' e 'mesAnoFim' (formato YYYY-MM)"})
		return
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.PessoaRecursosRecebidosUGMesDesnormalizada{}
	}
	c.JSON(http.StatusOK, resultado)
}
