package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/ibge"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualListarEstadosRequest struct{}

type EsferaEstadualListarEstadosResponse struct {
	Estados []types.EstadoIBGE
}

type EsferaEstadualListarEstadosUseCase struct {
	ibgeClient *ibge.IBGEClient
}

func NovoEsferaEstadualListarEstadosUseCase(ibge *ibge.IBGEClient) *EsferaEstadualListarEstadosUseCase {
	return &EsferaEstadualListarEstadosUseCase{
		ibgeClient: ibge,
	}
}

func (u *EsferaEstadualListarEstadosUseCase) Executar(ctx context.Context, req *EsferaEstadualListarEstadosRequest) (*EsferaEstadualListarEstadosResponse, error) {
	estados, err := u.ibgeClient.ListarEstados(ctx)
	if err != nil {
		return nil, err
	}
	return &EsferaEstadualListarEstadosResponse{Estados: estados}, nil
}
