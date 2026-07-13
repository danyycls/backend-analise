package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/api/ibge/usecase"
)

type ListarMunicipiosHandler struct {
	useCase *usecase.ListarMunicipiosUseCase
}

func NovoListarMunicipiosHandler(useCase *usecase.ListarMunicipiosUseCase) *ListarMunicipiosHandler {
	return &ListarMunicipiosHandler{useCase: useCase}
}

func (h *ListarMunicipiosHandler) ListarMunicipios(c *gin.Context) {
	uf := c.Param("uf")
	if uf == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF é obrigatória"})
		return
	}

	municipios, err := h.useCase.Executar(c.Request.Context(), uf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, municipios)
}
