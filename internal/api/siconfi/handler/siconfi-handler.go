package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	siconfiClient "github.com/danyele/podp/internal/sources/siconfi/client"
)

type SICONFIHandler struct {
	client *siconfiClient.SICONFIClient
}

func NovoSICONFIHandler(client *siconfiClient.SICONFIClient) *SICONFIHandler {
	return &SICONFIHandler{client: client}
}

func (h *SICONFIHandler) ListarEntes(c *gin.Context) {
	result, err := h.client.ListarEntes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarDCARequest struct {
	AnExercicio int64  `json:"an_exercicio"`
	IdEnte      int    `json:"id_ente"`
	NoAnexo     string `json:"no_anexo"`
}

func (h *SICONFIHandler) BuscarDCA(c *gin.Context) {
	var req buscarDCARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	var anexo *string
	if req.NoAnexo != "" {
		anexo = &req.NoAnexo
	}
	var result []siconfiClient.DCAItem
	var err error
	if anexo != nil {
		result, err = h.client.BuscarDCA(c.Request.Context(), req.AnExercicio, req.IdEnte, *anexo)
	} else {
		result, err = h.client.BuscarDCA(c.Request.Context(), req.AnExercicio, req.IdEnte)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarRGFRequest struct {
	AnExercicio         int64  `json:"an_exercicio"`
	InPeriodicidade     string `json:"in_periodicidade"`
	NrPeriodo           int    `json:"nr_periodo"`
	CoTipoDemonstrativo string `json:"co_tipo_demonstration"`
	CoPoder             string `json:"co_poder"`
	IdEnte              int    `json:"id_ente"`
	NoAnexo             string `json:"no_anexo"`
	CoEsfera            string `json:"co_esfera"`
}

func (h *SICONFIHandler) BuscarRGF(c *gin.Context) {
	var req buscarRGFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	params := siconfiClient.RGFParams{
		AnExercicio:         req.AnExercicio,
		InPeriodicidade:     req.InPeriodicidade,
		NrPeriodo:           req.NrPeriodo,
		CoTipoDemonstrativo: req.CoTipoDemonstrativo,
		CoPoder:             req.CoPoder,
		IdEnte:              req.IdEnte,
		NoAnexo:             req.NoAnexo,
		CoEsfera:            req.CoEsfera,
	}
	result, err := h.client.BuscarRGF(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarRREORequest struct {
	AnExercicio         int64  `json:"an_exercicio"`
	NrPeriodo           int    `json:"nr_periodo"`
	CoTipoDemonstrativo string `json:"co_tipo_demonstration"`
	IdEnte              int    `json:"id_ente"`
	NoAnexo             string `json:"no_anexo"`
	CoEsfera            string `json:"co_esfera"`
}

func (h *SICONFIHandler) BuscarRREO(c *gin.Context) {
	var req buscarRREORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	params := siconfiClient.RREOParams{
		AnExercicio:         req.AnExercicio,
		NrPeriodo:           req.NrPeriodo,
		CoTipoDemonstrativo: req.CoTipoDemonstrativo,
		IdEnte:              req.IdEnte,
		NoAnexo:             req.NoAnexo,
		CoEsfera:            req.CoEsfera,
	}
	result, err := h.client.BuscarRREO(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarMSCRequest struct {
	IdEnte       int    `json:"id_ente"`
	AnReferencia int64  `json:"an_referencia"`
	MeReferencia int64  `json:"me_referencia"`
	CoTipoMatriz string `json:"co_tipo_matriz"`
	ClasseConta  int    `json:"classe_conta"`
	IdTV         string `json:"id_tv"`
}

func (h *SICONFIHandler) BuscarMSCPatrimonial(c *gin.Context) {
	var req buscarMSCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	params := siconfiClient.MSCParams{
		IdEnte:       req.IdEnte,
		AnReferencia: req.AnReferencia,
		MeReferencia: req.MeReferencia,
		CoTipoMatriz: req.CoTipoMatriz,
		ClasseConta:  req.ClasseConta,
		IdTV:         req.IdTV,
	}
	result, err := h.client.BuscarMSCPatrimonial(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SICONFIHandler) BuscarMSCOrcamentaria(c *gin.Context) {
	var req buscarMSCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	params := siconfiClient.MSCParams{
		IdEnte:       req.IdEnte,
		AnReferencia: req.AnReferencia,
		MeReferencia: req.MeReferencia,
		CoTipoMatriz: req.CoTipoMatriz,
		ClasseConta:  req.ClasseConta,
		IdTV:         req.IdTV,
	}
	result, err := h.client.BuscarMSCOrcamentaria(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SICONFIHandler) BuscarMSCControle(c *gin.Context) {
	var req buscarMSCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	params := siconfiClient.MSCParams{
		IdEnte:       req.IdEnte,
		AnReferencia: req.AnReferencia,
		MeReferencia: req.MeReferencia,
		CoTipoMatriz: req.CoTipoMatriz,
		ClasseConta:  req.ClasseConta,
		IdTV:         req.IdTV,
	}
	result, err := h.client.BuscarMSCControle(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarExtratoRequest struct {
	IdEnte       int   `json:"id_ente"`
	AnReferencia int64 `json:"an_referencia"`
}

func (h *SICONFIHandler) BuscarExtratoEntregas(c *gin.Context) {
	var req buscarExtratoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	result, err := h.client.BuscarExtratoEntregas(c.Request.Context(), req.IdEnte, req.AnReferencia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *SICONFIHandler) ListarAnexosRelatorios(c *gin.Context) {
	result, err := h.client.ListarAnexosRelatorios(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Unused import guard
var _ = strconv.Itoa
