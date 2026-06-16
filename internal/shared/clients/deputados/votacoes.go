package deputados

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func (c *DeputadosClient) ListarVotacoes(ctx context.Context, params map[string]string) ([]Votacao, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []Votacao
	if err := c.doGet(ctx, "/votacoes", query, &out); err != nil {
		return nil, fmt.Errorf("listar votacoes: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarVotacao(ctx context.Context, id int) (*VotacaoDetalhe, error) {
	var out VotacaoDetalhe
	if err := c.doGet(ctx, "/votacoes/"+strconv.Itoa(id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar votacao: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarOrientacoes(ctx context.Context, idVotacao int) ([]Orientacao, error) {
	var out []Orientacao
	path := fmt.Sprintf("/votacoes/%d/orientacoes", idVotacao)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar orientacoes: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarVotos(ctx context.Context, idVotacao int) ([]Voto, error) {
	var out []Voto
	path := fmt.Sprintf("/votacoes/%d/votos", idVotacao)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar votos: %w", err)
	}
	return out, nil
}
