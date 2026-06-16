package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarEmendasUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarEmendasUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarEmendasUseCase {
	return &BuscarEmendasUseCase{client: c}
}

func (u *BuscarEmendasUseCase) Buscar(ctx context.Context, filtro portaltransparencia.EmendaQueryParams) ([]portaltransparencia.ConsultaEmendas, error) {
	return u.client.ListarEmendas(ctx, filtro)
}
