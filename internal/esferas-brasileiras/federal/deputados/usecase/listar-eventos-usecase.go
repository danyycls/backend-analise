package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalListarEventosUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarEventosUseCase(c *deputados.DeputadosClient) *EsferaFederalListarEventosUseCase {
	return &EsferaFederalListarEventosUseCase{client: c}
}

func (u *EsferaFederalListarEventosUseCase) Executar(ctx context.Context, params map[string]string) ([]deputados.Evento, error) {
	return u.client.ListarEventos(ctx, params)
}
