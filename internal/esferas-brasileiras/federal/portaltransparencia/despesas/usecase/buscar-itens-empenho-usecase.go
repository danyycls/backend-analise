package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarItensEmpenhoUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarItensEmpenhoUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarItensEmpenhoUseCase {
	return &BuscarItensEmpenhoUseCase{client: c}
}

func (u *BuscarItensEmpenhoUseCase) Buscar(ctx context.Context, codigoDocumento string, pagina int) ([]portaltransparencia.DetalhamentoDoGasto, error) {
	return u.client.ListarItensEmpenho(ctx, codigoDocumento, pagina)
}
