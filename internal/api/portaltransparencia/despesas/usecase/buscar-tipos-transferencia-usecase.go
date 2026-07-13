package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
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
