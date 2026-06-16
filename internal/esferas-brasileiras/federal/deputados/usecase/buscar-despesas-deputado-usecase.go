package usecase

import (
	"context"
	"strings"

	"github.com/danyele/laceu/internal/shared/clients/deputados"
)

type EsferaFederalBuscarDespesasDeputadoRequest struct {
	ID                int
	Params            map[string]string
	TipoDespesa       string
	CNPJCPFFornecedor string
}

type EsferaFederalBuscarDespesasDeputadoResponse struct {
	Despesas []deputados.DeputadoDespesa
}

type EsferaFederalBuscarDespesasDeputadoUseCase struct {
	client *deputados.DeputadosClient
}

func NovoEsferaFederalBuscarDespesasDeputadoUseCase(c *deputados.DeputadosClient) *EsferaFederalBuscarDespesasDeputadoUseCase {
	return &EsferaFederalBuscarDespesasDeputadoUseCase{client: c}
}

func (u *EsferaFederalBuscarDespesasDeputadoUseCase) Executar(ctx context.Context, req *EsferaFederalBuscarDespesasDeputadoRequest) (*EsferaFederalBuscarDespesasDeputadoResponse, error) {
	params := req.Params

	var despesas []deputados.DeputadoDespesa
	var err error

	if req.TipoDespesa != "" || req.CNPJCPFFornecedor != "" {
		despesas, err = u.client.ListarTodasDespesasPorDeputado(ctx, req.ID, params)
	} else {
		despesas, err = u.client.ListarDespesasPorDeputado(ctx, req.ID, params)
	}
	if err != nil {
		return nil, err
	}

	if req.TipoDespesa == "" && req.CNPJCPFFornecedor == "" {
		return &EsferaFederalBuscarDespesasDeputadoResponse{Despesas: despesas}, nil
	}

	filtered := make([]deputados.DeputadoDespesa, 0, len(despesas))
	for _, d := range despesas {
		if req.TipoDespesa != "" && d.TipoDespesa != req.TipoDespesa {
			continue
		}
		if req.CNPJCPFFornecedor != "" {
			normalized := strings.NewReplacer(".", "", "-", "", "/", "").Replace(req.CNPJCPFFornecedor)
			saved := strings.NewReplacer(".", "", "-", "", "/", "").Replace(d.CNPJCPFFornecedor)
			if saved != normalized {
				continue
			}
		}
		filtered = append(filtered, d)
	}

	return &EsferaFederalBuscarDespesasDeputadoResponse{Despesas: filtered}, nil
}
