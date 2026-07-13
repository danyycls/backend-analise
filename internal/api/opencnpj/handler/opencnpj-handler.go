package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	opencnpjClient "github.com/danyele/podp/internal/sources/opencnpj/client"
)

type OpenCNPJHandler struct {
	client *opencnpjClient.OpenCNPJClient
}

func NovoOpenCNPJHandler(client *opencnpjClient.OpenCNPJClient) *OpenCNPJHandler {
	return &OpenCNPJHandler{client: client}
}

func (h *OpenCNPJHandler) Buscar(c *gin.Context) {
	cnpj := c.Param("cnpj")
	if cnpj == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "CNPJ e obrigatorio"})
		return
	}
	result, err := h.client.Buscar(c.Request.Context(), cnpj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
