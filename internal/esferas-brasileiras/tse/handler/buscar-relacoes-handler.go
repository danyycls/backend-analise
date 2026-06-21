package handler

import (
	"net/http"

	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/esferas-brasileiras/tse/usecase"

	"github.com/gin-gonic/gin"
)

type BuscaRelacoesHandler struct {
	useCase *usecase.BuscarRelacoesUseCase
}

func NovoBuscarRelacoesHandler(useCase *usecase.BuscarRelacoesUseCase) *BuscaRelacoesHandler {
	return &BuscaRelacoesHandler{useCase: useCase}
}

func (h *BuscaRelacoesHandler) Buscar(c *gin.Context) {
	var req tsetypes.BuscaRelacoesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if len(req.CNPJ) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "CNPJ nao pode ser vazio"})
		return
	}

	result, err := h.useCase.Executar(c.Request.Context(), req.CNPJ)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
