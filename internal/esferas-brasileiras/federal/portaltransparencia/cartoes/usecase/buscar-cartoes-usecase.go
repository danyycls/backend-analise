package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarCartoesUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarCartoesUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarCartoesUseCase {
	return &BuscarCartoesUseCase{client: c}
}

func (u *BuscarCartoesUseCase) Buscar(ctx context.Context, filtro portaltransparencia.CartaoQueryParams) ([]portaltransparencia.Cartao, error) {
	return u.client.ListarCartoes(ctx, filtro)
}
