package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/shared/clients/senado"
)

type BuscarEncontroUseCase struct {
	client *senado.SenadoClient
}

func NovoBuscarEncontroUseCase(c *senado.SenadoClient) *BuscarEncontroUseCase {
	return &BuscarEncontroUseCase{client: c}
}

func (u *BuscarEncontroUseCase) Buscar(ctx context.Context, codigo string, params map[string]string) (*senado.PlenarioEncontro, error) {
	return u.client.BuscarEncontro(ctx, codigo, params)
}
