package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/shared/clients/senado"
)

type ListarComissoesUseCase struct {
	client *senado.SenadoClient
}

func NovoListarComissoesUseCase(c *senado.SenadoClient) *ListarComissoesUseCase {
	return &ListarComissoesUseCase{client: c}
}

func (u *ListarComissoesUseCase) Listar(ctx context.Context, codigo string) ([]senado.ComissaoMembro, error) {
	return u.client.ListarComissoes(ctx, codigo)
}
