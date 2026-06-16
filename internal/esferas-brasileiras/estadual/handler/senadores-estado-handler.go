package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/laceu/internal/esferas-brasileiras/estadual/usecase"
	"github.com/danyele/laceu/internal/shared/logger"
)

type EsferaEstadualBuscarSenadoresHandler struct {
	useCase *usecase.EsferaEstadualBuscarSenadoresUseCase
}

func NovoEsferaEstadualBuscarSenadoresHandler(useCase *usecase.EsferaEstadualBuscarSenadoresUseCase) *EsferaEstadualBuscarSenadoresHandler {
	return &EsferaEstadualBuscarSenadoresHandler{useCase: useCase}
}

func (h *EsferaEstadualBuscarSenadoresHandler) BuscarSenadoresEstado(c *gin.Context) {
	log := logger.New("Estadual: Handler: BuscarSenadoresEstado")
	uf := c.Param("uf")
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	resp, err := h.useCase.Executar(c.Request.Context(), &usecase.EsferaEstadualBuscarSenadoresRequest{UF: uf})
	if err != nil {
		log.Error("erro ao buscar senadores UF", "uf", uf, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar senadores"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp.Senadores})
}
