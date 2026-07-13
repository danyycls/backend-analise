package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type ListarTramitacoesUseCase struct {
	client *deputados.DeputadosClient
}

func NovoListarTramitacoesUseCase(c *deputados.DeputadosClient) *ListarTramitacoesUseCase {
	return &ListarTramitacoesUseCase{client: c}
}

func (u *ListarTramitacoesUseCase) Executar(ctx context.Context, idProposicao int) ([]deputados.Tramitacao, error) {
	return u.client.ListarTramitacoes(ctx, idProposicao)
}
