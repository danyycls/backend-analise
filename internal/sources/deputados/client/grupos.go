package deputados

import (
	"context"
	"fmt"
	"strconv"
)

func (c *DeputadosClient) ListarGrupos(ctx context.Context) ([]Grupo, error) {
	var out []Grupo
	if err := c.doGet(ctx, "/grupos", nil, &out); err != nil {
		return nil, fmt.Errorf("listar grupos: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarGrupo(ctx context.Context, id int) (*GrupoDetalhe, error) {
	var out GrupoDetalhe
	if err := c.doGet(ctx, "/grupos/"+strconv.Itoa(id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar grupo: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarHistoricoGrupo(ctx context.Context, id int) ([]HistoricoGrupo, error) {
	var out []HistoricoGrupo
	path := fmt.Sprintf("/grupos/%d/historico", id)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar historico grupo: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarMembrosGrupo(ctx context.Context, id int) ([]MembroGrupo, error) {
	var out []MembroGrupo
	path := fmt.Sprintf("/grupos/%d/membros", id)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar membros grupo: %w", err)
	}
	return out, nil
}
