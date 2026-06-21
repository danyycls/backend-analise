package usecase

import (
	"context"
	"fmt"

	deputados "github.com/danyele/podp/internal/shared/clients/deputados"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualBuscarDeputadosRequest struct {
	UF string
}

type EsferaEstadualBuscarDeputadosResponse struct {
	Deputados []types.DeputadoUF
}

type EsferaEstadualBuscarDeputadosUseCase struct {
	deputadosCli *deputados.DeputadosClient
}

func NovoEsferaEstadualBuscarDeputadosUseCase(deputadosCli *deputados.DeputadosClient) *EsferaEstadualBuscarDeputadosUseCase {
	return &EsferaEstadualBuscarDeputadosUseCase{
		deputadosCli: deputadosCli,
	}
}

func (u *EsferaEstadualBuscarDeputadosUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarDeputadosRequest) (*EsferaEstadualBuscarDeputadosResponse, error) {
	depParams := map[string]string{"siglaUf": req.UF}
	deputadosAPI, err := u.deputadosCli.ListarInfoDeputadosAtivos(ctx, depParams)
	if err != nil {
		return nil, fmt.Errorf("erro buscar deputados: %w", err)
	}

	result := make([]types.DeputadoUF, 0, len(deputadosAPI))
	for _, d := range deputadosAPI {
		result = append(result, types.DeputadoUF{
			ID:            d.ID,
			Nome:          d.Nome,
			SiglaPartido:  d.SiglaPartido,
			SiglaUF:       d.SiglaUF,
			URLFoto:       d.URLFoto,
			Email:         d.Email,
			NomeEleitoral: "",
		})
	}
	return &EsferaEstadualBuscarDeputadosResponse{Deputados: result}, nil
}
