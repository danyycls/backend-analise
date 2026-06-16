package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarRemuneracaoServidoresUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarRemuneracaoServidoresUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarRemuneracaoServidoresUseCase {
	return &BuscarRemuneracaoServidoresUseCase{client: c}
}

func (u *BuscarRemuneracaoServidoresUseCase) Buscar(ctx context.Context, filtro portaltransparencia.ServidorRemuneracaoQueryParams) ([]portaltransparencia.ServidorRemuneracao, error) {
	return u.client.ListarRemuneracaoServidores(ctx, filtro)
}
