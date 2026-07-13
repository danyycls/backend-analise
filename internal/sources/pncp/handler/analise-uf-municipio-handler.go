package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	pncp "github.com/danyele/podp/internal/sources/pncp/client"
	"github.com/danyele/podp/internal/sources/pncp/usecase"
)

type AnaliseUFMunicipioHandler struct {
	*JobManager
	useCase *usecase.ConsultaContratoUFMunicipioPNCPUseCase
}

func NovoAnaliseUFMunicipioHandler(useCase *usecase.ConsultaContratoUFMunicipioPNCPUseCase) *AnaliseUFMunicipioHandler {
	return &AnaliseUFMunicipioHandler{
		JobManager: NovoJobManager(),
		useCase:    useCase,
	}
}

func (h *AnaliseUFMunicipioHandler) AnaliseUFMunicipio(c *gin.Context) {
	var req pncp.AnaliseContratoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if req.Tipo != "uf" && req.Tipo != "municipio" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "tipo deve ser 'uf' ou 'municipio'"})
		return
	}

	jobID, eventChan, _ := h.CriarJob()

	go func() {
		defer close(eventChan)

		eventChan <- pncp.EventoAnalise{Type: "started", CNPJ: req.Tipo, Total: 1}

		var results []*pncp.AnaliseResultado
		var err error
		var paginasErro []int

		if req.Tipo == "municipio" {
			results, err = h.useCase.BuscarPorMunicipio(context.Background(), req.CodigoMunicipioIbge, req.DataInicial, req.DataFinal, &paginasErro)
		} else {
			results, err = h.useCase.BuscarPorUF(context.Background(), req.UF, req.DataInicial, req.DataFinal, &paginasErro)
		}

		if err != nil {
			eventChan <- pncp.EventoAnalise{Type: "error", Message: err.Error()}
			results = []*pncp.AnaliseResultado{}
		} else {
			totalContratos := 0
			totalEmpresas := 0
			var valorTotal float64
			for _, r := range results {
				if r.Resumo == nil {
					continue
				}
				if r.Resumo.TotalContratos != nil {
					totalContratos += *r.Resumo.TotalContratos
				}
				if r.Resumo.TotalEmpresas != nil {
					totalEmpresas += *r.Resumo.TotalEmpresas
				}
				if r.Resumo.ValorTotalContratos != nil {
					valorTotal += *r.Resumo.ValorTotalContratos
				}
			}

			eventChan <- pncp.EventoAnalise{
				Type:                "success",
				TotalContratos:      totalContratos,
				ValorTotalContratos: valorTotal,
				CNPJ:                req.UF,
				Orgao:               fmt.Sprintf("Busca por %s", req.Tipo),
			}
		}

		eventChan <- pncp.EventoAnalise{
			Type: "progress", Processed: 1, Total: 1,
			Success: 1, Errors: 0,
		}

		eventChan <- pncp.EventoAnalise{Type: "completed", Total: 1}

		eventChan <- pncp.EventoAnalise{
			Type:    "results",
			Results: results,
		}

		h.FinalizarJob(jobID, results, paginasErro)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"jobId":  jobID,
		"status": "processing",
		"total":  1,
	})
}
