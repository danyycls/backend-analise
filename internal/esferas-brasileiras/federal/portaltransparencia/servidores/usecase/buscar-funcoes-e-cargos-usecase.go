package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarFuncoesECargosUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarFuncoesECargosUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarFuncoesECargosUseCase {
	return &BuscarFuncoesECargosUseCase{client: c}
}

func (u *BuscarFuncoesECargosUseCase) Buscar(ctx context.Context, filtro portaltransparencia.FuncaoCargoQueryParams) ([]portaltransparencia.FuncaoServidor, error) {
	return u.client.ListarFuncoesECargos(ctx, filtro)
}
