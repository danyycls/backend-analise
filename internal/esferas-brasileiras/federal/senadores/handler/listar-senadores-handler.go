package handler

import (
	"net/http"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/senadores/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/shared/logger"
	redis "github.com/danyele/laceu/internal/shared/redis"
)

type ListarSenadoresHandler struct {
	useCase *usecase.ListarSenadoresUseCase
	redis   *redis.RedisCache
}

func NovoListarSenadoresHandler(useCase *usecase.ListarSenadoresUseCase, redis *redis.RedisCache) *ListarSenadoresHandler {
	return &ListarSenadoresHandler{useCase: useCase, redis: redis}
}

func (h *ListarSenadoresHandler) Listar(c *gin.Context) {
	log := logger.New("Senadores: Handler: Listar")
	chave := redis.ChaveCache("senado-senadores", []byte("all"))

	var cached []interface{}
	ok, err := h.redis.Get(c.Request.Context(), chave, &cached)
	if err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}
	if ok {
		c.JSON(http.StatusOK, gin.H{"dados": cached})
		return
	}

	senadores, err := h.useCase.Listar(c.Request.Context())
	if err != nil {
		log.Error("erro ao listar senadores", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar senadores: " + err.Error()})
		return
	}

	if err := h.redis.Set(c.Request.Context(), chave, senadores); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}

	c.JSON(http.StatusOK, gin.H{"dados": senadores})
}
