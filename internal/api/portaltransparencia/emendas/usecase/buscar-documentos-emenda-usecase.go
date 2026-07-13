package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarDocumentosEmendaUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarDocumentosEmendaUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarDocumentosEmendaUseCase {
	return &BuscarDocumentosEmendaUseCase{client: c}
}

func (u *BuscarDocumentosEmendaUseCase) Buscar(ctx context.Context, codigo string, pagina int) ([]portaltransparencia.DocumentoRelacionadoEmenda, error) {
	return u.client.ListarDocumentosEmenda(ctx, codigo, pagina)
}
