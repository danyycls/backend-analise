package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type ListarRelacionadasHandler struct {
	useCase *usecase.ListarRelacionadasUseCase
}

func NovoListarRelacionadasHandler(useCase *usecase.ListarRelacionadasUseCase) *ListarRelacionadasHandler {
	return &ListarRelacionadasHandler{useCase: useCase}
}

func (h *ListarRelacionadasHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarRelacionadas")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao listar proposicoes relacionadas", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar proposicoes relacionadas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
