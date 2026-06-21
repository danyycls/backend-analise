package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarAcoesUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarAcoesUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarAcoesUseCase {
	return &BuscarAcoesUseCase{client: c}
}

func (u *BuscarAcoesUseCase) Buscar(ctx context.Context, filtro portaltransparencia.ListarFuncionalProgramaticaQueryParams) ([]portaltransparencia.CodigoDescricao, error) {
	return u.client.ListarAcoes(ctx, filtro)
}
