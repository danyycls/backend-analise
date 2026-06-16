package usecase

import (
	"context"

	"github.com/danyele/laceu/internal/shared/clients/deputados"
)

type EsferaFederalBuscarOrgaoAssociadoDeputadoRequest struct {
	ID     int
	Params map[string]string
}

type EsferaFederalBuscarOrgaoAssociadoDeputadoResponse struct {
	Orgaos []deputados.DeputadoOrgao
}

type EsferaFederalBuscarOrgaoAssociadoDeputadoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarOrgaoAssociadoDeputadoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarOrgaoAssociadoDeputadoUseCase {
	return &EsferaFederalBuscarOrgaoAssociadoDeputadoUseCase{client: c}
}

func (u *EsferaFederalBuscarOrgaoAssociadoDeputadoUseCase) Executar(ctx context.Context, req *EsferaFederalBuscarOrgaoAssociadoDeputadoRequest) (*EsferaFederalBuscarOrgaoAssociadoDeputadoResponse, error) {
	orgaos, err := u.client.ListarOrgaos(ctx, req.ID, req.Params)
	if err != nil {
		return nil, err
	}
	return &EsferaFederalBuscarOrgaoAssociadoDeputadoResponse{Orgaos: orgaos}, nil
}
