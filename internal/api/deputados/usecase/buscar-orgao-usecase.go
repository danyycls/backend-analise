package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalBuscarOrgaoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarOrgaoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarOrgaoUseCase {
	return &EsferaFederalBuscarOrgaoUseCase{client: c}
}

func (u *EsferaFederalBuscarOrgaoUseCase) Executar(ctx context.Context, id int) (*deputados.OrgaoDetalhe, error) {
	return u.client.BuscarOrgaoCamara(ctx, id)
}
