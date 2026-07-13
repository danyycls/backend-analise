package parse

import (
	"github.com/danyele/podp/internal/shared/types"
	tipos "github.com/danyele/podp/internal/sources/tse/importacao/types"
	"github.com/google/uuid"
)

func ponteiroUUID(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	valor := id
	return &valor
}

func remapearEleicaoIDEmDados(d *tipos.DadosImportacao, antigo, novo uuid.UUID) {
	if d == nil || antigo == uuid.Nil || antigo == novo {
		return
	}
	for _, candidato := range d.Candidatos {
		if candidato.EleicaoID == antigo {
			candidato.EleicaoID = novo
		}
	}
	prestacoesAtualizadas := make(map[string]*types.PrestacaoContas, len(d.Prestacoes))
	for _, prestacao := range d.Prestacoes {
		if prestacao.EleicaoID == antigo {
			prestacao.EleicaoID = novo
		}
		chave := chavePrestacaoNatural(prestacao.TipoPrestador, prestacao.EleicaoID, prestacao.SQPrestadorContas)
		prestacoesAtualizadas[chave] = prestacao
	}
	d.Prestacoes = prestacoesAtualizadas
}

func remapearUnidadeEleitoralIDEmDados(d *tipos.DadosImportacao, antigo, novo uuid.UUID) {
	if d == nil || antigo == uuid.Nil || antigo == novo {
		return
	}
	for _, prestacao := range d.Prestacoes {
		if prestacao.UnidadeEleitoralID != nil && *prestacao.UnidadeEleitoralID == antigo {
			prestacao.UnidadeEleitoralID = ponteiroUUID(novo)
		}
	}
}

func remapearPartidoIDEmDados(d *tipos.DadosImportacao, antigo, novo uuid.UUID) {
	if d == nil || antigo == uuid.Nil || antigo == novo {
		return
	}
	for _, candidato := range d.Candidatos {
		if candidato.PartidoID != nil && *candidato.PartidoID == antigo {
			candidato.PartidoID = ponteiroUUID(novo)
		}
	}
	for _, prestacao := range d.Prestacoes {
		if prestacao.PartidoID != nil && *prestacao.PartidoID == antigo {
			prestacao.PartidoID = ponteiroUUID(novo)
		}
	}
	for _, despesa := range d.DespesasOrgaoPartidario {
		if despesa.PartidoID == antigo {
			despesa.PartidoID = novo
		}
	}
	for _, receita := range d.ReceitasOrgaoPartidario {
		if receita.PartidoID == antigo {
			receita.PartidoID = novo
		}
	}
}

func remapearCandidatoIDEmDados(d *tipos.DadosImportacao, antigo, novo uuid.UUID) {
	if d == nil || antigo == uuid.Nil || antigo == novo {
		return
	}
	for _, bem := range d.BensCandidato {
		if bem.CandidatoID == antigo {
			bem.CandidatoID = novo
		}
	}
	for _, prestacao := range d.Prestacoes {
		if prestacao.CandidatoID != nil && *prestacao.CandidatoID == antigo {
			prestacao.CandidatoID = ponteiroUUID(novo)
		}
	}
	for _, despesa := range d.DespesasCandidato {
		if despesa.CandidatoID == antigo {
			despesa.CandidatoID = novo
		}
	}
	for _, receita := range d.ReceitasCandidato {
		if receita.CandidatoID == antigo {
			receita.CandidatoID = novo
		}
	}
}

func appendSlice[T any](dst []T, src []T) []T {
	if len(src) == 0 {
		return dst
	}
	novaCap := len(dst) + len(src)
	if cap(dst) < novaCap {
		nova := make([]T, len(dst), novaCap)
		copy(nova, dst)
		dst = nova
	}
	return append(dst, src...)
}

func mergeMap[K comparable, V any](dst, src map[K]V) {
	if dst == nil || len(src) == 0 {
		return
	}
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			dst[k] = v
		}
	}
}

func limparMap[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

func limparSrc(dados *tipos.DadosImportacao, incluirDimensoes bool) {
	limparMap(dados.Fornecedores)
	limparMap(dados.Doadores)
	limparMap(dados.Prestacoes)
	limparMap(dados.PrestacoesPorTipoESQ)
	limparMap(dados.PrestacoesPorID)
	limparMap(dados.ReceitasCandidatoPorSQ)
	limparMap(dados.ReceitasOrgaoPorSQ)
	limparMap(dados.BensCandidato)
	dados.Convenios = nil
	dados.DespesasCandidato = nil
	dados.DespesasOrgaoPartidario = nil
	dados.ReceitasCandidato = nil
	dados.ReceitasOrgaoPartidario = nil
	dados.ReceitasDoadorOriginarioCandidato = nil
	dados.ReceitasDoadorOriginarioOrgaoPartidario = nil
	if incluirDimensoes {
		limparMap(dados.Eleicoes)
		limparMap(dados.UnidadesEleitorais)
		limparMap(dados.Partidos)
		limparMap(dados.Candidatos)
		limparMap(dados.CandidatosPorID)
	}
}

func mergeTransacional(dst, src *tipos.DadosImportacao) {
	mergeMap(dst.Fornecedores, src.Fornecedores)
	mergeMap(dst.Doadores, src.Doadores)
	mergeMap(dst.Prestacoes, src.Prestacoes)
	mergeMap(dst.PrestacoesPorTipoESQ, src.PrestacoesPorTipoESQ)
	mergeMap(dst.PrestacoesPorID, src.PrestacoesPorID)
	mergeMap(dst.ReceitasCandidatoPorSQ, src.ReceitasCandidatoPorSQ)
	mergeMap(dst.ReceitasOrgaoPorSQ, src.ReceitasOrgaoPorSQ)

	dst.DespesasCandidato = appendSlice(dst.DespesasCandidato, src.DespesasCandidato)
	dst.DespesasOrgaoPartidario = appendSlice(dst.DespesasOrgaoPartidario, src.DespesasOrgaoPartidario)
	dst.ReceitasCandidato = appendSlice(dst.ReceitasCandidato, src.ReceitasCandidato)
	dst.ReceitasOrgaoPartidario = appendSlice(dst.ReceitasOrgaoPartidario, src.ReceitasOrgaoPartidario)
	dst.ReceitasDoadorOriginarioCandidato = appendSlice(dst.ReceitasDoadorOriginarioCandidato, src.ReceitasDoadorOriginarioCandidato)
	dst.ReceitasDoadorOriginarioOrgaoPartidario = appendSlice(dst.ReceitasDoadorOriginarioOrgaoPartidario, src.ReceitasDoadorOriginarioOrgaoPartidario)
	dst.Convenios = appendSlice(dst.Convenios, src.Convenios)
	mergeMap(dst.BensCandidato, src.BensCandidato)
}

func MergeDados(dst, src *tipos.DadosImportacao) {
	if dst == nil || src == nil {
		return
	}
	mergeDimensoesComRemapeamento(dst, src)

	for k, v := range src.Candidatos {
		if existente, ok := dst.Candidatos[k]; ok {
			if existente.ID != v.ID {
				remapearCandidatoIDEmDados(src, v.ID, existente.ID)
				delete(src.CandidatosPorID, v.ID)
				src.CandidatosPorID[existente.ID] = v
			}
		} else {
			dst.Candidatos[k] = v
			dst.CandidatosPorID[v.ID] = v
		}
	}

	mergeTransacional(dst, src)
	limparSrc(src, true)
}

func MergeDadosTransacionais(dst, src *tipos.DadosImportacao) {
	if dst == nil || src == nil {
		return
	}
	mergeDimensoesComRemapeamento(dst, src)
	mergeTransacional(dst, src)
	limparSrc(src, false)
}

func mergeDimensoesComRemapeamento(dst, src *tipos.DadosImportacao) {
	for k, v := range src.Eleicoes {
		if existente, ok := dst.Eleicoes[k]; ok {
			if existente.ID != v.ID {
				remapearEleicaoIDEmDados(src, v.ID, existente.ID)
			}
		} else {
			dst.Eleicoes[k] = v
		}
	}
	for k, v := range src.UnidadesEleitorais {
		if existente, ok := dst.UnidadesEleitorais[k]; ok {
			if existente.ID != v.ID {
				remapearUnidadeEleitoralIDEmDados(src, v.ID, existente.ID)
			}
		} else {
			dst.UnidadesEleitorais[k] = v
		}
	}
	for k, v := range src.Partidos {
		if existente, ok := dst.Partidos[k]; ok {
			if existente.ID != v.ID {
				remapearPartidoIDEmDados(src, v.ID, existente.ID)
			}
		} else {
			dst.Partidos[k] = v
		}
	}
}
