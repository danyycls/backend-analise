package senado

import (
	"context"
	"fmt"
)

func (c *SenadoClient) ListarProcessos(ctx context.Context, params map[string]string) ([]ProcessoItem, error) {
	var resultado []ProcessoItem
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, "/processo", query, &resultado); err != nil {
		return nil, fmt.Errorf("listar processors: %w", err)
	}
	return resultado, nil
}

func (c *SenadoClient) ListarProcessoAssuntos(ctx context.Context) ([]ProcessoAssunto, error) {
	var resultado []ProcessoAssunto
	if err := c.doGetJSON(ctx, "/processo/assuntos", nil, &resultado); err != nil {
		return nil, fmt.Errorf("listar assuntos processo: %w", err)
	}
	return resultado, nil
}

func (c *SenadoClient) ListarProcessoEmendas(ctx context.Context, params map[string]string) ([]ProcessoEmenda, error) {
	var resultado []ProcessoEmenda
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, "/processo/emenda", query, &resultado); err != nil {
		return nil, fmt.Errorf("listar emendas processo: %w", err)
	}
	return resultado, nil
}

func (c *SenadoClient) BuscarProcesso(ctx context.Context, id string) (*ProcessoItem, error) {
	var resultado ProcessoItem
	path := fmt.Sprintf("/processo/%s", id)
	if err := c.doGetJSON(ctx, path, nil, &resultado); err != nil {
		return nil, fmt.Errorf("buscar processo: %w", err)
	}
	return &resultado, nil
}
