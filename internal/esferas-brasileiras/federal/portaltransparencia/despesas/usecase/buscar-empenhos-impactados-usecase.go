package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarEmpenhosImpactadosUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarEmpenhosImpactadosUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarEmpenhosImpactadosUseCase {
	return &BuscarEmpenhosImpactadosUseCase{client: c}
}

func (u *BuscarEmpenhosImpactadosUseCase) Buscar(ctx context.Context, codigoDocumento, fase string, pagina int) ([]portaltransparencia.EmpenhoImpactadoBasico, error) {
	return u.client.ListarEmpenhosImpactados(ctx, codigoDocumento, fase, pagina)
}
