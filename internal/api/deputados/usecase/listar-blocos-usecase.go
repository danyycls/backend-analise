package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalListarBlocosUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarBlocosUseCase(c *deputados.DeputadosClient) *EsferaFederalListarBlocosUseCase {
	return &EsferaFederalListarBlocosUseCase{client: c}
}

func (u *EsferaFederalListarBlocosUseCase) Executar(ctx context.Context, params map[string]string) ([]deputados.Bloco, error) {
	return u.client.ListarBlocos(ctx, params)
}
