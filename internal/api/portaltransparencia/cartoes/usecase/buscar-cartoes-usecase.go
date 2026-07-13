package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
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
