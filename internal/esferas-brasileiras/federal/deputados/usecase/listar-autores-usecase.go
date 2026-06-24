package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type ListarAutoresUseCase struct {
	client *deputados.DeputadosClient
}

func NovoListarAutoresUseCase(c *deputados.DeputadosClient) *ListarAutoresUseCase {
	return &ListarAutoresUseCase{client: c}
}

func (u *ListarAutoresUseCase) Executar(ctx context.Context, idProposicao int) ([]deputados.Author, error) {
	return u.client.ListarAutores(ctx, idProposicao)
}
