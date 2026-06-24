package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	"github.com/danyele/podp/internal/shared/logger"
)

type EsferaFederalListarMembrosOrgaoHandler struct {
	useCase *usecase.EsferaFederalListarMembrosOrgaoUseCase
}

func NovoEsferaFederalListarMembrosOrgaoHandler(useCase *usecase.EsferaFederalListarMembrosOrgaoUseCase) *EsferaFederalListarMembrosOrgaoHandler {
	return &EsferaFederalListarMembrosOrgaoHandler{useCase: useCase}
}

func (h *EsferaFederalListarMembrosOrgaoHandler) Listar(c *gin.Context) {
	log := logger.New("Deputados: Handler: ListarMembrosOrgao")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id invalido"})
		return
	}

	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resp, err := h.useCase.Executar(c.Request.Context(), id, params)
	if err != nil {
		log.Error("erro ao listar membros do orgao", "id", id, "erro", err)
		c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao listar membros do orgao: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dados": resp})
}
