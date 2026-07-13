package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalListarMembrosOrgaoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarMembrosOrgaoUseCase(c *deputados.DeputadosClient) *EsferaFederalListarMembrosOrgaoUseCase {
	return &EsferaFederalListarMembrosOrgaoUseCase{client: c}
}

func (u *EsferaFederalListarMembrosOrgaoUseCase) Executar(ctx context.Context, id int, params map[string]string) ([]deputados.MembroOrgao, error) {
	return u.client.ListarMembrosOrgaoCamara(ctx, id, params)
}
