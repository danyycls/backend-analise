package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/portaltransparencia"
)

type BuscarPessoasJuridicasUseCase struct {
	client *portaltransparencia.PortalTransparenciaClient
}

func NovoBuscarPessoasJuridicasUseCase(c *portaltransparencia.PortalTransparenciaClient) *BuscarPessoasJuridicasUseCase {
	return &BuscarPessoasJuridicasUseCase{client: c}
}

func (u *BuscarPessoasJuridicasUseCase) Buscar(ctx context.Context, filtro portaltransparencia.PessoaJuridicaQueryParams) (*portaltransparencia.PessoaJuridica, error) {
	return u.client.ListarPessoasJuridicas(ctx, filtro)
}
