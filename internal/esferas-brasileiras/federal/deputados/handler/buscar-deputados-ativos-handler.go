package handler

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/laceu/internal/shared/logger"
	redis "github.com/danyele/laceu/internal/shared/redis"
)

type EsferaFederalBuscarDeputadosAtivosHandler struct {
	useCase *usecase.EsferaFederalBuscarDeputadosAtivosUseCase
	redis   *redis.RedisCache
}

func NovoEsferaFederalBuscarDeputadosAtivosHandler(useCase *usecase.EsferaFederalBuscarDeputadosAtivosUseCase, redis *redis.RedisCache) *EsferaFederalBuscarDeputadosAtivosHandler {
	return &EsferaFederalBuscarDeputadosAtivosHandler{useCase: useCase, redis: redis}
}

func (h *EsferaFederalBuscarDeputadosAtivosHandler) BuscarDeputadosAtivos(c *gin.Context) {
	log := logger.New("Deputados: Handler: BuscarDeputadosAtivos")
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	var chave string
	if len(params) > 0 {
		keys := make([]string, 0, len(params))
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		normalized := make(map[string]string)
		for _, k := range keys {
			normalized[k] = params[k]
		}
		raw, _ := json.Marshal(normalized)
		chave = redis.ChaveCache("deputados", raw)
	} else {
		chave = redis.ChaveCache("deputados", []byte("all"))
	}

	var cached []interface{}
	ok, err := h.redis.Get(c.Request.Context(), chave, &cached)
	if err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}
	if ok {
		c.JSON(http.StatusOK, gin.H{"dados": cached})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaFederalBuscarDeputadosAtivosRequest{Params: params})
	if err != nil {
		log.Error("erro ao listar deputados", "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar deputados: " + err.Error()})
		return
	}

	if err := h.redis.Set(c.Request.Context(), chave, resp.Deputados); err != nil {
		log.Warn("cache indisponivel", "erro", err)
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp.Deputados})
}
