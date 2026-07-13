package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/portaltransparencia/client"
)

type BuscarPessoasFisicasUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarPessoasFisicasUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarPessoasFisicasUseCase {
	return &BuscarPessoasFisicasUseCase{client: c}
}

func (u *BuscarPessoasFisicasUseCase) Buscar(ctx context.Context, filtro portaltransparencia.PessoaFisicaQueryParams) (*portaltransparencia.PessoaFisica, error) {
	return u.client.ListarPessoasFisicas(ctx, filtro)
}
