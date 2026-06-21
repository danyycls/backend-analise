package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type BuscarFavorecidosFinaisPorDocumentoUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarFavorecidosFinaisPorDocumentoUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarFavorecidosFinaisPorDocumentoUseCase {
	return &BuscarFavorecidosFinaisPorDocumentoUseCase{client: c}
}

func (u *BuscarFavorecidosFinaisPorDocumentoUseCase) Buscar(ctx context.Context, codigoDocumento string, pagina int) ([]portaltransparencia.ConsultaFavorecidosFinaisPorDocumento, error) {
	return u.client.ListarFavorecidosFinaisPorDocumento(ctx, codigoDocumento, pagina)
}
