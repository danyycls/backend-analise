package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarSubfuncoesUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarSubfuncoesUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarSubfuncoesUseCase {
	return &BuscarSubfuncoesUseCase{client: c}
}

func (u *BuscarSubfuncoesUseCase) Buscar(ctx context.Context, filtro portaltransparencia.ListarFuncionalProgramaticaQueryParams) ([]portaltransparencia.Subfuncao, error) {
	return u.client.ListarSubfuncoes(ctx, filtro)
}
