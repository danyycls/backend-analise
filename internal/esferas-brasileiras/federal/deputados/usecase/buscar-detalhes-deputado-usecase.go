package usecase

import (
	"context"

	"github.com/danyele/podp/internal/shared/clients/deputados"
)

type EsferaFederalBuscarDetalhesDeputadoRequest struct {
	ID int
}

type EsferaFederalBuscarDetalhesDeputadoResponse struct {
	Deputado *deputados.DeputadoResponse
}

type EsferaFederalBuscarDetalhesDeputadoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarDetalhesDeputadoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarDetalhesDeputadoUseCase {
	return &EsferaFederalBuscarDetalhesDeputadoUseCase{client: c}
}

func (u *EsferaFederalBuscarDetalhesDeputadoUseCase) Executar(ctx context.Context, req *EsferaFederalBuscarDetalhesDeputadoRequest) (*EsferaFederalBuscarDetalhesDeputadoResponse, error) {
	deputado, err := u.client.BuscarDeputado(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	result := &deputados.DeputadoResponse{
		Deputado:         deputado,
		Frentes:          []deputados.Frente{},
		Historico:        []deputados.DeputadoHistorico{},
		MandatosExternos: []deputados.DeputadoMandatoExterno{},
	}

	if frentes, err := u.client.ListarFrentesDeputado(ctx, req.ID); err == nil {
		result.Frentes = frentes
	}

	if historico, err := u.client.ListarHistorico(ctx, req.ID); err == nil {
		result.Historico = historico
	}

	if mandatos, err := u.client.ListarMandatosExternos(ctx, req.ID); err == nil {
		result.MandatosExternos = mandatos
	}

	return &EsferaFederalBuscarDetalhesDeputadoResponse{Deputado: result}, nil
}
