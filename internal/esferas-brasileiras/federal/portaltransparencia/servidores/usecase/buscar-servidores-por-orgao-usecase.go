package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarServidoresPorOrgaoUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarServidoresPorOrgaoUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarServidoresPorOrgaoUseCase {
	return &BuscarServidoresPorOrgaoUseCase{client: c}
}

func (u *BuscarServidoresPorOrgaoUseCase) Buscar(ctx context.Context, filtro portaltransparencia.ServidorPorOrgaoQueryParams) ([]portaltransparencia.ServidorPorOrgao, error) {
	return u.client.ListarServidoresPorOrgao(ctx, filtro)
}
