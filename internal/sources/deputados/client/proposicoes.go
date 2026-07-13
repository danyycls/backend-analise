package deputados

import (
	"context"
	"fmt"
	"net/url"
)

func (c *DeputadosClient) ListarProposicoes(ctx context.Context, params map[string]string) ([]Proposicao, error) {
	query := make(url.Values)
	for k, v := range params {
		query.Set(k, v)
	}
	var out []Proposicao
	if err := c.doGet(ctx, "/proposicoes", query, &out); err != nil {
		return nil, fmt.Errorf("listar proposicoes: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) BuscarProposicao(ctx context.Context, id int) (*ProposicaoDetalhe, error) {
	var out ProposicaoDetalhe
	if err := c.doGet(ctx, fmt.Sprintf("/proposicoes/%d", id), nil, &out); err != nil {
		return nil, fmt.Errorf("buscar proposicao: %w", err)
	}
	return &out, nil
}

func (c *DeputadosClient) ListarTramitacoes(ctx context.Context, idProposicao int) ([]Tramitacao, error) {
	var out []Tramitacao
	path := fmt.Sprintf("/proposicoes/%d/tramitacoes", idProposicao)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar tramitacoes: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarAutores(ctx context.Context, idProposicao int) ([]Author, error) {
	var out []Author
	path := fmt.Sprintf("/proposicoes/%d/autores", idProposicao)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar autores: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarTemas(ctx context.Context, idProposicao int) ([]Tema, error) {
	var out []Tema
	path := fmt.Sprintf("/proposicoes/%d/temas", idProposicao)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar temas: %w", err)
	}
	return out, nil
}

func (c *DeputadosClient) ListarRelacionadas(ctx context.Context, idProposicao int) ([]Proposicao, error) {
	var out []Proposicao
	path := fmt.Sprintf("/proposicoes/%d/relacionadas", idProposicao)
	if err := c.doGet(ctx, path, nil, &out); err != nil {
		return nil, fmt.Errorf("listar relacionadas: %w", err)
	}
	return out, nil
}
