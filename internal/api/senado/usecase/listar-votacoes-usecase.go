package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type ListarVotacoesUseCase struct {
	client *senado.SenadoClient
}

func NovoListarVotacoesUseCase(c *senado.SenadoClient) *ListarVotacoesUseCase {
	return &ListarVotacoesUseCase{client: c}
}

func (u *ListarVotacoesUseCase) Listar(ctx context.Context, params map[string]string) ([]senado.VotacaoItem, error) {
	return u.client.ListarVotacoes(ctx, params)
}
