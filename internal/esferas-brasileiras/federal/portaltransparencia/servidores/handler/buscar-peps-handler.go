package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/servidores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarPEPsHandler struct {
	useCase *usecase.BuscarPEPsUseCase
}

func NovoBuscarPEPsHandler(useCase *usecase.BuscarPEPsUseCase) *BuscarPEPsHandler {
	return &BuscarPEPsHandler{useCase: useCase}
}

func (h *BuscarPEPsHandler) BuscarPEPs(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	filtro := portaltransparencia.PEPQueryParams{
		Pagina:                 pagina,
		CPF:                    c.Query("cpf"),
		Nome:                   c.Query("nome"),
		DescricaoFuncao:        c.Query("descricaoFuncao"),
		OrgaoServidorLotacao:   c.Query("orgaoServidorLotacao"),
		DataInicioExercicioDe:  c.Query("dataInicioExercicioDe"),
		DataInicioExercicioAte: c.Query("dataInicioExercicioAte"),
		DataFimExercicioDe:     c.Query("dataFimExercicioDe"),
		DataFimExercicioAte:    c.Query("dataFimExercicioAte"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.PEP{}
	}
	c.JSON(http.StatusOK, resultado)
}
