package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/ligacao-politica/usecase"
	"github.com/danyele/podp/internal/shared/logger"
	redis "github.com/danyele/podp/internal/shared/redis"
)

type AnalisarLigacaoPoliticaHandler struct {
	useCase usecase.AnalisarLigacaoPoliticaUseCaseInterface
	redis   redis.Cache
}

func NovoAnalisarLigacaoPoliticaHandler(useCase usecase.AnalisarLigacaoPoliticaUseCaseInterface, redis redis.Cache) *AnalisarLigacaoPoliticaHandler {
	return &AnalisarLigacaoPoliticaHandler{useCase: useCase, redis: redis}
}

func (h *AnalisarLigacaoPoliticaHandler) Analisar(c *gin.Context) {
	log := logger.New("LigacaoPolitica: Handler: Analisar")
	var req struct {
		Licitacoes []usecase.AnalisarLigacaoPoliticaRequest `json:"licitacoes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "corpo inválido: " + err.Error()})
		return
	}

	if len(req.Licitacoes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "array licitacoes é obrigatório"})
		return
	}

	raw, _ := json.Marshal(req.Licitacoes)
	chave := redis.ChaveCache("ligacao-politica", raw)

	var cached usecase.AnalisarLigacaoPoliticaResponse
	ok, err := h.redis.Get(c.Request.Context(), chave, &cached)
	if err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}
	if ok {
		c.JSON(http.StatusOK, &cached)
		return
	}

	resultado, err := h.useCase.Executar(c.Request.Context(), req.Licitacoes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	if err := h.redis.Set(c.Request.Context(), chave, resultado); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}

	c.JSON(http.StatusOK, resultado)
}
