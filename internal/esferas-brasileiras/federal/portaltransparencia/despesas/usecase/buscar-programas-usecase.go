package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarProgramasUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarProgramasUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarProgramasUseCase {
	return &BuscarProgramasUseCase{client: c}
}

func (u *BuscarProgramasUseCase) Buscar(ctx context.Context, filtro portaltransparencia.ListarFuncionalProgramaticaQueryParams) ([]portaltransparencia.CodigoDescricao, error) {
	return u.client.ListarProgramas(ctx, filtro)
}
