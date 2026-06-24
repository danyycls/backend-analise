package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	usecase "github.com/danyele/podp/internal/esferas-brasileiras/federal/pncp/usecase"
	pncp "github.com/danyele/podp/internal/shared/clients/pncp"
	"github.com/danyele/podp/internal/shared/logger"
)

type BuscarLicitacoesUFHandler struct {
	useCase *usecase.ConsultaPublicacaoPNCPUseCase
}

func NovoBuscarLicitacoesUFHandler(useCase *usecase.ConsultaPublicacaoPNCPUseCase) *BuscarLicitacoesUFHandler {
	return &BuscarLicitacoesUFHandler{useCase: useCase}
}

func (h *BuscarLicitacoesUFHandler) Buscar(c *gin.Context) {
	log := logger.New("PNCP: Handler: BuscarLicitacoesUF")
	uf := c.Param("uf")
	if uf == "" || len(uf) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "UF invalida"})
		return
	}

	codigo := c.Param("codigo")

	anoStr := c.Query("ano")
	if anoStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "parametro 'ano' obrigatorio"})
		return
	}

	ano, err := strconv.ParseInt(anoStr, 10, 64)
	if err != nil || ano < 2000 || ano > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ano invalido"})
		return
	}

	trimestreStr := c.Query("trimestre")

	var dataInicial, dataFinal string

	if trimestreStr != "" {
		trimestre, err := strconv.ParseInt(trimestreStr, 10, 64)
		if err != nil || trimestre < 1 || trimestre > 4 {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "trimestre invalido (use 1, 2, 3 ou 4)"})
			return
		}
		switch trimestre {
		case 1:
			dataInicial = anoStr + "0101"
			dataFinal = anoStr + "0331"
		case 2:
			dataInicial = anoStr + "0401"
			dataFinal = anoStr + "0630"
		case 3:
			dataInicial = anoStr + "0701"
			dataFinal = anoStr + "0930"
		case 4:
			dataInicial = anoStr + "1001"
			dataFinal = anoStr + "1231"
		}
	} else {
		dataInicial = anoStr + "0101"
		dataFinal = anoStr + "1231"
	}

	var resultados []*pncp.AnaliseResultado

	if codigo != "" {
		resultados, err = h.useCase.BuscarPorMunicipio(c.Request.Context(), codigo, dataInicial, dataFinal, "", nil)
		if err != nil {
			log.Error("erro ao buscar licitacoes por municipio", "uf", uf, "codigo", codigo, "ano", ano, "trimestre", trimestreStr, "erro", err)
			c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar licitacoes"})
			return
		}
	} else {
		resultados, err = h.useCase.BuscarPorUF(c.Request.Context(), uf, dataInicial, dataFinal, "", nil)
		if err != nil {
			log.Error("erro ao buscar licitacoes por UF", "uf", uf, "ano", ano, "trimestre", trimestreStr, "erro", err)
			c.JSON(http.StatusBadGateway, gin.H{"erro": "falha ao buscar licitacoes"})
			return
		}
	}

	contratos := make([]pncp.Contrato, 0)
	for _, r := range resultados {
		if r.Contratos != nil {
			contratos = append(contratos, r.Contratos...)
		}
	}

	c.JSON(http.StatusOK, contratos)
}
