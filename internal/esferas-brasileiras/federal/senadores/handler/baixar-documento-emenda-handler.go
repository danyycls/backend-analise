package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type BaixarDocumentoEmendaHandler struct {
	useCase *usecase.BaixarDocumentoEmendaUseCase
}

func NovoBaixarDocumentoEmendaHandler(useCase *usecase.BaixarDocumentoEmendaUseCase) *BaixarDocumentoEmendaHandler {
	return &BaixarDocumentoEmendaHandler{useCase: useCase}
}

func (h *BaixarDocumentoEmendaHandler) Baixar(c *gin.Context) {
	log := logger.New("Senadores: Handler: BaixarDocumentoEmenda")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	dados, contentType, err := h.useCase.Executar(c.Request.Context(), id)
	if err != nil {
		log.Error("erro ao baixar documento", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao baixar documento: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, contentType, dados)
}
