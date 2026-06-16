package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarFuncoesUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarFuncoesUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarFuncoesUseCase {
	return &BuscarFuncoesUseCase{client: c}
}

func (u *BuscarFuncoesUseCase) Buscar(ctx context.Context, filtro portaltransparencia.ListarFuncionalProgramaticaQueryParams) ([]portaltransparencia.Funcao, error) {
	return u.client.ListarFuncoes(ctx, filtro)
}
