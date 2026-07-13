package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type BuscarProposicaoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoBuscarProposicaoUseCase(c *deputados.DeputadosClient) *BuscarProposicaoUseCase {
	return &BuscarProposicaoUseCase{client: c}
}

func (u *BuscarProposicaoUseCase) Executar(ctx context.Context, id int) (*deputados.ProposicaoDetalhe, error) {
	return u.client.BuscarProposicao(ctx, id)
}
