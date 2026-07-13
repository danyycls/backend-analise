package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarHistoricoItemEmpenhoUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarHistoricoItemEmpenhoUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarHistoricoItemEmpenhoUseCase {
	return &BuscarHistoricoItemEmpenhoUseCase{client: c}
}

func (u *BuscarHistoricoItemEmpenhoUseCase) Buscar(ctx context.Context, codigoDocumento string, sequencial, pagina int) ([]portaltransparencia.HistoricoSubItemEmpenho, error) {
	return u.client.ListarHistoricoItemEmpenho(ctx, codigoDocumento, sequencial, pagina)
}
