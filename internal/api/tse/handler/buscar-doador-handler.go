package handler

import (
	"net/http"

	tsetypes "github.com/danyele/podp/internal/api/tse/types"
	"github.com/danyele/podp/internal/api/tse/usecase"

	"github.com/gin-gonic/gin"
)

type BuscarDoadorHandler struct {
	useCase *usecase.BuscarDoadorUseCase
}

func NovoBuscarDoadorHandler(useCase *usecase.BuscarDoadorUseCase) *BuscarDoadorHandler {
	return &BuscarDoadorHandler{useCase: useCase}
}

func (h *BuscarDoadorHandler) BuscarDoador(c *gin.Context) {
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
