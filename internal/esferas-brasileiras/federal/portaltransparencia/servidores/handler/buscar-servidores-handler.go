package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/portaltransparencia/servidores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarServidoresHandler struct {
	useCase *usecase.BuscarServidoresUseCase
}

func NovoBuscarServidoresHandler(useCase *usecase.BuscarServidoresUseCase) *BuscarServidoresHandler {
	return &BuscarServidoresHandler{useCase: useCase}
}

func (h *BuscarServidoresHandler) Buscar(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.ServidorQueryParams{
		Pagina:                 pagina,
		CPF:                    c.Query("cpf"),
		Nome:                   c.Query("nome"),
		OrgaoServidorLotacao:   c.Query("orgaoServidorLotacao"),
		OrgaoServidorExercicio: c.Query("orgaoServidorExercicio"),
		SituacaoServidor:       c.Query("situacaoServidor"),
		TipoServidor:           c.Query("tipoServidor"),
		CodigoFuncaoCargo:      c.Query("codigoFuncaoCargo"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.CadastroServidor{}
	}
	c.JSON(http.StatusOK, resultado)
}
