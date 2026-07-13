package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarDocumentosPorFavorecidoUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarDocumentosPorFavorecidoUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarDocumentosPorFavorecidoUseCase {
	return &BuscarDocumentosPorFavorecidoUseCase{client: c}
}

func (u *BuscarDocumentosPorFavorecidoUseCase) Buscar(ctx context.Context, filtro portaltransparencia.DespesaDocumentosPorFavorecidoQueryParams) ([]interface{}, error) {
	return u.client.ListarDocumentosPorFavorecido(ctx, filtro)
}
