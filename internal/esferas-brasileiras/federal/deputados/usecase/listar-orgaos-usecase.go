package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalListarOrgaosUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarOrgaosUseCase(c *deputados.DeputadosClient) *EsferaFederalListarOrgaosUseCase {
	return &EsferaFederalListarOrgaosUseCase{client: c}
}

func (u *EsferaFederalListarOrgaosUseCase) Executar(ctx context.Context, params map[string]string) ([]deputados.Orgao, error) {
	return u.client.ListarOrgaosCamara(ctx, params)
}
