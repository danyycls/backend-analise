package portaltransparencia

import (
	"context"
	"strconv"
)

func (c *PortalTransparenciaClient) ListarOrgaosSIAPE(ctx context.Context, filtro OrgaoQueryParams) ([]Orgao, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.Codigo != "" {
		params["codigo"] = filtro.Codigo
	}
	if filtro.Descricao != "" {
		params["descricao"] = filtro.Descricao
	}
	var result []Orgao
	err := c.doGet(ctx, "/api-de-dados/orgaos-siape", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarOrgaosSIAFI(ctx context.Context, filtro OrgaoQueryParams) ([]Orgao, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.Codigo != "" {
		params["codigo"] = filtro.Codigo
	}
	if filtro.Descricao != "" {
		params["descricao"] = filtro.Descricao
	}
	var result []Orgao
	err := c.doGet(ctx, "/api-de-dados/orgaos-siafi", params, &result)
	return result, err
}
