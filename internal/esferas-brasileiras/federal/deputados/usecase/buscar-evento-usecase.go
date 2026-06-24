package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalBuscarEventoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarEventoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarEventoUseCase {
	return &EsferaFederalBuscarEventoUseCase{client: c}
}

func (u *EsferaFederalBuscarEventoUseCase) Executar(ctx context.Context, id int) (*deputados.EventoDetalhe, error) {
	return u.client.BuscarEvento(ctx, id)
}
