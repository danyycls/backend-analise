package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarDocumentosRelacionadosUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarDocumentosRelacionadosUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarDocumentosRelacionadosUseCase {
	return &BuscarDocumentosRelacionadosUseCase{client: c}
}

func (u *BuscarDocumentosRelacionadosUseCase) Buscar(ctx context.Context, codigoDocumento, fase string) ([]portaltransparencia.DocumentoRelacionado, error) {
	return u.client.ListarDocumentosRelacionados(ctx, codigoDocumento, fase)
}
