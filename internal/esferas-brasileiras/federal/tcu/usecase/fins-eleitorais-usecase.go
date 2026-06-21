package usecase

import (
	"context"

	client "github.com/danyele/podp/internal/shared/clients/tcu"
)

type FinsEleitoraisUseCase struct {
	client *client.TCUClient
}

func NovoFinsEleitoraisUseCase(c *client.TCUClient) *FinsEleitoraisUseCase {
	return &FinsEleitoraisUseCase{client: c}
}

func (u *FinsEleitoraisUseCase) Buscar(ctx context.Context, filter client.TCUQueryParams) ([]client.FinsEleitorais, error) {
	return u.client.BuscarFinsEleitorais(ctx, filter)
}
