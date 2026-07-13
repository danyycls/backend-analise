package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/sources/ibge/client"
)

type ListarEstadosUseCase struct {
	client *ibge.IBGEClient
}

func NovoListarEstadosUseCase(client *ibge.IBGEClient) *ListarEstadosUseCase {
	return &ListarEstadosUseCase{client: client}
}

func (u *ListarEstadosUseCase) Executar(ctx context.Context) ([]types.EstadoIBGE, error) {
	return u.client.ListarEstados(ctx)
}
