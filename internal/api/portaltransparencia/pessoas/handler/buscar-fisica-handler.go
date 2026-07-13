package handler

import (
	"net/http"

	"github.com/danyele/podp/internal/api/portaltransparencia/pessoas/usecase"

	"github.com/gin-gonic/gin"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

func boolQuery(c *gin.Context, key string) string {
	if c.Query(key) != "" {
		return c.Query(key)
	}
	return ""
}

type BuscarFisicaHandler struct {
	useCase *usecase.BuscarPessoasFisicasUseCase
}

func NovoBuscarFisicaHandler(useCase *usecase.BuscarPessoasFisicasUseCase) *BuscarFisicaHandler {
	return &BuscarFisicaHandler{useCase: useCase}
}

func (h *BuscarFisicaHandler) BuscarFisica(c *gin.Context) {
	filtro := portaltransparencia.PessoaFisicaQueryParams{
		CPF:                       c.Query("cpf"),
		Nome:                      c.Query("nome"),
		NIS:                       c.Query("nis"),
		FavorecidoDespesas:        boolQuery(c, "favorecidoDespesas"),
		Servidor:                  boolQuery(c, "servidor"),
		BeneficiarioDiarias:       boolQuery(c, "beneficiarioDiarias"),
		Permissionario:            boolQuery(c, "permissionario"),
		Contratado:                boolQuery(c, "contratado"),
		SancionadoCEIS:            boolQuery(c, "sancionadoCEIS"),
		SancionadoCNEP:            boolQuery(c, "sancionadoCNEP"),
		SancionadoCEPIM:           boolQuery(c, "sancionadoCEPIM"),
		SancionadoCEAF:            boolQuery(c, "sancionadoCEAF"),
		SancionadoAcordoLeniencia: boolQuery(c, "sancionadoAcordoLeniencia"),
		Ordenacao:                 c.Query("ordenacao"),
		OrdenacaoDirecao:          c.Query("ordenacaoDirecao"),
	}
	resultado, err := h.useCase.Buscar(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"erro": "portaltransparencia: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}
