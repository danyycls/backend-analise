package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarOrgaosSIAPEUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarOrgaosSIAPEUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarOrgaosSIAPEUseCase {
	return &BuscarOrgaosSIAPEUseCase{client: c}
}

func (u *BuscarOrgaosSIAPEUseCase) Buscar(ctx context.Context, filtro portaltransparencia.OrgaoQueryParams) ([]portaltransparencia.OrgaoSIAPE, error) {
	return u.client.ListarOrgaosSIAPE(ctx, filtro)
}
