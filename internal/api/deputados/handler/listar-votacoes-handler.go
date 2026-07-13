package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarVotacoesHandler struct {
	useCase *usecase.EsferaFederalListarVotacoesUseCase
}

func NovoEsferaFederalListarVotacoesHandler(useCase *usecase.EsferaFederalListarVotacoesUseCase) *EsferaFederalListarVotacoesHandler {
	return &EsferaFederalListarVotacoesHandler{useCase: useCase}
}

func (h *EsferaFederalListarVotacoesHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarVotacoes")
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resp, err := h.useCase.Executar(c.Request.Context(), params)
	if err != nil {
		log.Error("erro ao listar votacoes", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar votacoes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
