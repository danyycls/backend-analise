package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarPEPsUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarPEPsUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarPEPsUseCase {
	return &BuscarPEPsUseCase{client: c}
}

func (u *BuscarPEPsUseCase) Buscar(ctx context.Context, filtro portaltransparencia.PEPQueryParams) ([]portaltransparencia.PEP, error) {
	return u.client.ListarPEPs(ctx, filtro)
}
