package usecase

import (
	"context"

	client "github.com/danyele/laceu/internal/shared/clients/tcu"
)

type InabilitadosUseCase struct {
	client *client.TCUClient
}

func NovoInabilitadosUseCase(c *client.TCUClient) *InabilitadosUseCase {
	return &InabilitadosUseCase{client: c}
}

func (u *InabilitadosUseCase) Buscar(ctx context.Context, filter client.TCUQueryParams) ([]client.Sancoes, error) {
	return u.client.BuscarInabilitados(ctx, filter)
}
