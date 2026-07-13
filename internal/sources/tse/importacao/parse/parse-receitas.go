package parse

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
	tipos "github.com/danyele/podp/internal/sources/tse/importacao/types"
)

func (p *ProcessadorLeitorCSV) processarReceitaCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoCandidatoContratada(ctx, registro)
		if err != nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado", "caminho", caminho, "linha", numeroLinha,
				"sq_receita", registro["SQ_RECEITA"], "sq_candidato", registro["SQ_CANDIDATO"], "erro", err)
			return nil
		}
		candidatoID := uuidOpcional(prestacao.CandidatoID, "prestacao_contas.candidato_id")
		if candidatoID == nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado: candidato_id vazio",
				"caminho", caminho, "linha", numeroLinha,
				"sq_receita", registro["SQ_RECEITA"], "sq_candidato", registro["SQ_CANDIDATO"])
			return nil
		}

		doadorID, err := p.garantirDoadorReceitaOpcional(ctx, registro)
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqReceita := inteiro64OuZero(registro["SQ_RECEITA"])
		receita := &types.ReceitaCandidato{
			PrestacaoContasID:        prestacao.ID,
			SQPrestadorContas:        prestacao.SQPrestadorContas,
			CandidatoID:              *candidatoID,
			DoadorID:                 doadorID,
			SQReceita:                sqReceita,
			FonteReceitaCodigo:       inteiroOpcional(registro["CD_FONTE_RECEITA"]),
			FonteReceitaDescricao:    textoOpcional(registro["DS_FONTE_RECEITA"]),
			OrigemReceitaCodigo:      inteiroOpcional(registro["CD_ORIGEM_RECEITA"]),
			OrigemReceitaDescricao:   textoOpcional(registro["DS_ORIGEM_RECEITA"]),
			NaturezaReceitaCodigo:    inteiroOpcional(registro["CD_NATUREZA_RECEITA"]),
			NaturezaReceitaDescricao: textoOpcional(registro["DS_NATUREZA_RECEITA"]),
			EspecieReceitaCodigo:     inteiroOpcional(registro["CD_ESPECIE_RECEITA"]),
			EspecieReceitaDescricao:  textoOpcional(registro["DS_ESPECIE_RECEITA"]),
			NumeroReciboDoacao:       textoOpcional(registro["NR_RECIBO_DOACAO"]),
			NumeroDocumentoDoacao:    textoOpcional(registro["NR_DOCUMENTO_DOACAO"]),
			DataReceita:              dataOpcional(registro["DT_RECEITA"]),
			Descricao:                textoComFallback(registro["DS_RECEITA"], "RECEITA SEM DESCRICAO"),
			Valor:                    decimalOuZero(registro["VR_RECEITA"]),
			NaturezaRecursoEstimavel: textoOpcional(registro["DS_NATUREZA_RECURSO_ESTIMAVEL"]),
			Genero:                   textoOpcional(registro["DS_GENERO"]),
			CorRaca:                  textoOpcional(registro["DS_COR_RACA"]),
		}
		receita.ID = uuid.Must(uuid.NewV7())
		p.dados.ReceitasCandidato = append(p.dados.ReceitasCandidato, receita)
		p.dados.ReceitasCandidatoPorSQ[sqReceita] = receita
		return nil
	})
}

func (p *ProcessadorLeitorCSV) processarReceitaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoOrgaoContratada(ctx, registro)
		if err != nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado", "caminho", caminho, "linha", numeroLinha,
				"sq_receita", registro["SQ_RECEITA"], "sq_prestador_contas", registro["SQ_PRESTADOR_CONTAS"], "erro", err)
			return nil
		}
		partidoID := uuidOpcional(prestacao.PartidoID, "prestacao_contas.partido_id")
		if partidoID == nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado: partido_id vazio",
				"caminho", caminho, "linha", numeroLinha,
				"sq_receita", registro["SQ_RECEITA"], "sq_prestador_contas", registro["SQ_PRESTADOR_CONTAS"])
			return nil
		}

		doadorID, err := p.garantirDoadorReceitaOpcional(ctx, registro)
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqReceita := inteiro64OuZero(registro["SQ_RECEITA"])
		receita := &types.ReceitaOrgaoPartidario{
			PrestacaoContasID:        prestacao.ID,
			SQPrestadorContas:        prestacao.SQPrestadorContas,
			PartidoID:                *partidoID,
			DoadorID:                 doadorID,
			SQReceita:                sqReceita,
			FonteReceitaCodigo:       inteiroOpcional(registro["CD_FONTE_RECEITA"]),
			FonteReceitaDescricao:    textoOpcional(registro["DS_FONTE_RECEITA"]),
			OrigemReceitaCodigo:      inteiroOpcional(registro["CD_ORIGEM_RECEITA"]),
			OrigemReceitaDescricao:   textoOpcional(registro["DS_ORIGEM_RECEITA"]),
			NaturezaReceitaCodigo:    inteiroOpcional(registro["CD_NATUREZA_RECEITA"]),
			NaturezaReceitaDescricao: textoOpcional(registro["DS_NATUREZA_RECEITA"]),
			EspecieReceitaCodigo:     inteiroOpcional(registro["CD_ESPECIE_RECEITA"]),
			EspecieReceitaDescricao:  textoOpcional(registro["DS_ESPECIE_RECEITA"]),
			NumeroReciboDoacao:       textoOpcional(registro["NR_RECIBO_DOACAO"]),
			NumeroDocumentoDoacao:    textoOpcional(registro["NR_DOCUMENTO_DOACAO"]),
			DataReceita:              dataOpcional(registro["DT_RECEITA"]),
			Descricao:                textoComFallback(registro["DS_RECEITA"], "RECEITA SEM DESCRICAO"),
			Valor:                    decimalOuZero(registro["VR_RECEITA"]),
		}
		receita.ID = uuid.Must(uuid.NewV7())
		p.dados.ReceitasOrgaoPartidario = append(p.dados.ReceitasOrgaoPartidario, receita)
		p.dados.ReceitasOrgaoPorSQ[sqReceita] = receita
		return nil
	})
}

func (p *ProcessadorLeitorCSV) processarReceitaCandidatoDoadorOriginario(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao := p.garantirOuCriarPrestacaoCandidato(ctx, registro)
		if prestacao == nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado", "caminho", caminho, "linha", numeroLinha,
				"sq_receita", registro["SQ_RECEITA"], "sq_prestador_contas", registro["SQ_PRESTADOR_CONTAS"])
			return nil
		}

		sqReceita := inteiro64OuZero(registro["SQ_RECEITA"])
		receitaID := receitaCandidatoIDPorSQ(p.dados, sqReceita)
		origem := &types.ReceitaDoadorOriginarioCandidato{
			PrestacaoContasID:  prestacao.ID,
			SQPrestadorContas:  prestacao.SQPrestadorContas,
			ReceitaCandidatoID: receitaID,
			SQReceita:          sqReceita,
			DocumentoDoador:    documentoOpcional(registro["NR_CPF_CNPJ_DOADOR_ORIGINARIO"]),
			NomeDoador:         textoComFallback(registro["NM_DOADOR_ORIGINARIO"], "DOADOR ORIGINARIO SEM NOME"),
			NomeDoadorRFB:      textoOpcional(registro["NM_DOADOR_ORIGINARIO_RFB"]),
			TipoDoador:         textoOpcional(registro["TP_DOADOR_ORIGINARIO"]),
			CNAECodigo:         textoOpcional(registro["CD_CNAE_DOADOR_ORIGINARIO"]),
			CNAEDescricao:      textoOpcional(registro["DS_CNAE_DOADOR_ORIGINARIO"]),
			DataReceita:        dataOpcional(registro["DT_RECEITA"]),
			Descricao:          textoComFallback(registro["DS_RECEITA"], "RECEITA SEM DESCRICAO"),
			Valor:              decimalOuZero(registro["VR_RECEITA"]),
		}
		origem.ID = uuid.Must(uuid.NewV7())
		p.dados.ReceitasDoadorOriginarioCandidato = append(p.dados.ReceitasDoadorOriginarioCandidato, origem)
		return nil
	})
}

func (p *ProcessadorLeitorCSV) processarReceitaOrgaoPartidarioDoadorOriginario(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao := p.garantirOuCriarPrestacaoOrgao(ctx, registro)
		if prestacao == nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado", "caminho", caminho, "linha", numeroLinha,
				"sq_receita", registro["SQ_RECEITA"], "sq_prestador_contas", registro["SQ_PRESTADOR_CONTAS"])
			return nil
		}

		sqReceita := inteiro64OuZero(registro["SQ_RECEITA"])
		receitaID := receitaOrgaoIDPorSQ(p.dados, sqReceita)
		origem := &types.ReceitaDoadorOriginarioOrgaoPartidario{
			PrestacaoContasID:        prestacao.ID,
			SQPrestadorContas:        prestacao.SQPrestadorContas,
			ReceitaOrgaoPartidarioID: receitaID,
			SQReceita:                sqReceita,
			DocumentoDoador:          documentoOpcional(registro["NR_CPF_CNPJ_DOADOR_ORIGINARIO"]),
			NomeDoador:               textoComFallback(registro["NM_DOADOR_ORIGINARIO"], "DOADOR ORIGINARIO SEM NOME"),
			NomeDoadorRFB:            textoOpcional(registro["NM_DOADOR_ORIGINARIO_RFB"]),
			TipoDoador:               textoOpcional(registro["TP_DOADOR_ORIGINARIO"]),
			CNAECodigo:               textoOpcional(registro["CD_CNAE_DOADOR_ORIGINARIO"]),
			CNAEDescricao:            textoOpcional(registro["DS_CNAE_DOADOR_ORIGINARIO"]),
			DataReceita:              dataOpcional(registro["DT_RECEITA"]),
			Descricao:                textoComFallback(registro["DS_RECEITA"], "RECEITA SEM DESCRICAO"),
			Valor:                    decimalOuZero(registro["VR_RECEITA"]),
		}
		origem.ID = uuid.Must(uuid.NewV7())
		p.dados.ReceitasDoadorOriginarioOrgaoPartidario = append(p.dados.ReceitasDoadorOriginarioOrgaoPartidario, origem)
		return nil
	})
}

func receitaCandidatoIDPorSQ(dados *tipos.DadosImportacao, sq int64) *uuid.UUID {
	if receita, ok := dados.ReceitasCandidatoPorSQ[sq]; ok {
		return ponteiroUUIDLocal(receita.ID)
	}
	return nil
}

func receitaOrgaoIDPorSQ(dados *tipos.DadosImportacao, sq int64) *uuid.UUID {
	if receita, ok := dados.ReceitasOrgaoPorSQ[sq]; ok {
		return ponteiroUUIDLocal(receita.ID)
	}
	return nil

}

func ponteiroUUIDLocal(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	valor := id
	return &valor
}

func (p *ProcessadorLeitorCSV) resolverEleicaoOpcional(registro map[string]string) uuid.UUID {
	codigo := inteiroOpcional(registro["CD_ELEICAO"])
	if codigo == nil {
		return garantirEleicaoSentinela(p.dados)
	}
	if existente, ok := p.dados.Eleicoes[*codigo]; ok {
		return existente.ID
	}
	ano := inteiroOpcional(registro["AA_ELEICAO"])
	if ano == nil {
		return garantirEleicaoSentinela(p.dados)
	}
	e := &types.Eleicao{
		Ano:       int16(*ano),
		CodigoTSE: *codigo,
		Descricao: textoComFallback(registro["DS_ELEICAO"], "Eleicao sem descricao"),
	}
	e.ID = uuid.Must(uuid.NewV7())
	p.dados.Eleicoes[*codigo] = e
	return e.ID
}

func (p *ProcessadorLeitorCSV) garantirOuCriarPrestacaoCandidato(ctx context.Context, registro map[string]string) *types.PrestacaoContas {
	prestacao, err := p.garantirPrestacaoPorTipoESQ(ctx, tipos.TipoPrestadorCandidato, registro["SQ_PRESTADOR_CONTAS"], registro["CD_ELEICAO"])
	if err == nil {
		return prestacao
	}

	sqPrestador := inteiro64Opcional(registro["SQ_PRESTADOR_CONTAS"])
	if sqPrestador == nil {
		return nil
	}

	eleicaoID := p.resolverEleicaoOpcional(registro)
	chave := chavePrestacaoNatural(tipos.TipoPrestadorCandidato, eleicaoID, *sqPrestador)
	if existente, ok := p.dados.Prestacoes[chave]; ok {
		return existente
	}

	prestacao = &types.PrestacaoContas{
		SQPrestadorContas: *sqPrestador,
		EleicaoID:         eleicaoID,
		TipoPrestador:     tipos.TipoPrestadorCandidato,
	}
	candidatoID, err := p.garantirIDCandidato(ctx, registro["SQ_PRESTADOR_CONTAS"])
	if err == nil {
		prestacao.CandidatoID = &candidatoID
	} else {
		candidato := &types.Candidato{
			SQCandidato:  *sqPrestador,
			EleicaoID:    eleicaoID,
			NomeCompleto: fmt.Sprintf("__CAND_PLACEHOLDER_%d", *sqPrestador),
		}
		candidato.ID = uuid.Must(uuid.NewV7())
		p.dados.Candidatos[*sqPrestador] = candidato
		p.dados.CandidatosPorID[candidato.ID] = candidato
		prestacao.CandidatoID = &candidato.ID
	}
	prestacao.ID = uuid.Must(uuid.NewV7())
	p.cachearPrestacao(prestacao)
	return prestacao
}

func (p *ProcessadorLeitorCSV) garantirOuCriarPrestacaoOrgao(ctx context.Context, registro map[string]string) *types.PrestacaoContas {
	prestacao, err := p.garantirPrestacaoPorTipoESQ(ctx, tipos.TipoPrestadorOrgaoPartidario, registro["SQ_PRESTADOR_CONTAS"], registro["CD_ELEICAO"])
	if err == nil {
		return prestacao
	}

	sqPrestador := inteiro64Opcional(registro["SQ_PRESTADOR_CONTAS"])
	if sqPrestador == nil {
		return nil
	}

	eleicaoID := p.resolverEleicaoOpcional(registro)
	chave := chavePrestacaoNatural(tipos.TipoPrestadorOrgaoPartidario, eleicaoID, *sqPrestador)
	if existente, ok := p.dados.Prestacoes[chave]; ok {
		return existente
	}

	prestacao = &types.PrestacaoContas{
		SQPrestadorContas: *sqPrestador,
		EleicaoID:         eleicaoID,
		TipoPrestador:     tipos.TipoPrestadorOrgaoPartidario,
	}
	if partido, ok := p.dados.Partidos[int16(*sqPrestador)]; ok {
		partidoID := partido.ID
		prestacao.PartidoID = &partidoID
	} else {
		partido := &types.Partido{
			Numero: int16(*sqPrestador),
			Sigla:  fmt.Sprintf("PH_%d", *sqPrestador),
			Nome:   fmt.Sprintf("__PARTIDO_PLACEHOLDER_%d", *sqPrestador),
		}
		partido.ID = uuid.Must(uuid.NewV7())
		p.dados.Partidos[partido.Numero] = partido
		prestacao.PartidoID = &partido.ID
	}
	prestacao.ID = uuid.Must(uuid.NewV7())
	p.cachearPrestacao(prestacao)
	return prestacao
}
