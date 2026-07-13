package handler

import (
	"net/http"

	tsetypes "github.com/danyele/podp/internal/api/tse/types"
	"github.com/danyele/podp/internal/api/tse/usecase"

	"github.com/gin-gonic/gin"
)

type ConsultaEntidadeHandler struct {
	useCase *usecase.ConsultarEntidadeUseCase
}

func NovoConsultarEntidadeHandler(useCase *usecase.ConsultarEntidadeUseCase) *ConsultaEntidadeHandler {
	return &ConsultaEntidadeHandler{useCase: useCase}
}

func (h *ConsultaEntidadeHandler) Consultar(c *gin.Context) {
	var req tsetypes.ConsultaEntidadeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	if req.Tipo == "" || req.Chave == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "tipo e chave sao obrigatorios"})
		return
	}

	result := h.useCase.Executar(c.Request.Context(), &req)
	c.JSON(http.StatusOK, result)
}
