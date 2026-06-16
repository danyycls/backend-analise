package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/logger"
)

type BuscarComissaoHandler struct {
	useCase *usecase.BuscarComissaoUseCase
}

func NovoBuscarComissaoHandler(useCase *usecase.BuscarComissaoUseCase) *BuscarComissaoHandler {
	return &BuscarComissaoHandler{useCase: useCase}
}

func (h *BuscarComissaoHandler) Buscar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Buscar")
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo invalido"})
		return
	}

	comissao, err := h.useCase.Buscar(c.Request.Context(), codigo)
	if err != nil {
		log.Error("erro ao buscar comissao", "codigo", codigo, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar comissao: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": comissao})
}
