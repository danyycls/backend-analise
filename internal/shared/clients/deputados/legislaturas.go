package deputados

import (
	"context"
	"fmt"
	"strconv"
)

func (c *DeputadosClient) ListarLegislaturas(ctx context.Context) ([]Legislatura, error) {
	var out []Legislatura
	if err := c.doGet(ctx, "/legislaturas", nil, &out); err != nil {
		return nil, fmt.Errorf("listar legislaturas: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarLegislatura(ctx context.Context, id int) (*Legislatura, error) {
	var out Legislatura
	if err := c.doGet(ctx, "/legislaturas/"+strconv.Itoa(id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar legislatura: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarLideres(ctx context.Context, idLegislatura int) ([]Lider, error) {
	var out []Lider
	path := fmt.Sprintf("/legislaturas/%d/lideres", idLegislatura)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar lideres: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarMesa(ctx context.Context, idLegislatura int) ([]MembroMesa, error) {
	var out []MembroMesa
	path := fmt.Sprintf("/legislaturas/%d/mesa", idLegislatura)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar mesa: %w", err)
	}
	return out, nil
}
