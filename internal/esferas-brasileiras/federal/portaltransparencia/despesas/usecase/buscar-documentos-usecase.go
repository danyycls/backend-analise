package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarDocumentosUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarDocumentosUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarDocumentosUseCase {
	return &BuscarDocumentosUseCase{client: c}
}

func (u *BuscarDocumentosUseCase) Buscar(ctx context.Context, filtro portaltransparencia.DespesaDocumentosQueryParams) ([]interface{}, error) {
	return u.client.ListarDocumentos(ctx, filtro)
}
