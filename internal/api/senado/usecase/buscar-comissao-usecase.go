package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type BuscarComissaoUseCase struct {
	client *senado.SenadoClient
}

func NovoBuscarComissaoUseCase(c *senado.SenadoClient) *BuscarComissaoUseCase {
	return &BuscarComissaoUseCase{client: c}
}

func (u *BuscarComissaoUseCase) Buscar(ctx context.Context, codigo string) (*senado.ComissaoDetalhe, error) {
	return u.client.BuscarComissao(ctx, codigo)
}
