package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	repositorio "github.com/danyele/laceu/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/laceu/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/laceu/internal/shared/types"
)

func MontarDespesaCandidatoDetalhada(ctx context.Context, repo *repositorio.Repositorio, d *types.DespesaCandidato) tsetypes.DespesaCandidatoDetalhada {
	res := tsetypes.DespesaCandidatoDetalhada{Despesa: d}
	if d.CandidatoID != uuid.Nil {
		cand, err := repo.CandidatoBuscarPorID(ctx, d.CandidatoID)
		if err == nil && cand != nil {
			res.SQCandidato = cand.SQCandidato
		}
	}
	if d.PrestacaoContasID != uuid.Nil {
		p, err := repo.PrestacaoBuscarPorID(ctx, d.PrestacaoContasID)
		if err == nil && p != nil {
			res.SQPrestacao = p.SQPrestadorContas
		}
	}
	return res
}

func MontarDespesaOrgaoPartidarioDetalhada(ctx context.Context, repo *repositorio.Repositorio, d *types.DespesaOrgaoPartidario) tsetypes.DespesaOrgaoPartidarioDetalhada {
	res := tsetypes.DespesaOrgaoPartidarioDetalhada{Despesa: d}
	if d.PrestacaoContasID != uuid.Nil {
		p, err := repo.PrestacaoBuscarPorID(ctx, d.PrestacaoContasID)
		if err == nil && p != nil {
			res.SQPrestacao = p.SQPrestadorContas
		}
	}
	if d.PartidoID != uuid.Nil {
		part, err := repo.PartidosBuscarPorID(ctx, d.PartidoID)
		if err == nil && part != nil {
			res.PartidoNumero = part.Numero
			res.PartidoNome = part.Nome
		}
	}
	res.DescricaoDeVinculo = montarDescricaoDespesaOrgaoPartidario(d, res.PartidoNome)
	return res
}

func montarDescricaoDespesaOrgaoPartidario(d *types.DespesaOrgaoPartidario, partidoNome string) string {
	nomePart := "partido"
	if partidoNome != "" {
		nomePart = partidoNome
	}

	valor := d.Valor
	data := ""
	if d.DataDespesa != nil {
		data = fmt.Sprintf(" em %s", d.DataDespesa.Format("02/01/2006"))
	}

	return fmt.Sprintf("Despesa partidária de R$ %.2f para %s do partido %s%s", valor, d.Descricao, nomePart, data)
}
