package usecase

import (
	"context"

	"github.com/danyele/podp/internal/sources/deputados/client"
)

type EsferaFederalBuscarDeputadosAtivosRequest struct {
	Params map[string]string
}

type EsferaFederalBuscarDeputadosAtivosResponse struct {
	Deputados []deputados.Deputado
}

type EsferaFederalBuscarDeputadosAtivosUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarDeputadosAtivosUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarDeputadosAtivosUseCase {
	return &EsferaFederalBuscarDeputadosAtivosUseCase{client: c}
}

func (u *EsferaFederalBuscarDeputadosAtivosUseCase) Executar(ctx context.Context, req *EsferaFederalBuscarDeputadosAtivosRequest) (*EsferaFederalBuscarDeputadosAtivosResponse, error) {
	deputados, err := u.client.ListarInfoDeputadosAtivos(ctx, req.Params)
	if err != nil {
		return nil, err
	}
	return &EsferaFederalBuscarDeputadosAtivosResponse{Deputados: deputados}, nil
}
