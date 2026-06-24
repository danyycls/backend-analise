package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type ListarProposicoesUseCase struct {
	client *deputados.DeputadosClient
}

func NovoListarProposicoesUseCase(c *deputados.DeputadosClient) *ListarProposicoesUseCase {
	return &ListarProposicoesUseCase{client: c}
}

func (u *ListarProposicoesUseCase) Executar(ctx context.Context, params map[string]string) ([]deputados.Proposicao, error) {
	return u.client.ListarProposicoes(ctx, params)
}
