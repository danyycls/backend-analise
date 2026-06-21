package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/shared/clients/senado"
)

type ListarVotacoesComissaoParlamentarUseCase struct {
	client *senado.SenadoClient
}

func NovoListarVotacoesComissaoParlamentarUseCase(c *senado.SenadoClient) *ListarVotacoesComissaoParlamentarUseCase {
	return &ListarVotacoesComissaoParlamentarUseCase{client: c}
}

func (u *ListarVotacoesComissaoParlamentarUseCase) Listar(ctx context.Context, codigo string, params map[string]string) ([]senado.VotacaoComissao, error) {
	return u.client.ListarVotacoesComissaoParlamentar(ctx, codigo, params)
}
