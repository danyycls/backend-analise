package utils

import (
	"context"

	repositorio "github.com/danyele/laceu/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/laceu/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/laceu/internal/shared/types"
)

func MontarFornecedorDetalhado(ctx context.Context, repo *repositorio.Repositorio, f *types.Fornecedor) *tsetypes.FornecedorDetalhado {
	dto := &tsetypes.FornecedorDetalhado{
		Fornecedor: FornecedorParaEnriquecido(f),
	}
	desps, _ := repo.DespesasCandidatoBuscarPorFornecedorID(ctx, f.ID)
	for _, d := range desps {
		dto.DespesasCandidato = append(dto.DespesasCandidato,
			MontarDespesaCandidatoDetalhada(ctx, repo, d))
	}
	despsPart, _ := repo.DespesasPartidoBuscarPorFornecedorID(ctx, f.ID)
	for _, d := range despsPart {
		dto.DespesasOrgaoPartidario = append(dto.DespesasOrgaoPartidario,
			MontarDespesaOrgaoPartidarioDetalhada(ctx, repo, d))
	}
	return dto
}

func FornecedorParaEnriquecido(f *types.Fornecedor) tsetypes.FornecedorEnriquecido {
	return tsetypes.FornecedorEnriquecido{
		CPFCNPJ:                    f.CPFCNPJ,
		Nome:                       f.Nome,
		NomeRFB:                    f.NomeRFB,
		TipoFornecedorCodigo:       f.TipoFornecedorCodigo,
		TipoFornecedorDescricao:    f.TipoFornecedorDescricao,
		CNAECodigo:                 f.CNAECodigo,
		CNAEDescricao:              f.CNAEDescricao,
		EsferaPartidariaCodigo:     f.EsferaPartidariaCodigo,
		EsferaPartidariaDescricao:  f.EsferaPartidariaDescricao,
		UFSigla:                    f.UFSigla,
		MunicipioNome:              f.MunicipioNome,
		SQCandidatoRelacionado:     f.SQCandidatoRelacionado,
		NumeroCandidatoRelacionado: f.NumeroCandidatoRelacionado,
		CargoCodigoRelacionado:     f.CargoCodigoRelacionado,
		CargoDescricaoRelacionada:  f.CargoDescricaoRelacionada,
		PartidoNumeroRelacionado:   f.PartidoNumeroRelacionado,
		PartidoSiglaRelacionado:    f.PartidoSiglaRelacionado,
		PartidoNomeRelacionado:     f.PartidoNomeRelacionado,
	}
}
