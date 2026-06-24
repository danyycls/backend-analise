package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalBuscarPartidoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarPartidoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarPartidoUseCase {
	return &EsferaFederalBuscarPartidoUseCase{client: c}
}

func (u *EsferaFederalBuscarPartidoUseCase) Executar(ctx context.Context, id int) (*deputados.PartidoDetalhe, error) {
	return u.client.BuscarPartido(ctx, id)
}
