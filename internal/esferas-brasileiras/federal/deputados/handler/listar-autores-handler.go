package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type ListarAutoresHandler struct {
	useCase *usecase.ListarAutoresUseCase
}

func NovoListarAutoresHandler(useCase *usecase.ListarAutoresUseCase) *ListarAutoresHandler {
	return &ListarAutoresHandler{useCase: useCase}
}

func (h *ListarAutoresHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarAutores")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao listar autores", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar autores: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
