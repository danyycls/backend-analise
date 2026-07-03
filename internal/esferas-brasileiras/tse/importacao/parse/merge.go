package parse

import (
	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/podp/internal/shared/types"
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

func MergeDados(dst, src *tipos.DadosImportacao) {
	if dst == nil || src == nil {
		return
	}

	mergeDimensoesComRemapeamento(dst, src)

	for k, v := range src.Candidatos {
		if existente, ok := dst.Candidatos[k]; ok {
			if existente.ID != v.ID {
				remapearCandidatoIDEmDados(src, v.ID, existente.ID)
			}
		} else {
			dst.Candidatos[k] = v
		}
	}
	for k, v := range src.Fornecedores {
		if _, ok := dst.Fornecedores[k]; !ok {
			dst.Fornecedores[k] = v
		}
	}
	for k, v := range src.Doadores {
		if _, ok := dst.Doadores[k]; !ok {
			dst.Doadores[k] = v
		}
	}
	for k, v := range src.Prestacoes {
		if _, ok := dst.Prestacoes[k]; !ok {
			dst.Prestacoes[k] = v
		}
	}
	for k, v := range src.PrestacoesPorTipoESQ {
		if _, ok := dst.PrestacoesPorTipoESQ[k]; !ok {
			dst.PrestacoesPorTipoESQ[k] = v
		}
	}
	for k, v := range src.ReceitasCandidatoPorSQ {
		if _, ok := dst.ReceitasCandidatoPorSQ[k]; !ok {
			dst.ReceitasCandidatoPorSQ[k] = v
		}
	}
	for k, v := range src.ReceitasOrgaoPorSQ {
		if _, ok := dst.ReceitasOrgaoPorSQ[k]; !ok {
			dst.ReceitasOrgaoPorSQ[k] = v
		}
	}

	if len(src.DespesasCandidato) > 0 {
		dst.DespesasCandidato = append(dst.DespesasCandidato, src.DespesasCandidato...)
	}
	if len(src.DespesasOrgaoPartidario) > 0 {
		dst.DespesasOrgaoPartidario = append(dst.DespesasOrgaoPartidario, src.DespesasOrgaoPartidario...)
	}
	if len(src.ReceitasCandidato) > 0 {
		dst.ReceitasCandidato = append(dst.ReceitasCandidato, src.ReceitasCandidato...)
	}
	if len(src.ReceitasOrgaoPartidario) > 0 {
		dst.ReceitasOrgaoPartidario = append(dst.ReceitasOrgaoPartidario, src.ReceitasOrgaoPartidario...)
	}
	if len(src.ReceitasDoadorOriginarioCandidato) > 0 {
		dst.ReceitasDoadorOriginarioCandidato = append(dst.ReceitasDoadorOriginarioCandidato, src.ReceitasDoadorOriginarioCandidato...)
	}
	if len(src.ReceitasDoadorOriginarioOrgaoPartidario) > 0 {
		dst.ReceitasDoadorOriginarioOrgaoPartidario = append(dst.ReceitasDoadorOriginarioOrgaoPartidario, src.ReceitasDoadorOriginarioOrgaoPartidario...)
	}
	if len(src.Convenios) > 0 {
		dst.Convenios = append(dst.Convenios, src.Convenios...)
	}
	if len(src.BensCandidato) > 0 {
		dst.BensCandidato = append(dst.BensCandidato, src.BensCandidato...)
	}

	src.Convenios = nil
	src.Eleicoes = nil
	src.UnidadesEleitorais = nil
	src.Partidos = nil
	src.Candidatos = nil
	src.Fornecedores = nil
	src.Doadores = nil
	src.Prestacoes = nil
	src.PrestacoesPorTipoESQ = nil
	src.ReceitasCandidatoPorSQ = nil
	src.ReceitasOrgaoPorSQ = nil
	src.DespesasCandidato = nil
	src.DespesasOrgaoPartidario = nil
	src.ReceitasCandidato = nil
	src.ReceitasOrgaoPartidario = nil
	src.ReceitasDoadorOriginarioCandidato = nil
	src.ReceitasDoadorOriginarioOrgaoPartidario = nil
	src.BensCandidato = nil
}

func MergeDadosTransacionais(dst, src *tipos.DadosImportacao) {
	if dst == nil || src == nil {
		return
	}

	mergeDimensoesComRemapeamento(dst, src)

	for k, v := range src.Fornecedores {
		if _, ok := dst.Fornecedores[k]; !ok {
			dst.Fornecedores[k] = v
		}
	}
	for k, v := range src.Doadores {
		if _, ok := dst.Doadores[k]; !ok {
			dst.Doadores[k] = v
		}
	}
	for k, v := range src.Prestacoes {
		if _, ok := dst.Prestacoes[k]; !ok {
			dst.Prestacoes[k] = v
		}
	}
	for k, v := range src.PrestacoesPorTipoESQ {
		if _, ok := dst.PrestacoesPorTipoESQ[k]; !ok {
			dst.PrestacoesPorTipoESQ[k] = v
		}
	}
	for k, v := range src.ReceitasCandidatoPorSQ {
		if _, ok := dst.ReceitasCandidatoPorSQ[k]; !ok {
			dst.ReceitasCandidatoPorSQ[k] = v
		}
	}
	for k, v := range src.ReceitasOrgaoPorSQ {
		if _, ok := dst.ReceitasOrgaoPorSQ[k]; !ok {
			dst.ReceitasOrgaoPorSQ[k] = v
		}
	}

	if len(src.DespesasCandidato) > 0 {
		dst.DespesasCandidato = append(dst.DespesasCandidato, src.DespesasCandidato...)
	}
	if len(src.DespesasOrgaoPartidario) > 0 {
		dst.DespesasOrgaoPartidario = append(dst.DespesasOrgaoPartidario, src.DespesasOrgaoPartidario...)
	}
	if len(src.ReceitasCandidato) > 0 {
		dst.ReceitasCandidato = append(dst.ReceitasCandidato, src.ReceitasCandidato...)
	}
	if len(src.ReceitasOrgaoPartidario) > 0 {
		dst.ReceitasOrgaoPartidario = append(dst.ReceitasOrgaoPartidario, src.ReceitasOrgaoPartidario...)
	}
	if len(src.ReceitasDoadorOriginarioCandidato) > 0 {
		dst.ReceitasDoadorOriginarioCandidato = append(dst.ReceitasDoadorOriginarioCandidato, src.ReceitasDoadorOriginarioCandidato...)
	}
	if len(src.ReceitasDoadorOriginarioOrgaoPartidario) > 0 {
		dst.ReceitasDoadorOriginarioOrgaoPartidario = append(dst.ReceitasDoadorOriginarioOrgaoPartidario, src.ReceitasDoadorOriginarioOrgaoPartidario...)
	}
	if len(src.Convenios) > 0 {
		dst.Convenios = append(dst.Convenios, src.Convenios...)
	}
	if len(src.BensCandidato) > 0 {
		dst.BensCandidato = append(dst.BensCandidato, src.BensCandidato...)
	}

	src.Convenios = nil
	src.Fornecedores = nil
	src.Doadores = nil
	src.Prestacoes = nil
	src.PrestacoesPorTipoESQ = nil
	src.ReceitasCandidatoPorSQ = nil
	src.ReceitasOrgaoPorSQ = nil
	src.DespesasCandidato = nil
	src.DespesasOrgaoPartidario = nil
	src.ReceitasCandidato = nil
	src.ReceitasOrgaoPartidario = nil
	src.ReceitasDoadorOriginarioCandidato = nil
	src.ReceitasDoadorOriginarioOrgaoPartidario = nil
	src.BensCandidato = nil
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
