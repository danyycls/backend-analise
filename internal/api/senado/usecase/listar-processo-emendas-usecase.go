package usecase

import (
	"context"

	senado "github.com/danyele/podp/internal/sources/senado/client"
)

type ListarProcessoEmendasUseCase struct {
	client *senado.SenadoClient
}

func NovoListarProcessoEmendasUseCase(c *senado.SenadoClient) *ListarProcessoEmendasUseCase {
	return &ListarProcessoEmendasUseCase{client: c}
}

func (u *ListarProcessoEmendasUseCase) Listar(ctx context.Context, params map[string]string) ([]senado.ProcessoEmenda, error) {
	return u.client.ListarProcessoEmendas(ctx, params)
}
