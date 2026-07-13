package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarDocumentoPorCodigoUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarDocumentoPorCodigoUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarDocumentoPorCodigoUseCase {
	return &BuscarDocumentoPorCodigoUseCase{client: c}
}

func (u *BuscarDocumentoPorCodigoUseCase) Buscar(ctx context.Context, codigo string) (*portaltransparencia.DespesasPorDocumento, error) {
	return u.client.BuscarDocumentoPorCodigo(ctx, codigo)
}
