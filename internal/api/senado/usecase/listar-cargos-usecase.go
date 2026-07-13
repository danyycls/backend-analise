package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type ListarCargosUseCase struct {
	client *senado.SenadoClient
}

func NovoListarCargosUseCase(c *senado.SenadoClient) *ListarCargosUseCase {
	return &ListarCargosUseCase{client: c}
}

func (u *ListarCargosUseCase) Listar(ctx context.Context, codigo string) ([]senado.Cargo, error) {
	return u.client.ListarCargos(ctx, codigo)
}
