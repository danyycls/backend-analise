package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarDespesasPorFuncionalProgramaticaUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarDespesasPorFuncionalProgramaticaUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarDespesasPorFuncionalProgramaticaUseCase {
	return &BuscarDespesasPorFuncionalProgramaticaUseCase{client: c}
}

func (u *BuscarDespesasPorFuncionalProgramaticaUseCase) Buscar(ctx context.Context, filtro portaltransparencia.DespesaFuncionalProgramaticaQueryParams) ([]portaltransparencia.DespesaAnualPorFuncaoESubfuncao, error) {
	return u.client.ListarDespesasPorFuncionalProgramatica(ctx, filtro)
}
