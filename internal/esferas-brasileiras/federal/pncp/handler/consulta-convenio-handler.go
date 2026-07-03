package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/esferas-brasileiras/federal/pncp/usecase"
)

type ConsultaConvenioHandler struct {
	useCase *usecase.ConsultaConvenioUseCase
}

func NovoConsultaConvenioHandler(useCase *usecase.ConsultaConvenioUseCase) *ConsultaConvenioHandler {
	return &ConsultaConvenioHandler{useCase: useCase}
}

func (h *ConsultaConvenioHandler) Listar(c *gin.Context) {
	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.DefaultQuery("por_pagina", "10"))
	uf := c.Query("uf")
	municipio := c.Query("municipio")
	tipo := c.Query("tipo")

	result, err := h.useCase.Listar(c.Request.Context(), pagina, porPagina, uf, municipio, tipo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
