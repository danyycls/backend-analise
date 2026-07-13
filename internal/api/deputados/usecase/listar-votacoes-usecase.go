package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalListarVotacoesUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalListarVotacoesUseCase(c *deputados.DeputadosClient) *EsferaFederalListarVotacoesUseCase {
	return &EsferaFederalListarVotacoesUseCase{client: c}
}

func (u *EsferaFederalListarVotacoesUseCase) Executar(ctx context.Context, params map[string]string) ([]deputados.Votacao, error) {
	return u.client.ListarVotacoes(ctx, params)
}
