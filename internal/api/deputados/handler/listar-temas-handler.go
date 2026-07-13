package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type ListarTemasHandler struct {
	useCase *usecase.ListarTemasUseCase
}

func NovoListarTemasHandler(useCase *usecase.ListarTemasUseCase) *ListarTemasHandler {
	return &ListarTemasHandler{useCase: useCase}
}

func (h *ListarTemasHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarTemas")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao listar temas", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar temas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
