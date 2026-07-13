package portaltransparencia

import (
	"context"
	"strconv"
)

func (c *PortalTransparenciaClient) ListarEmendas(ctx context.Context, filtro EmendaQueryParams) ([]ConsultaEmendas, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.CodigoEmenda != "" {
		params["codigoEmenda"] = filtro.CodigoEmenda
	}
	if filtro.NumeroEmenda != "" {
		params["numeroEmenda"] = filtro.NumeroEmenda
	}
	if filtro.NomeAutor != "" {
		params["nomeAutor"] = filtro.NomeAutor
	}
	if filtro.TipoEmenda != "" {
		params["tipoEmenda"] = filtro.TipoEmenda
	}
	if filtro.Ano != "" {
		params["ano"] = filtro.Ano
	}
	if filtro.CodigoFuncao != "" {
		params["codigoFuncao"] = filtro.CodigoFuncao
	}
	if filtro.CodigoSubfuncao != "" {
		params["codigoSubfuncao"] = filtro.CodigoSubfuncao
	}
	var result []ConsultaEmendas
	err := c.doGet(ctx, "/api-de-dados/emendas", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarDocumentosEmenda(ctx context.Context, codigo string, pagina int) ([]DocumentoRelacionadoEmenda, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(pagina),
	}
	var result []DocumentoRelacionadoEmenda
	err := c.doGet(ctx, "/api-de-dados/emendas/documentos/"+codigo, params, &result)
	return result, err
}
