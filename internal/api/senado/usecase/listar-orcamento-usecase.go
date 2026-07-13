package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type ListarOrcamentoUseCase struct {
	client *senado.SenadoClient
}

func NovoListarOrcamentoUseCase(c *senado.SenadoClient) *ListarOrcamentoUseCase {
	return &ListarOrcamentoUseCase{client: c}
}

func (u *ListarOrcamentoUseCase) Listar(ctx context.Context) ([]senado.LoteEmendasOrcamento, error) {
	return u.client.ListarOrcamento(ctx)
}
