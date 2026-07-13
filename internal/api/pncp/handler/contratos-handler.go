package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pncpClient "github.com/danyele/podp/internal/sources/pncp/client"
)

type PNCPContratosHandler struct {
	client *pncpClient.PNCPClient
}

func NovoPNCPContratosHandler(client *pncpClient.PNCPClient) *PNCPContratosHandler {
	return &PNCPContratosHandler{client: client}
}

type buscarContratosRequest struct {
	DataInicial      string `json:"data_inicial"`
	DataFinal        string `json:"data_final"`
	CodigoModalidade string `json:"codigo_modalidade"`
	Pagina           int    `json:"pagina"`
	Tamanho          int    `json:"tamanho"`
}

func (h *PNCPContratosHandler) BuscarPorMunicipio(c *gin.Context) {
	codigoIBGE := c.Param("codigoIBGE")
	if codigoIBGE == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "codigo IBGE e obrigatorio"})
		return
	}
	var req buscarContratosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	if req.Pagina == 0 {
		req.Pagina = 1
	}
	if req.Tamanho == 0 {
		req.Tamanho = 500
	}
	result, err := h.client.BuscarContratosPorMunicipio(c.Request.Context(), codigoIBGE, req.DataInicial, req.DataFinal, req.CodigoModalidade, req.Pagina, req.Tamanho)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *PNCPContratosHandler) BuscarPorUF(c *gin.Context) {
	uf := c.Param("uf")
	if uf == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF e obrigatoria"})
		return
	}
	var req buscarContratosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	if req.Pagina == 0 {
		req.Pagina = 1
	}
	if req.Tamanho == 0 {
		req.Tamanho = 500
	}
	result, err := h.client.BuscarContratosPorUF(c.Request.Context(), uf, req.DataInicial, req.DataFinal, req.CodigoModalidade, req.Pagina, req.Tamanho)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *PNCPContratosHandler) BuscarPorOrgao(c *gin.Context) {
	cnpj := c.Param("cnpj")
	if cnpj == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "CNPJ e obrigatorio"})
		return
	}
	var req buscarContratosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	if req.Pagina == 0 {
		req.Pagina = 1
	}
	if req.Tamanho == 0 {
		req.Tamanho = 500
	}
	result, err := h.client.BuscarContratos(c.Request.Context(), cnpj, req.DataInicial, req.DataFinal, req.Pagina, req.Tamanho)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
