package handler

import (
	"net/http"

	tsetypes "github.com/danyele/podp/internal/api/tse/types"
	"github.com/danyele/podp/internal/api/tse/usecase"

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

	if len(req.Documento) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "documento nao pode ser vazio"})
		return
	}

	result, err := h.useCase.Executar(c.Request.Context(), req.Documento)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
