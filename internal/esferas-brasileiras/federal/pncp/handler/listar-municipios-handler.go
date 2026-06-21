package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/clients/ibge"
)

type ListarMunicipiosHandler struct {
	ibgeClient *ibge.IBGEClient
}

func NovoListarMunicipiosHandler(ibgeClient *ibge.IBGEClient) *ListarMunicipiosHandler {
	return &ListarMunicipiosHandler{
		ibgeClient: ibgeClient,
	}
}

func (h *ListarMunicipiosHandler) ListarMunicipios(c *gin.Context) {
	uf := c.Param("uf")
	if uf == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF é obrigatória"})
		return
	}

	municipios, err := h.ibgeClient.ListarMunicipios(c.Request.Context(), uf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, municipios)
}
