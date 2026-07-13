package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type BuscarProcessoUseCase struct {
	client *senado.SenadoClient
}

func NovoBuscarProcessoUseCase(c *senado.SenadoClient) *BuscarProcessoUseCase {
	return &BuscarProcessoUseCase{client: c}
}

func (u *BuscarProcessoUseCase) Buscar(ctx context.Context, id string) (*senado.ProcessoItem, error) {
	return u.client.BuscarProcesso(ctx, id)
}
