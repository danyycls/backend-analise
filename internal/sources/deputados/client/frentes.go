package deputados

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func (c *DeputadosClient) ListarFrentes(ctx context.Context, idLegislatura int) ([]Frente, error) {
	query := url.Values{}
	query.Set("idLegislatura", strconv.Itoa(idLegislatura))
	var out []Frente
	if err := c.doGet(ctx, "/frentes", query, &out); err != nil {
		return nil, fmt.Errorf("listar frentes: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarFrente(ctx context.Context, id int) (*FrenteDetalhe, error) {
	var out FrenteDetalhe
	if err := c.doGet(ctx, "/frentes/"+strconv.Itoa(id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar frente: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarMembrosFrente(ctx context.Context, idFrente int) ([]MembroFrente, error) {
	var out []MembroFrente
	path := fmt.Sprintf("/frentes/%d/membros", idFrente)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar membros frente: %w", err)
	}
	return out, nil
}
