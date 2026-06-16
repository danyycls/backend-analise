package usecase

import (
	"context"

	senado "github.com/danyele/laceu/internal/shared/clients/senado"
)

type ListarAgendaDiaUseCase struct {
	client *senado.SenadoClient
}

func NovoListarAgendaDiaUseCase(c *senado.SenadoClient) *ListarAgendaDiaUseCase {
	return &ListarAgendaDiaUseCase{client: c}
}

func (u *ListarAgendaDiaUseCase) Listar(ctx context.Context, data string, params map[string]string) ([]senado.Reuniao, error) {
	return u.client.ListarAgendaDia(ctx, data, params)
}
