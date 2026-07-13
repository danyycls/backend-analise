package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/ibge/usecase"
)

type BuscarPopulacaoHandler struct {
	useCase *usecase.BuscarPopulacaoUseCase
}

func NovoBuscarPopulacaoHandler(useCase *usecase.BuscarPopulacaoUseCase) *BuscarPopulacaoHandler {
	return &BuscarPopulacaoHandler{useCase: useCase}
}

func (h *BuscarPopulacaoHandler) BuscarPopulacao(c *gin.Context) {
	var req usecase.BuscarPopulacaoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if len(req.MunicipioIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "municipio_ids e obrigatorio"})
		return
	}

	populacao, err := h.useCase.Executar(c.Request.Context(), req.MunicipioIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, populacao)
}
