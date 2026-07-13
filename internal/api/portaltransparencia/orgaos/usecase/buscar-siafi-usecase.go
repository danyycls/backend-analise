package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarOrgaosSIAFIUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarOrgaosSIAFIUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarOrgaosSIAFIUseCase {
	return &BuscarOrgaosSIAFIUseCase{client: c}
}

func (u *BuscarOrgaosSIAFIUseCase) Buscar(ctx context.Context, filtro portaltransparencia.OrgaoQueryParams) ([]portaltransparencia.OrgaoSIAFI, error) {
	return u.client.ListarOrgaosSIAFI(ctx, filtro)
}
