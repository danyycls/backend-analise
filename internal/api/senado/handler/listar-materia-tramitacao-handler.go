package handler

import (
	"net/http"

	senado "github.com/danyele/podp/internal/api/senado"
	"github.com/danyele/podp/internal/api/senado/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
)

type ListarMateriaTramitacaoHandler struct {
	useCase *usecase.ListarMateriaTramitacaoUseCase
}

func NovoListarMateriaTramitacaoHandler(useCase *usecase.ListarMateriaTramitacaoUseCase) *ListarMateriaTramitacaoHandler {
	return &ListarMateriaTramitacaoHandler{useCase: useCase}
}

func (h *ListarMateriaTramitacaoHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	materials, err := h.useCase.Listar(c.Request.Context(), senado.QueryParams(c))
	if err != nil {
		log.Error("erro ao listar materia tramitacao", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar materia tramitacao: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dados": materials})
}
