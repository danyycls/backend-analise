package deputados

import (
	"context"
	"fmt"
	"net/url"
)

func (c *DeputadosClient) ListarPartidos(ctx context.Context, params map[string]string) ([]Partido, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []Partido
	if err := c.doGet(ctx, "/partidos", query, &out); err != nil {
		return nil, fmt.Errorf("listar partidos: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarPartido(ctx context.Context, id int) (*PartidoDetalhe, error) {
	var out PartidoDetalhe
	if err := c.doGet(ctx, fmt.Sprintf("/partidos/%d", id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar partido: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarMembrosPartido(ctx context.Context, id int, params map[string]string) ([]Deputado, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []Deputado
	if err := c.doGet(ctx, fmt.Sprintf("/partidos/%d/membros", id), query, &out); err != nil {
		return nil, fmt.Errorf("listar membros partido: %w", err)
	}
	return out, nil
}
