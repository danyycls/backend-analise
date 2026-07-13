package deputados

import (
	"context"
	"fmt"
	"net/url"
)

func (c *DeputadosClient) ListarEventos(ctx context.Context, params map[string]string) ([]Evento, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []Evento
	if err := c.doGet(ctx, "/eventos", query, &out); err != nil {
		return nil, fmt.Errorf("listar eventos: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarEvento(ctx context.Context, id int) (*EventoDetalhe, error) {
	var out EventoDetalhe
	if err := c.doGet(ctx, fmt.Sprintf("/eventos/%d", id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar evento: %w", err)
	}
	return &out, nil
}
