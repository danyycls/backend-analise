package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type ListarProcessosUseCase struct {
	client *senado.SenadoClient
}

func NovoListarProcessosUseCase(c *senado.SenadoClient) *ListarProcessosUseCase {
	return &ListarProcessosUseCase{client: c}
}

func (u *ListarProcessosUseCase) Listar(ctx context.Context, params map[string]string) ([]senado.ProcessoItem, error) {
	return u.client.ListarProcessos(ctx, params)
}
