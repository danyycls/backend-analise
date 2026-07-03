package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	repositorio "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/shared/types"
)

func MontarReceitaCandidatoDetalhada(ctx context.Context, repo *repositorio.Repositorio, r *types.ReceitaCandidato) tsetypes.ReceitaCandidatoDetalhada {
	res := tsetypes.ReceitaCandidatoDetalhada{Receita: r}

	if r.CandidatoID != uuid.Nil {
		cand, err := repo.CandidatoBuscarPorID(ctx, r.CandidatoID)
		if err == nil && cand != nil {
			res.SQCandidato = cand.SQCandidato
			res.NumeroCandidato = cand.NumeroCandidato
			res.NomeCandidato = cand.NomeCompleto
			res.NomeUrnaCandidato = cand.NomeUrna
			res.CargoCandidato = cand.CargoNome
			res.UFCandidato = cand.UFSigla
			res.Candidato = cand

			if cand.PartidoID != nil && *cand.PartidoID != uuid.Nil {
				part, err := repo.PartidosBuscarPorID(ctx, *cand.PartidoID)
				if err == nil && part != nil {
					res.PartidoSigla = part.Sigla
					res.PartidoNome = part.Nome
				}
			}
		}
	}

	res.DescricaoDeVinculo = montarDescricaoReceitaCandidato(r, res)
	return res
}

func MontarReceitaOrgaoPartidarioDetalhada(ctx context.Context, repo *repositorio.Repositorio, r *types.ReceitaOrgaoPartidario) tsetypes.ReceitaOrgaoPartidarioDetalhada {
	res := tsetypes.ReceitaOrgaoPartidarioDetalhada{Receita: r}

	if r.PartidoID != uuid.Nil {
		part, err := repo.PartidosBuscarPorID(ctx, r.PartidoID)
		if err == nil && part != nil {
			res.PartidoNumero = part.Numero
			res.PartidoNome = part.Nome
			res.Partido = part
		}
	}

	res.DescricaoDeVinculo = montarDescricaoReceitaOrgaoPartidario(r, res.PartidoNome)
	return res
}

func montarDescricaoReceitaCandidato(r *types.ReceitaCandidato, det tsetypes.ReceitaCandidatoDetalhada) string {
	dest := "candidato"
	if det.NomeCandidato != "" {
		dest = det.NomeCandidato
	}
	if det.NomeUrnaCandidato != "" && det.NomeUrnaCandidato != det.NomeCandidato {
		dest = det.NomeUrnaCandidato + " (" + det.NomeCandidato + ")"
	}
	if det.PartidoSigla != "" {
		dest += " / " + det.PartidoSigla
	}
	if det.CargoCandidato != "" {
		dest += " para " + det.CargoCandidato
	}
	if det.UFCandidato != "" {
		dest += " - " + det.UFCandidato
	}

	valor := r.Valor
	data := ""
	if r.DataReceita != nil {
		data = fmt.Sprintf(" em %s", r.DataReceita.Format("02/01/2006"))
	}

	return fmt.Sprintf("Doação de R$ %.2f a %s%s", valor, dest, data)
}

func montarDescricaoReceitaOrgaoPartidario(r *types.ReceitaOrgaoPartidario, partidoNome string) string {
	dest := "partido"
	if partidoNome != "" {
		dest = partidoNome
	}

	valor := r.Valor
	data := ""
	if r.DataReceita != nil {
		data = fmt.Sprintf(" em %s", r.DataReceita.Format("02/01/2006"))
	}

	return fmt.Sprintf("Doação de R$ %.2f ao %s%s", valor, dest, data)
}
