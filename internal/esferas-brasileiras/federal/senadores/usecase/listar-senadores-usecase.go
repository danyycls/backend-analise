package usecase

import (
	"context"

	senado "github.com/danyele/laceu/internal/shared/clients/senado"
)

type ListarSenadoresUseCase struct {
	client *senado.SenadoClient
}

func NovoListarSenadoresUseCase(c *senado.SenadoClient) *ListarSenadoresUseCase {
	return &ListarSenadoresUseCase{client: c}
}

func (u *ListarSenadoresUseCase) Listar(ctx context.Context) ([]senado.ParlamentarResumo, error) {
	return u.client.ListarSenadores(ctx)
}
