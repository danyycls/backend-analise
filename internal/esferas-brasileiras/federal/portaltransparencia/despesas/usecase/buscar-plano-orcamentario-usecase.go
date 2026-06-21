package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarPlanoOrcamentarioUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarPlanoOrcamentarioUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarPlanoOrcamentarioUseCase {
	return &BuscarPlanoOrcamentarioUseCase{client: c}
}

func (u *BuscarPlanoOrcamentarioUseCase) Buscar(ctx context.Context, filtro portaltransparencia.DespesaPlanoOrcamentarioQueryParams) ([]portaltransparencia.DespesasPorPlanoOrcamentario, error) {
	return u.client.ListarDespesasPlanoOrcamentario(ctx, filtro)
}
