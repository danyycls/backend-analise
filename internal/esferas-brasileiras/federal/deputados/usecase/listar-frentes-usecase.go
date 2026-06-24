package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalListarFrentesUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarFrentesUseCase(c *deputados.DeputadosClient) *EsferaFederalListarFrentesUseCase {
	return &EsferaFederalListarFrentesUseCase{client: c}
}

func (u *EsferaFederalListarFrentesUseCase) Executar(ctx context.Context, idLegislatura int) ([]deputados.Frente, error) {
	return u.client.ListarFrentes(ctx, idLegislatura)
}
