package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalListarVotosUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarVotosUseCase(c *deputados.DeputadosClient) *EsferaFederalListarVotosUseCase {
	return &EsferaFederalListarVotosUseCase{client: c}
}

func (u *EsferaFederalListarVotosUseCase) Executar(ctx context.Context, id int) ([]deputados.Voto, error) {
	return u.client.ListarVotos(ctx, id)
}
