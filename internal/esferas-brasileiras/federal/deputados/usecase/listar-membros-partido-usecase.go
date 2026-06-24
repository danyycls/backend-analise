package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalListarMembrosPartidoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarMembrosPartidoUseCase(c *deputados.DeputadosClient) *EsferaFederalListarMembrosPartidoUseCase {
	return &EsferaFederalListarMembrosPartidoUseCase{client: c}
}

func (u *EsferaFederalListarMembrosPartidoUseCase) Executar(ctx context.Context, id int, params map[string]string) ([]deputados.Deputado, error) {
	return u.client.ListarMembrosPartido(ctx, id, params)
}
