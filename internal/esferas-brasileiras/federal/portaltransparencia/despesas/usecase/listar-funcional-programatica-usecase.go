package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
)

type ListarFuncionalProgramaticaUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoListarFuncionalProgramaticaUseCase(c *portaltransparencia.PortalTransparenciaClient) *ListarFuncionalProgramaticaUseCase {
	return &ListarFuncionalProgramaticaUseCase{client: c}
}

func (u *ListarFuncionalProgramaticaUseCase) Buscar(ctx context.Context, ano, pagina int) ([]portaltransparencia.FuncionalProgramatica, error) {
	return u.client.ListarFuncionalProgramatica(ctx, ano, pagina)
}
