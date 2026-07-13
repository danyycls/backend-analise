package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarDespesasPorOrgaoUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarDespesasPorOrgaoUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarDespesasPorOrgaoUseCase {
	return &BuscarDespesasPorOrgaoUseCase{client: c}
}

func (u *BuscarDespesasPorOrgaoUseCase) Buscar(ctx context.Context, filtro portaltransparencia.DespesaPorOrgaoQueryParams) ([]portaltransparencia.DespesaAnualPorOrgao, error) {
	return u.client.ListarDespesasPorOrgao(ctx, filtro)
}
