package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type ListarTemasUseCase struct {
	client *deputados.DeputadosClient
}

func NovoListarTemasUseCase(c *deputados.DeputadosClient) *ListarTemasUseCase {
	return &ListarTemasUseCase{client: c}
}

func (u *ListarTemasUseCase) Executar(ctx context.Context, idProposicao int) ([]deputados.Tema, error) {
	return u.client.ListarTemas(ctx, idProposicao)
}
