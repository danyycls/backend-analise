package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarRecursosRecebidosUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarRecursosRecebidosUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarRecursosRecebidosUseCase {
	return &BuscarRecursosRecebidosUseCase{client: c}
}

func (u *BuscarRecursosRecebidosUseCase) Buscar(ctx context.Context, filtro portaltransparencia.DespesaRecursosRecebidosQueryParams) ([]portaltransparencia.PessoaRecursosRecebidosUGMesDesnormalizada, error) {
	return u.client.ListarRecursosRecebidos(ctx, filtro)
}
