package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/shared/logger"
	anomalia "github.com/danyele/podp/internal/worker/anomalia"
	"github.com/danyele/podp/internal/worker/anomalia/usecase"
)

type AnomaliaConsultaHandler struct {
	useCase *usecase.AnomaliaConsultaUseCase
}

func NovoAnomaliaConsultaHandler(useCase *usecase.AnomaliaConsultaUseCase) *AnomaliaConsultaHandler {
	return &AnomaliaConsultaHandler{useCase: useCase}
}

func (h *AnomaliaConsultaHandler) Listar(c *gin.Context) {
	log := logger.New("Worker: Anomalia: ConsultaHandler: Listar")

	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	porPagina, _ := strconv.Atoi(c.DefaultQuery("por_pagina", "50"))

	filtro := usecase.AnomaliaFiltro{
		Documento: c.Query("documento"),
		Uf:        c.Query("uf"),
		Municipio: c.Query("municipio"),
		Tag:       c.Query("tag"),
		Categoria: c.Query("categoria"),
		Pagina:    pagina,
		PorPagina: porPagina,
	}

	resultado, err := h.useCase.Listar(c.Request.Context(), filtro)
	if err != nil {
		log.Error("erro ao listar anomalias", "erro", err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "erro ao listar anomalias"})
		return
	}

	c.JSON(http.StatusOK, anomalia.ListarAnomaliasResponse{
		Total:     resultado.Total,
		Pagina:    pagina,
		PorPagina: porPagina,
		Anomalias: resultado.Anomalias,
	})
}
