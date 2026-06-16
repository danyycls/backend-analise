package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarMovimentacaoLiquidaUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarMovimentacaoLiquidaUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarMovimentacaoLiquidaUseCase {
	return &BuscarMovimentacaoLiquidaUseCase{client: c}
}

func (u *BuscarMovimentacaoLiquidaUseCase) Buscar(ctx context.Context, filtro portaltransparencia.DespesaMovimentacaoLiquidaQueryParams) ([]portaltransparencia.DespesaLiquidaAnualPorFuncaoESubfuncao, error) {
	return u.client.ListarDespesasMovimentacaoLiquida(ctx, filtro)
}
