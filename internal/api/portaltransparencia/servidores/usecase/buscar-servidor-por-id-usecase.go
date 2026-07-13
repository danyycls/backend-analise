package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarServidorPorIDUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarServidorPorIDUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarServidorPorIDUseCase {
	return &BuscarServidorPorIDUseCase{client: c}
}

func (u *BuscarServidorPorIDUseCase) Buscar(ctx context.Context, id int) (*portaltransparencia.CadastroServidor, error) {
	return u.client.BuscarServidorPorID(ctx, id)
}
