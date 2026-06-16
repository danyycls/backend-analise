package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarTiposTransferenciaUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarTiposTransferenciaUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarTiposTransferenciaUseCase {
	return &BuscarTiposTransferenciaUseCase{client: c}
}

func (u *BuscarTiposTransferenciaUseCase) Buscar(ctx context.Context) ([]portaltransparencia.CodigoDescricao, error) {
	return u.client.ListarTiposTransferencia(ctx)
}
