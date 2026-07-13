package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type ListarMateriaTramitacaoUseCase struct {
	client *senado.SenadoClient
}

func NovoListarMateriaTramitacaoUseCase(c *senado.SenadoClient) *ListarMateriaTramitacaoUseCase {
	return &ListarMateriaTramitacaoUseCase{client: c}
}

func (u *ListarMateriaTramitacaoUseCase) Listar(ctx context.Context, params map[string]string) ([]senado.MateriaItem, error) {
	return u.client.ListarMateriaTramitacao(ctx, params)
}
