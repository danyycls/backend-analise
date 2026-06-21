package usecase

import (
	"context"

	client "github.com/danyele/podp/internal/shared/clients/tcu"
)

type InidoneosUseCase struct {
	client *client.TCUClient
}

func NovoInidoneosUseCase(c *client.TCUClient) *InidoneosUseCase {
	return &InidoneosUseCase{client: c}
}

func (u *InidoneosUseCase) Buscar(ctx context.Context, filter client.TCUQueryParams) ([]client.Sancoes, error) {
	return u.client.BuscarInidoneos(ctx, filter)
}
