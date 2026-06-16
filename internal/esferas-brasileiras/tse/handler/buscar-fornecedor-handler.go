package handler

import (
	"net/http"

	tsetypes "github.com/danyele/laceu/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/laceu/internal/esferas-brasileiras/tse/usecase"

	"github.com/gin-gonic/gin"
)

type BuscarFornecedorHandler struct {
	useCase *usecase.BuscarFornecedorUseCase
}

func NovoBuscarFornecedorHandler(useCase *usecase.BuscarFornecedorUseCase) *BuscarFornecedorHandler {
	return &BuscarFornecedorHandler{useCase: useCase}
}

func (h *BuscarFornecedorHandler) BuscarFornecedor(c *gin.Context) {
	var req tsetypes.BuscaDocumentoRequest
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
