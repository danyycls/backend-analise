package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/shared/clients/senado"
)

type ListarVotacoesComissaoUseCase struct {
	client *senado.SenadoClient
}

func NovoListarVotacoesComissaoUseCase(c *senado.SenadoClient) *ListarVotacoesComissaoUseCase {
	return &ListarVotacoesComissaoUseCase{client: c}
}

func (u *ListarVotacoesComissaoUseCase) Listar(ctx context.Context, sigla string, params map[string]string) ([]senado.VotacaoComissao, error) {
	return u.client.ListarVotacoesComissao(ctx, sigla, params)
}
