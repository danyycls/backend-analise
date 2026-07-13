package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalBuscarFrenteUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarFrenteUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarFrenteUseCase {
	return &EsferaFederalBuscarFrenteUseCase{client: c}
}

func (u *EsferaFederalBuscarFrenteUseCase) Executar(ctx context.Context, id int) (*deputados.FrenteDetalhe, error) {
	return u.client.BuscarFrente(ctx, id)
}
