package usecase

import (
	"context"

	client "github.com/danyele/podp/internal/shared/clients/tcu"
)

type ContasIrregularesUseCase struct {
	client *client.TCUClient
}

func NovoContasIrregularesUseCase(c *client.TCUClient) *ContasIrregularesUseCase {
	return &ContasIrregularesUseCase{client: c}
}

func (u *ContasIrregularesUseCase) Buscar(ctx context.Context, filter client.TCUQueryParams) ([]client.ContasIrregulares, error) {
	return u.client.BuscarContasIrregulares(ctx, filter)
}
