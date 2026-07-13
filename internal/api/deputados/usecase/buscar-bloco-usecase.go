package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalBuscarBlocoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarBlocoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarBlocoUseCase {
	return &EsferaFederalBuscarBlocoUseCase{client: c}
}

func (u *EsferaFederalBuscarBlocoUseCase) Executar(ctx context.Context, id string) (*deputados.Bloco, error) {
	return u.client.BuscarBloco(ctx, id)
}
