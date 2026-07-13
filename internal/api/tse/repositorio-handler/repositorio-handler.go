package repositoriotse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/danyele/podp/internal/api/tse/repositorio"
)

type TSERepositorioHandler struct {
	repo *repositorio.Repositorio
}

func NovoTSERepositorioHandler(repo *repositorio.Repositorio) *TSERepositorioHandler {
	return &TSERepositorioHandler{repo: repo}
}

func (h *TSERepositorioHandler) CargosDistintos(c *gin.Context) {
	result, err := h.repo.CargosDistintos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarPorFiltrosRequest struct {
	CargoNome           *string `json:"cargo_nome"`
	PartidoID           *string `json:"partido_id"`
	UFSigla             *string `json:"uf_sigla"`
	SituacaoTotalizacao *string `json:"situacao_totalizacao"`
}

func (h *TSERepositorioHandler) CandidatoBuscarPorFiltros(c *gin.Context) {
	var req buscarPorFiltrosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	result, err := h.repo.CandidatoBuscarPorFiltros(c.Request.Context(), req.CargoNome, req.PartidoID, req.UFSigla, req.SituacaoTotalizacao)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarPorCPFRequest struct {
	CPF string `json:"cpf"`
}

func (h *TSERepositorioHandler) CandidatosBuscarPorCPF(c *gin.Context) {
	var req buscarPorCPFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	result, err := h.repo.CandidatosBuscarPorCPF(c.Request.Context(), req.CPF)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarPorIDRequest struct {
	ID string `json:"id"`
}

func (h *TSERepositorioHandler) CandidatoBuscarPorID(c *gin.Context) {
	var req buscarPorIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UUID invalido"})
		return
	}
	result, err := h.repo.CandidatoBuscarPorID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarPorDocumentoRequest struct {
	Documentos []string `json:"documentos"`
}

func (h *TSERepositorioHandler) FornecedoresBuscarPorDocumento(c *gin.Context) {
	var req buscarPorDocumentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	result, err := h.repo.FornecedoresBuscarPorDocumento(c.Request.Context(), req.Documentos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TSERepositorioHandler) DoadoresBuscarPorDocumento(c *gin.Context) {
	var req buscarPorDocumentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	result, err := h.repo.DoadoresBuscarPorDocumento(c.Request.Context(), req.Documentos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarPorDoadorIDRequest struct {
	DoadorID string `json:"doador_id"`
}

func (h *TSERepositorioHandler) ReceitasCandidatoBuscarPorDoadorID(c *gin.Context) {
	var req buscarPorDoadorIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	id, err := uuid.Parse(req.DoadorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UUID invalido"})
		return
	}
	result, err := h.repo.ReceitasCandidatoBuscarPorDoadorID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TSERepositorioHandler) ReceitasOrgaoBuscarPorDoadorID(c *gin.Context) {
	var req buscarPorDoadorIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	id, err := uuid.Parse(req.DoadorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UUID invalido"})
		return
	}
	result, err := h.repo.ReceitasOrgaoBuscarPorDoadorID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarPorFornecedorIDRequest struct {
	FornecedorID string `json:"fornecedor_id"`
}

func (h *TSERepositorioHandler) DespesasCandidatoBuscarPorFornecedorID(c *gin.Context) {
	var req buscarPorFornecedorIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	id, err := uuid.Parse(req.FornecedorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UUID invalido"})
		return
	}
	result, err := h.repo.DespesasCandidatoBuscarPorFornecedorID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TSERepositorioHandler) DespesasPartidoBuscarPorFornecedorID(c *gin.Context) {
	var req buscarPorFornecedorIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	id, err := uuid.Parse(req.FornecedorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UUID invalido"})
		return
	}
	result, err := h.repo.DespesasPartidoBuscarPorFornecedorID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type buscarPorIDsRequest struct {
	IDs []string `json:"ids"`
}

func (h *TSERepositorioHandler) PartidosBuscarPorIDs(c *gin.Context) {
	var req buscarPorIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, s := range req.IDs {
		id, err := uuid.Parse(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "UUID invalido: " + s})
			return
		}
		ids = append(ids, id)
	}
	result, err := h.repo.PartidosBuscarPorIDs(c.Request.Context(), ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TSERepositorioHandler) EleicoesBuscarPorIDs(c *gin.Context) {
	var req buscarPorIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, s := range req.IDs {
		id, err := uuid.Parse(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "UUID invalido: " + s})
			return
		}
		ids = append(ids, id)
	}
	result, err := h.repo.EleicoesBuscarPorIDs(c.Request.Context(), ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type candidatosEleitosRequest struct {
	UFSigla string   `json:"uf_sigla"`
	Cargos  []string `json:"cargos"`
}

func (h *TSERepositorioHandler) CandidatosEleitosPorUF(c *gin.Context) {
	var req candidatosEleitosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}
	result, err := h.repo.CandidatosEleitosPorUF(c.Request.Context(), req.UFSigla, req.Cargos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
