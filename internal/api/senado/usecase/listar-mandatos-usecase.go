package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type ListarMandatosUseCase struct {
	client *senado.SenadoClient
}

func NovoListarMandatosUseCase(c *senado.SenadoClient) *ListarMandatosUseCase {
	return &ListarMandatosUseCase{client: c}
}

func (u *ListarMandatosUseCase) Listar(ctx context.Context, codigo string) ([]senado.MandatoDetalhe, error) {
	return u.client.ListarMandatos(ctx, codigo)
}
