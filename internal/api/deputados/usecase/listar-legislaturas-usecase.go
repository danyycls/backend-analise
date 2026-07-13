package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalListarLegislaturasUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarLegislaturasUseCase(c *deputados.DeputadosClient) *EsferaFederalListarLegislaturasUseCase {
	return &EsferaFederalListarLegislaturasUseCase{client: c}
}

func (u *EsferaFederalListarLegislaturasUseCase) Executar(ctx context.Context) ([]deputados.Legislatura, error) {
	return u.client.ListarLegislaturas(ctx)
}
