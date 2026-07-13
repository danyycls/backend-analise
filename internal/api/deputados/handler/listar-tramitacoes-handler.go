package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type ListarTramitacoesHandler struct {
	useCase *usecase.ListarTramitacoesUseCase
}

func NovoListarTramitacoesHandler(useCase *usecase.ListarTramitacoesUseCase) *ListarTramitacoesHandler {
	return &ListarTramitacoesHandler{useCase: useCase}
}

func (h *ListarTramitacoesHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarTramitacoes")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao listar tramitacoes", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar tramitacoes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
