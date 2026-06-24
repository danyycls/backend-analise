package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalListarPartidosUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarPartidosUseCase(c *deputados.DeputadosClient) *EsferaFederalListarPartidosUseCase {
	return &EsferaFederalListarPartidosUseCase{client: c}
}

func (u *EsferaFederalListarPartidosUseCase) Executar(ctx context.Context, params map[string]string) ([]deputados.Partido, error) {
	return u.client.ListarPartidos(ctx, params)
}
