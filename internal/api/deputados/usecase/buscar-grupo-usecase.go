package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalBuscarGrupoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarGrupoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarGrupoUseCase {
	return &EsferaFederalBuscarGrupoUseCase{client: c}
}

func (u *EsferaFederalBuscarGrupoUseCase) Executar(ctx context.Context, id int) (*deputados.GrupoDetalhe, error) {
	return u.client.BuscarGrupo(ctx, id)
}
