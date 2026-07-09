package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/ibge"
	"github.com/danyele/podp/internal/shared/types"
)

type ListarMunicipiosUseCase struct {
	client *ibge.IBGEClient
}

func NovoListarMunicipiosUseCase(client *ibge.IBGEClient) *ListarMunicipiosUseCase {
	return &ListarMunicipiosUseCase{client: client}
}

func (u *ListarMunicipiosUseCase) Executar(ctx context.Context, uf string) ([]types.MunicipioIBGE, error) {
	return u.client.ListarMunicipios(ctx, uf)
}
