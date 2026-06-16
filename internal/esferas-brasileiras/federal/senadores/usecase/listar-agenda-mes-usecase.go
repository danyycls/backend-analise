package usecase

import (
	"context"

	senado "github.com/danyele/laceu/internal/shared/clients/senado"
)

type ListarAgendaMesUseCase struct {
	client *senado.SenadoClient
}

func NovoListarAgendaMesUseCase(c *senado.SenadoClient) *ListarAgendaMesUseCase {
	return &ListarAgendaMesUseCase{client: c}
}

func (u *ListarAgendaMesUseCase) Listar(ctx context.Context, data string, params map[string]string) ([]senado.Reuniao, error) {
	return u.client.ListarAgendaMes(ctx, data, params)
}
