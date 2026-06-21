package handler

import (
	"net/http"

	senadores "github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores"
	"github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type BuscarEncontroHandler struct {
	useCase *usecase.BuscarEncontroUseCase
}

func NovoBuscarEncontroHandler(useCase *usecase.BuscarEncontroUseCase) *BuscarEncontroHandler {
	return &BuscarEncontroHandler{useCase: useCase}
}

func (h *BuscarEncontroHandler) Buscar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Buscar")
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo invalido"})
		return
	}

	encontro, err := h.useCase.Buscar(c.Request.Context(), codigo, senadores.QueryParams(c))
	if err != nil {
		log.Error("erro ao buscar encontro", "codigo", codigo, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar encontro: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, encontro)
}
