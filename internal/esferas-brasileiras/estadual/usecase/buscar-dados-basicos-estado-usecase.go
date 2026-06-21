package usecase

import (
	"context"
	"sort"

	"github.com/danyele/podp/internal/shared/clients/ibge"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualBuscarDadosBasicosEstadoRequest struct {
	UF string
}

type EsferaEstadualBuscarDadosBasicosEstadoResponse struct {
	Dados *types.DadosEstadoConsolidado
}

type EsferaEstadualBuscarDadosBasicosEstadoUseCase struct {
	ibgeClient *ibge.IBGEClient
}

func NovoEsferaEstadualBuscarDadosBasicosEstadoUseCase(ibge *ibge.IBGEClient) *EsferaEstadualBuscarDadosBasicosEstadoUseCase {
	return &EsferaEstadualBuscarDadosBasicosEstadoUseCase{
		ibgeClient: ibge,
	}
}

func (u *EsferaEstadualBuscarDadosBasicosEstadoUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarDadosBasicosEstadoRequest) (*EsferaEstadualBuscarDadosBasicosEstadoResponse, error) {
	nomeEstado := u.buscarNomeEstado(ctx, req.UF)

	result := &types.DadosEstadoConsolidado{
		UF:   req.UF,
		Nome: nomeEstado,
	}

	u.buscarMunicipiosComPopulacao(ctx, req.UF, result)
	u.ordenarMunicipiosPorPopulacao(result)

	return &EsferaEstadualBuscarDadosBasicosEstadoResponse{Dados: result}, nil
}

func (u *EsferaEstadualBuscarDadosBasicosEstadoUseCase) buscarNomeEstado(ctx context.Context, uf string) string {
	estados, err := u.ibgeClient.ListarEstados(ctx)
	if err != nil {
		return uf
	}
	for _, e := range estados {
		if e.Sigla == uf {
			return e.Nome
		}
	}
	return uf
}

func (u *EsferaEstadualBuscarDadosBasicosEstadoUseCase) buscarMunicipiosComPopulacao(ctx context.Context, uf string, result *types.DadosEstadoConsolidado) {
	municipiosIBGE, err := u.ibgeClient.ListarMunicipios(ctx, uf)
	if err != nil {
		return
	}

	ids := make([]int, len(municipiosIBGE))
	for i, m := range municipiosIBGE {
		ids[i] = m.ID
	}

	populacaoMap, _ := u.ibgeClient.BuscarPopulacao(ctx, ids)

	var totalPop int64
	for _, m := range municipiosIBGE {
		pop := populacaoMap[m.ID]
		totalPop += pop
		result.Municipios = append(result.Municipios, types.MunicipioComDados{
			ID:        m.ID,
			Nome:      m.Nome,
			Populacao: pop,
		})
	}
	result.Populacao = totalPop
}

func (u *EsferaEstadualBuscarDadosBasicosEstadoUseCase) ordenarMunicipiosPorPopulacao(result *types.DadosEstadoConsolidado) {
	sort.Slice(result.Municipios, func(i, j int) bool {
		return result.Municipios[i].Populacao < result.Municipios[j].Populacao
	})
}
