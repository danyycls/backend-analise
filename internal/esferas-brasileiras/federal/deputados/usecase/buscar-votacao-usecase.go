package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalBuscarVotacaoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarVotacaoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarVotacaoUseCase {
	return &EsferaFederalBuscarVotacaoUseCase{client: c}
}

func (u *EsferaFederalBuscarVotacaoUseCase) Executar(ctx context.Context, id int) (*deputados.VotacaoDetalhe, error) {
	return u.client.BuscarVotacao(ctx, id)
}
