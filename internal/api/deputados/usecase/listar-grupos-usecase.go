package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalListarGruposUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarGruposUseCase(c *deputados.DeputadosClient) *EsferaFederalListarGruposUseCase {
	return &EsferaFederalListarGruposUseCase{client: c}
}

func (u *EsferaFederalListarGruposUseCase) Executar(ctx context.Context) ([]deputados.Grupo, error) {
	return u.client.ListarGrupos(ctx)
}
