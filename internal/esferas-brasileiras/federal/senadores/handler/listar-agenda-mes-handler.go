package handler

import (
	"net/http"

	senadores "github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores"
	"github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarAgendaMesHandler struct {
	useCase *usecase.ListarAgendaMesUseCase
}

func NovoListarAgendaMesHandler(useCase *usecase.ListarAgendaMesUseCase) *ListarAgendaMesHandler {
	return &ListarAgendaMesHandler{useCase: useCase}
}

func (h *ListarAgendaMesHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	data := c.Param("data")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "data invalida"})
		return
	}

	reunioes, err := h.useCase.Listar(c.Request.Context(), data, senadores.QueryParams(c))
	if err != nil {
		log.Error("erro ao listar agenda mes", "data", data, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar agenda mes: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": reunioes})
}
