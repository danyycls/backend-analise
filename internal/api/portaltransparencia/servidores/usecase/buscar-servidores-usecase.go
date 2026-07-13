package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarServidoresUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarServidoresUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarServidoresUseCase {
	return &BuscarServidoresUseCase{client: c}
}

func (u *BuscarServidoresUseCase) Buscar(ctx context.Context, filtro portaltransparencia.ServidorQueryParams) ([]portaltransparencia.CadastroServidor, error) {
	return u.client.ListarServidores(ctx, filtro)
}
