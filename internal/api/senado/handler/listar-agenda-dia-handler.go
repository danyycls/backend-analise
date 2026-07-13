package handler

import (
	"net/http"

	senado "github.com/danyele/podp/internal/api/senado"
	"github.com/danyele/podp/internal/api/senado/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarAgendaDiaHandler struct {
	useCase *usecase.ListarAgendaDiaUseCase
}

func NovoListarAgendaDiaHandler(useCase *usecase.ListarAgendaDiaUseCase) *ListarAgendaDiaHandler {
	return &ListarAgendaDiaHandler{useCase: useCase}
}

func (h *ListarAgendaDiaHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	data := c.Param("data")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "data invalida"})
		return
	}

	reunioes, err := h.useCase.Listar(c.Request.Context(), data, senado.QueryParams(c))
	if err != nil {
		log.Error("erro ao listar agenda dia", "data", data, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar agenda dia: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": reunioes})
}
