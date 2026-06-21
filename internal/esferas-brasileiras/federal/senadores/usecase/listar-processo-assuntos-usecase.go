package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/shared/clients/senado"
)

type ListarProcessoAssuntosUseCase struct {
	client *senado.SenadoClient
}

func NovoListarProcessoAssuntosUseCase(c *senado.SenadoClient) *ListarProcessoAssuntosUseCase {
	return &ListarProcessoAssuntosUseCase{client: c}
}

func (u *ListarProcessoAssuntosUseCase) Listar(ctx context.Context) ([]senado.ProcessoAssunto, error) {
	return u.client.ListarProcessoAssuntos(ctx)
}
