package portaltransparencia

import (
	"context"
	"strconv"
	"strings"
	"time"
)

func converterDataISOparaBR(dataISO string) string {
	if dataISO == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", dataISO)
	if err != nil {
		return dataISO
	}
	return t.Format("02/01/2006")
}

func converterMesAnoISOparaBR(mesAnoISO string) string {
	if mesAnoISO == "" {
		return ""
	}
	t, err := time.Parse("2006-01", mesAnoISO)
	if err != nil {
		return strings.ReplaceAll(mesAnoISO, "-", "/")
	}
	return t.Format("01/2006")
}

func mapearTipoCartao(tipo string) string {
	switch strings.ToLower(tipo) {
	case "1", "cpgf", "corporate":
		return "1"
	case "2", "cpcc", "compras":
		return "2"
	case "3", "cpdc":
		return "3"
	}
	return tipo
}

func (c *PortalTransparenciaClient) ListarCartoes(ctx context.Context, filtro CartaoQueryParams) ([]Cartao, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.MesExtratoInicio != "" {
		params["mesExtratoInicio"] = converterMesAnoISOparaBR(filtro.MesExtratoInicio)
	}
	if filtro.MesExtratoFim != "" {
		params["mesExtratoFim"] = converterMesAnoISOparaBR(filtro.MesExtratoFim)
	}
	if filtro.DataTransacaoInicio != "" {
		params["dataTransacaoInicio"] = converterDataISOparaBR(filtro.DataTransacaoInicio)
	}
	if filtro.DataTransacaoFim != "" {
		params["dataTransacaoFim"] = converterDataISOparaBR(filtro.DataTransacaoFim)
	}
	if filtro.TipoCartao != "" {
		params["tipoCartao"] = mapearTipoCartao(filtro.TipoCartao)
	}
	if filtro.CodigoOrgao != "" {
		params["codigoOrgao"] = filtro.CodigoOrgao
	}
	if filtro.CPFPortador != "" {
		params["cpfPortador"] = filtro.CPFPortador
	}
	if filtro.CPFCNPJFavorecido != "" {
		params["cpfCnpjFavorecido"] = filtro.CPFCNPJFavorecido
	}
	if filtro.ValorDe != "" {
		params["valorDe"] = filtro.ValorDe
	}
	if filtro.ValorAte != "" {
		params["valorAte"] = filtro.ValorAte
	}
	var result []Cartao
	err := c.doGet(ctx, "/api-de-dados/cartoes", params, &result)
	return result, err
}
