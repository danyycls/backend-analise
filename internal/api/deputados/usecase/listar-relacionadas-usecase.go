package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type ListarRelacionadasUseCase struct {
	client *deputados.DeputadosClient
}

func NovoListarRelacionadasUseCase(c *deputados.DeputadosClient) *ListarRelacionadasUseCase {
	return &ListarRelacionadasUseCase{client: c}
}

func (u *ListarRelacionadasUseCase) Executar(ctx context.Context, idProposicao int) ([]deputados.Proposicao, error) {
	return u.client.ListarRelacionadas(ctx, idProposicao)
}
