package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalListarPartidosDoBlocoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarPartidosDoBlocoUseCase(c *deputados.DeputadosClient) *EsferaFederalListarPartidosDoBlocoUseCase {
	return &EsferaFederalListarPartidosDoBlocoUseCase{client: c}
}

func (u *EsferaFederalListarPartidosDoBlocoUseCase) Executar(ctx context.Context, id string) ([]deputados.Partido, error) {
	return u.client.ListarPartidosDoBloco(ctx, id)
}
