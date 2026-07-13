package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalBuscarLegislaturaUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarLegislaturaUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarLegislaturaUseCase {
	return &EsferaFederalBuscarLegislaturaUseCase{client: c}
}

func (u *EsferaFederalBuscarLegislaturaUseCase) Executar(ctx context.Context, id int) (*deputados.Legislatura, error) {
	return u.client.BuscarLegislatura(ctx, id)
}
