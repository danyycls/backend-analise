package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalListarMembrosFrenteUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarMembrosFrenteUseCase(c *deputados.DeputadosClient) *EsferaFederalListarMembrosFrenteUseCase {
	return &EsferaFederalListarMembrosFrenteUseCase{client: c}
}

func (u *EsferaFederalListarMembrosFrenteUseCase) Executar(ctx context.Context, id int) ([]deputados.MembroFrente, error) {
	return u.client.ListarMembrosFrente(ctx, id)
}
