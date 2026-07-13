package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/ibge/client"
)

type BuscarPopulacaoRequest struct {
	MunicipioIDs []int `json:"municipio_ids"`
}

type BuscarPopulacaoUseCase struct {
	client *ibge.IBGEClient
}

func NovoBuscarPopulacaoUseCase(client *ibge.IBGEClient) *BuscarPopulacaoUseCase {
	return &BuscarPopulacaoUseCase{client: client}
}

func (u *BuscarPopulacaoUseCase) Executar(ctx context.Context, municipioIDs []int) (map[int]int64, error) {
	return u.client.BuscarPopulacao(ctx, municipioIDs)
}
