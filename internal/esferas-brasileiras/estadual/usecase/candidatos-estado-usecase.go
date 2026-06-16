package usecase

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"

	estadual "github.com/danyele/laceu/internal/esferas-brasileiras/estadual"
	repositorio "github.com/danyele/laceu/internal/esferas-brasileiras/tse/repositorio"
	"github.com/danyele/laceu/internal/shared/database"
	"github.com/danyele/laceu/internal/shared/types"
)

type EsferaEstadualBuscarCandidatosRequest struct {
	UF string
}

type EsferaEstadualBuscarCandidatosResponse struct {
	Dados *types.DadosCandidatosEstado
}

type EsferaEstadualBuscarCandidatosUseCase struct {
	db database.DB
}

func NovoEsferaEstadualBuscarCandidatosUseCase(db database.DB) *EsferaEstadualBuscarCandidatosUseCase {
	return &EsferaEstadualBuscarCandidatosUseCase{
		db: db,
	}
}

func (u *EsferaEstadualBuscarCandidatosUseCase) Executar(ctx context.Context, req *EsferaEstadualBuscarCandidatosRequest) (*EsferaEstadualBuscarCandidatosResponse, error) {
	repo := repositorio.Novo(u.db)
	result := &types.DadosCandidatosEstado{}

	candidatos, err := u.buscarCandidatosPorUF(ctx, repo, req.UF)
	if err != nil {
		return &EsferaEstadualBuscarCandidatosResponse{Dados: result}, fmt.Errorf("erro buscar candidatos: %w", err)
	}

	idsPartido, idsEleicao := u.coletarIDsDistintos(candidatos)
	partidoMap, eleicaoMap := u.buscarPartidosEEleicoes(ctx, repo, idsPartido, idsEleicao)

	u.montarCandidatosEnriquecidos(candidatos, partidoMap, eleicaoMap, result)
	u.ordenarResultadosPorAno(result)

	return &EsferaEstadualBuscarCandidatosResponse{Dados: result}, nil
}

func (u *EsferaEstadualBuscarCandidatosUseCase) buscarCandidatosPorUF(
	ctx context.Context,
	repo *repositorio.Repositorio,
	uf string,
) ([]*types.Candidato, error) {
	return repo.CandidatosEleitosPorUF(ctx, uf, []string{"PREFEITO", "VICE-PREFEITO", "VEREADOR"})
}

func (u *EsferaEstadualBuscarCandidatosUseCase) coletarIDsDistintos(
	candidatos []*types.Candidato,
) ([]uuid.UUID, []uuid.UUID) {
	idsPartido := make([]uuid.UUID, 0)
	idsEleicao := make([]uuid.UUID, 0)
	seenPartido := make(map[uuid.UUID]bool)
	seenEleicao := make(map[uuid.UUID]bool)

	for _, c := range candidatos {
		if c.PartidoID != nil && !seenPartido[*c.PartidoID] {
			idsPartido = append(idsPartido, *c.PartidoID)
			seenPartido[*c.PartidoID] = true
		}
		if c.ID != [16]byte{} && !seenEleicao[c.EleicaoID] {
			idsEleicao = append(idsEleicao, c.EleicaoID)
			seenEleicao[c.EleicaoID] = true
		}
	}

	return idsPartido, idsEleicao
}

func (u *EsferaEstadualBuscarCandidatosUseCase) buscarPartidosEEleicoes(
	ctx context.Context,
	repo *repositorio.Repositorio,
	idsPartido []uuid.UUID,
	idsEleicao []uuid.UUID,
) (map[uuid.UUID]*types.Partido, map[uuid.UUID]*types.Eleicao) {
	partidoMap, _ := repo.PartidosBuscarPorIDs(ctx, idsPartido)
	eleicaoMap, _ := repo.EleicoesBuscarPorIDs(ctx, idsEleicao)
	return partidoMap, eleicaoMap
}

func (u *EsferaEstadualBuscarCandidatosUseCase) montarCandidatosEnriquecidos(
	candidatos []*types.Candidato,
	partidoMap map[uuid.UUID]*types.Partido,
	eleicaoMap map[uuid.UUID]*types.Eleicao,
	result *types.DadosCandidatosEstado,
) {
	for _, c := range candidatos {
		dto := u.montarCandidatoDTO(c, partidoMap, eleicaoMap)
		u.classificarCandidatoPorCargo(dto, result)
	}
}

func (u *EsferaEstadualBuscarCandidatosUseCase) montarCandidatoDTO(
	c *types.Candidato,
	partidoMap map[uuid.UUID]*types.Partido,
	eleicaoMap map[uuid.UUID]*types.Eleicao,
) types.CandidatoEleito {
	partidoSigla := ""
	partidoNome := ""
	if c.PartidoID != nil && partidoMap != nil {
		if p, ok := partidoMap[*c.PartidoID]; ok {
			partidoSigla = p.Sigla
			partidoNome = p.Nome
		}
	}

	ano := int16(0)
	eleicaoDesc := ""
	eleicaoData := ""
	eleicaoTipo := ""
	if eleicaoMap != nil {
		if e, ok := eleicaoMap[c.EleicaoID]; ok {
			ano = e.Ano
			eleicaoDesc = e.Descricao
			if e.DataEleicao != nil {
				eleicaoData = e.DataEleicao.Format("02/01/2006")
			}
			eleicaoTipo = estadual.NormalizarTipoEleicao(e.NomeTipoEleicao)
		}
	}

	return types.CandidatoEleito{
		ID:                           c.ID.String(),
		SQCandidato:                  c.SQCandidato,
		NomeUrna:                     c.NomeUrna,
		NomeCompleto:                 c.NomeCompleto,
		PartidoSigla:                 partidoSigla,
		PartidoNome:                  partidoNome,
		CargoNome:                    c.CargoNome,
		SituacaoTotalizacaoDescricao: c.SituacaoTotalizacaoDescricao,
		AnoEleicao:                   ano,
		NumeroCandidato:              c.NumeroCandidato,
		CPF:                          c.CPF,
		EleicaoDescricao:             eleicaoDesc,
		EleicaoData:                  eleicaoData,
		EleicaoTipo:                  eleicaoTipo,
	}
}

func (u *EsferaEstadualBuscarCandidatosUseCase) classificarCandidatoPorCargo(
	dto types.CandidatoEleito,
	result *types.DadosCandidatosEstado,
) {
	switch dto.CargoNome {
	case "PREFEITO":
		result.Prefeitos = append(result.Prefeitos, dto)
	case "VICE-PREFEITO":
		result.VicePrefeitos = append(result.VicePrefeitos, dto)
	case "VEREADOR":
		result.Vereadores = append(result.Vereadores, dto)
	}
}

func (u *EsferaEstadualBuscarCandidatosUseCase) ordenarResultadosPorAno(result *types.DadosCandidatosEstado) {
	sort.Slice(result.Prefeitos, func(i, j int) bool {
		return result.Prefeitos[i].AnoEleicao > result.Prefeitos[j].AnoEleicao
	})
	sort.Slice(result.VicePrefeitos, func(i, j int) bool {
		return result.VicePrefeitos[i].AnoEleicao > result.VicePrefeitos[j].AnoEleicao
	})
	sort.Slice(result.Vereadores, func(i, j int) bool {
		return result.Vereadores[i].AnoEleicao > result.Vereadores[j].AnoEleicao
	})
}
