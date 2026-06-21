package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/shared/clients/senado"
)

type ListarTodasComissoesUseCase struct {
	client *senado.SenadoClient
}

func NovoListarTodasComissoesUseCase(c *senado.SenadoClient) *ListarTodasComissoesUseCase {
	return &ListarTodasComissoesUseCase{client: c}
}

func (u *ListarTodasComissoesUseCase) Listar(ctx context.Context) ([]senado.ComissaoResumo, error) {
	return u.client.ListarTodasComissoes(ctx)
}
