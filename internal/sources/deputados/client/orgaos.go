package deputados

import (
	"context"
	"fmt"
	"net/url"
)

func (c *DeputadosClient) ListarOrgaosCamara(ctx context.Context, params map[string]string) ([]Orgao, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []Orgao
	if err := c.doGet(ctx, "/orgaos", query, &out); err != nil {
		return nil, fmt.Errorf("listar orgaos camara: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarOrgaoCamara(ctx context.Context, id int) (*OrgaoDetalhe, error) {
	var out OrgaoDetalhe
	if err := c.doGet(ctx, fmt.Sprintf("/orgaos/%d", id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar orgao camara: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarMembrosOrgaoCamara(ctx context.Context, id int, params map[string]string) ([]MembroOrgao, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []MembroOrgao
	if err := c.doGet(ctx, fmt.Sprintf("/orgaos/%d/membros", id), query, &out); err != nil {
		return nil, fmt.Errorf("listar membros orgao: %w", err)
	}
	return out, nil
}
