package usecase

import (
	"context"

	repositorio "github.com/danyele/laceu/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/laceu/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/laceu/internal/shared/database"
	"github.com/danyele/laceu/internal/shared/logger"

	"github.com/google/uuid"
)

type BuscarCandidatosUseCase struct {
	db database.DB
}

func NovoBuscarCandidatosUseCase(db database.DB) *BuscarCandidatosUseCase {
	return &BuscarCandidatosUseCase{db: db}
}

func (u *BuscarCandidatosUseCase) ExecutarListarCargos(ctx context.Context) (*tsetypes.OpcoesFiltroResponse, error) {
	repo := repositorio.Novo(u.db)
	cargos, err := repo.CargosDistintos(ctx)
	if err != nil {
		return nil, err
	}
	var opcoes []tsetypes.OpcaoFiltro
	seen := make(map[string]bool)
	for _, c := range cargos {
		if c.CargoNome == "" || seen[c.CargoNome] {
			continue
		}
		seen[c.CargoNome] = true
		opcoes = append(opcoes, tsetypes.OpcaoFiltro{
			Valor: c.CargoNome,
			Label: c.CargoNome,
		})
	}
	return &tsetypes.OpcoesFiltroResponse{Opcoes: opcoes}, nil
}

func (u *BuscarCandidatosUseCase) ExecutarListarPartidos(ctx context.Context) (*tsetypes.OpcoesFiltroResponse, error) {
	repo := repositorio.Novo(u.db)
	partidos, err := repo.PartidosListarDistintos(ctx)
	if err != nil {
		return nil, err
	}
	var opcoes []tsetypes.OpcaoFiltro
	for _, p := range partidos {
		opcoes = append(opcoes, tsetypes.OpcaoFiltro{
			Valor: p.ID.String(),
			Label: p.Sigla + " - " + p.Nome,
		})
	}
	return &tsetypes.OpcoesFiltroResponse{Opcoes: opcoes}, nil
}

func (u *BuscarCandidatosUseCase) Executar(ctx context.Context, req *tsetypes.BuscaCandidatosRequest) (*tsetypes.CandidatosResponse, error) {
	log := logger.New("TSE: UseCase: BuscarCandidatos")
	repo := repositorio.Novo(u.db)
	situacao := req.Eleito
	log.Info("requisicao recebida", "cargo_nome", req.CargoNome, "partido_id", req.PartidoID, "eleito", req.Eleito, "uf_sigla", req.UFSigla)
	candidatos, err := repo.CandidatoBuscarPorFiltros(ctx, req.CargoNome, req.PartidoID, req.UFSigla, situacao)
	if err != nil {
		log.Error("erro ao buscar candidatos", "erro", err)
		return nil, err
	}
	var dtos []tsetypes.CandidatoLista
	for _, c := range candidatos {
		dto := tsetypes.CandidatoLista{
			SQCandidato:                  c.SQCandidato,
			NomeCompleto:                 c.NomeCompleto,
			NomeUrna:                     c.NomeUrna,
			CPF:                          c.CPF,
			NumeroCandidato:              c.NumeroCandidato,
			CargoCodigo:                  c.CargoCodigo,
			CargoNome:                    c.CargoNome,
			UFSigla:                      c.UFSigla,
			SituacaoTotalizacaoDescricao: c.SituacaoTotalizacaoDescricao,
			Eleito:                       isEleito(c.SituacaoTotalizacaoDescricao),
		}
		if c.PartidoID != nil && *c.PartidoID != uuid.Nil {
			part, err := repo.PartidosBuscarPorID(ctx, *c.PartidoID)
			if err == nil && part != nil {
				dto.Partido = &tsetypes.PartidoResumido{
					Numero: part.Numero,
					Sigla:  part.Sigla,
					Nome:   part.Nome,
				}
			}
		}
		dtos = append(dtos, dto)
	}
	log.Info("resultados encontrados", "total", len(dtos))
	return &tsetypes.CandidatosResponse{
		Candidatos: dtos,
		Total:      len(dtos),
	}, nil
}

func isEleito(situacao string) bool {
	switch situacao {
	case "ELEITO", "ELEITO POR QP", "ELEITO POR MÉDIA":
		return true
	}
	return false
}

func init() {
	log := logger.New("TSE: UseCase: BuscarCandidatosUseCase")
	log.Info("BuscarCandidatosUseCase carregado")
}
