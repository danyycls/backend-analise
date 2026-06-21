package usecase

import (
	"context"
	"fmt"
	"sort"

	"github.com/danyele/podp/internal/shared/clients/ibge"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/types"
)

type EsferaEstadualBuscarMunicipiosPopulacaoRequest struct {
	UF string
}

type EsferaEstadualBuscarMunicipiosPopulacaoResponse struct {
	Municipios []types.MunicipioComDados
}

type EsferaEstadualBuscarMunicipiosPopulacaoUseCase struct {
	ibgeClient *ibge.IBGEClient
}

func NovoEsferaEstadualBuscarMunicipiosPopulacaoUseCase(ibgeClient *ibge.IBGEClient) *EsferaEstadualBuscarMunicipiosPopulacaoUseCase {
	return &EsferaEstadualBuscarMunicipiosPopulacaoUseCase{
		ibgeClient: ibgeClient,
	}
}

func (u *EsferaEstadualBuscarMunicipiosPopulacaoUseCase) Executar(
	ctx context.Context,
	req *EsferaEstadualBuscarMunicipiosPopulacaoRequest,
) (*EsferaEstadualBuscarMunicipiosPopulacaoResponse, error) {
	log := logger.New("Estadual: UseCase: BuscarMunicipiosPopulacao")

	municipios, err := u.buscarMunicipios(ctx, req.UF)
	if err != nil {
		return nil, fmt.Errorf("erro IBGE municipios: %w", err)
	}

	populacaoMap := u.buscarPopulacaoPorIDs(ctx, municipios, log)

	result := u.montarResultadoMunicipios(municipios, populacaoMap)
	u.ordenarMunicipiosPorPopulacao(result)

	return &EsferaEstadualBuscarMunicipiosPopulacaoResponse{Municipios: result}, nil
}

func (u *EsferaEstadualBuscarMunicipiosPopulacaoUseCase) buscarMunicipios(ctx context.Context, uf string) ([]types.MunicipioIBGE, error) {
	return u.ibgeClient.ListarMunicipios(ctx, uf)
}

func (u *EsferaEstadualBuscarMunicipiosPopulacaoUseCase) buscarPopulacaoPorIDs(
	ctx context.Context,
	municipios []types.MunicipioIBGE,
	log *logger.Logger,
) map[int]int64 {
	ids := make([]int, len(municipios))
	for i, m := range municipios {
		ids[i] = m.ID
	}

	populacaoMap, err := u.ibgeClient.BuscarPopulacao(ctx, ids)
	if err != nil {
		log.Error("erro ao buscar populacao IBGE", "erro", err)
	}
	return populacaoMap
}

func (u *EsferaEstadualBuscarMunicipiosPopulacaoUseCase) montarResultadoMunicipios(
	municipios []types.MunicipioIBGE,
	populacaoMap map[int]int64,
) []types.MunicipioComDados {
	result := make([]types.MunicipioComDados, len(municipios))
	for i, m := range municipios {
		result[i] = types.MunicipioComDados{
			ID:        m.ID,
			Nome:      m.Nome,
			Populacao: populacaoMap[m.ID],
		}
	}
	return result
}

func (u *EsferaEstadualBuscarMunicipiosPopulacaoUseCase) ordenarMunicipiosPorPopulacao(result []types.MunicipioComDados) {
	sort.Slice(result, func(i, j int) bool {
		return result[i].Populacao < result[j].Populacao
	})
}
