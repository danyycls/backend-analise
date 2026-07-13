package deputados

import (
	"context"
	"fmt"
	"net/url"
)

func (c *DeputadosClient) ListarBlocos(ctx context.Context, params map[string]string) ([]Bloco, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []Bloco
	if err := c.doGet(ctx, "/blocos", query, &out); err != nil {
		return nil, fmt.Errorf("listar blocos: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarBloco(ctx context.Context, id string) (*Bloco, error) {
	var out Bloco
	if err := c.doGet(ctx, "/blocos/"+id, nil, &out); err != nil {
		return nil, fmt.Errorf("buscar bloco: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarPartidosDoBloco(ctx context.Context, idBloco string) ([]Partido, error) {
	var out []Partido
	path := fmt.Sprintf("/blocos/%s/partidos", idBloco)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar partidos do bloco: %w", err)
	}
	return out, nil
}
