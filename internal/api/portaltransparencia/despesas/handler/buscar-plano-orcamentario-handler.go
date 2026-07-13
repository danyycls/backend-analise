package handler

import (
	"net/http"
	"strconv"

	"github.com/danyele/podp/internal/api/portaltransparencia/despesas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarPlanoOrcamentarioHandler struct {
	useCase *usecase.BuscarPlanoOrcamentarioUseCase
}

func NovoBuscarPlanoOrcamentarioHandler(useCase *usecase.BuscarPlanoOrcamentarioUseCase) *BuscarPlanoOrcamentarioHandler {
	return &BuscarPlanoOrcamentarioHandler{useCase: useCase}
}

func (h *BuscarPlanoOrcamentarioHandler) BuscarPlanoOrcamentario(c *gin.Context) {
	filtro := portaltransparencia.DespesaPlanoOrcamentarioQueryParams{
		Ano:                       c.Query("ano"),
		Pagina:                    func() int { p, _ := strconv.Atoi(c.DefaultQuery("pagina", "1")); return p }(),
		CodPlanoOrcamentario:      c.Query("codPlanoOrcamentario"),
		DescPlanoOrcamentario:     c.Query("descPlanoOrcamentario"),
		CodPOIdentfAcompanhamento: c.Query("codPOIdentfAcompanhamento"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	if resultado == nil {
		resultado = []portaltransparencia.DespesasPorPlanoOrcamentario{}
	}
	c.JSON(http.StatusOK, resultado)
}
