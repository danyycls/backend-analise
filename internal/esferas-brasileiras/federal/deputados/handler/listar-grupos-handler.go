package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarGruposHandler struct {
	useCase *usecase.EsferaFederalListarGruposUseCase
}

func NovoEsferaFederalListarGruposHandler(useCase *usecase.EsferaFederalListarGruposUseCase) *EsferaFederalListarGruposHandler {
	return &EsferaFederalListarGruposHandler{useCase: useCase}
}

func (h *EsferaFederalListarGruposHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarGrupos")

	resp, err := h.useCase.Executar(c.Request.Context())
	if err != nil {
		log.Error("erro ao listar grupos", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar grupos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
