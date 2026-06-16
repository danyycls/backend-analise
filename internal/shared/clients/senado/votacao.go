package senado

import (
	"context"
	"fmt"
)

func (c *SenadoClient) ListarVotacoes(ctx context.Context, params map[string]string) ([]VotacaoItem, error) {
	var resultado []VotacaoItem
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, "/votacao", query, &resultado); err != nil {
		return nil, fmt.Errorf("listar votacoes: %w", err)
	}
	return resultado, nil
}

func (c *SenadoClient) ListarVotacoesComissao(ctx context.Context, siglaComissao string, params map[string]string) ([]VotacaoComissao, error) {
	var resultado VotacaoComissaoWrapper
	path := fmt.Sprintf("/votacaoComissao/comissao/%s", siglaComissao)
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, path, query, &resultado); err != nil {
		return nil, fmt.Errorf("listar votacoes comissao: %w", err)
	}
	return resultado.VotacoesComissao.Votacoes.Votacao, nil
}

func (c *SenadoClient) ListarVotacoesComissaoParlamentar(ctx context.Context, codigo string, params map[string]string) ([]VotacaoComissao, error) {
	var resultado VotacaoComissaoParlamentar
	path := fmt.Sprintf("/votacaoComissao/parlamentar/%s", codigo)
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, path, query, &resultado); err != nil {
		return nil, fmt.Errorf("listar votacoes comissao parlamentar: %w", err)
	}
	return resultado.VotacoesParlamentar.Votacoes.Votacao, nil
}

func (c *SenadoClient) ListarMateriaTramitacao(ctx context.Context, params map[string]string) ([]MateriaItem, error) {
	var resultado MateriaTramitacao
	path := "/materia/lista/tramitacao"
	query := makeQueryParams(params)
	if err := c.doGetJSON(ctx, path, query, &resultado); err != nil {
		return nil, fmt.Errorf("listar materia tramitacao: %w", err)
	}
	return resultado.Materials.Materia, nil
}
