package handler

import (
	"net/http"

	senado "github.com/danyele/podp/internal/api/senado"
	"github.com/danyele/podp/internal/api/senado/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarProcessoEmendasHandler struct {
	useCase *usecase.ListarProcessoEmendasUseCase
}

func NovoListarProcessoEmendasHandler(useCase *usecase.ListarProcessoEmendasUseCase) *ListarProcessoEmendasHandler {
	return &ListarProcessoEmendasHandler{useCase: useCase}
}

func (h *ListarProcessoEmendasHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	emendas, err := h.useCase.Listar(c.Request.Context(), senado.QueryParams(c))
	if err != nil {
		log.Error("erro ao listar emendas processo", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar emendas: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": emendas})
}
