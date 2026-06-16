package parse

import (
	"context"

	"github.com/danyele/laceu/internal/shared/logger"

	"github.com/google/uuid"

	tipos "github.com/danyele/laceu/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/laceu/internal/shared/types"
)

func (p *ProcessadorLeitorCSV) processarReceitaCandidato(ctx context.Context, caminho string) (int, error) {
	log := logger.New("LeitorCSV: Service: processarReceitaCandidato")
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoCandidatoContratada(ctx, registro)
		if err != nil {
			log.Info("registro de receita de candidato ignorado",
				"caminho", caminho, "linha", numeroLinha, "erro", err)
			return nil
		}
		candidatoID := uuidOpcional(prestacao.CandidatoID, "prestacao_contas.candidato_id")
		if candidatoID == nil {
			log.Info("registro de receita de candidato ignorado - candidato_id vazio",
				"caminho", caminho, "linha", numeroLinha)
			return nil
		}

		doadorID, err := p.garantirDoadorReceitaOpcional(ctx, registro)
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqReceita := inteiro64OuZero(registro["SQ_RECEITA"])
		receita := &types.ReceitaCandidato{
			PrestacaoContasID:        prestacao.ID,
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
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (p *ProcessadorLeitorCSV) processarReceitaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	log := logger.New("LeitorCSV: Service: processarReceitaOrgaoPartidario")
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoOrgaoContratada(ctx, registro)
		if err != nil {
			log.Info("registro de receita de orgao partidario ignorado",
				"caminho", caminho, "linha", numeroLinha, "erro", err)
			return nil
		}
		partidoID := uuidOpcional(prestacao.PartidoID, "prestacao_contas.partido_id")
		if partidoID == nil {
			log.Info("registro de receita de orgao partidario ignorado - partido_id vazio",
				"caminho", caminho, "linha", numeroLinha)
			return nil
		}

		doadorID, err := p.garantirDoadorReceitaOpcional(ctx, registro)
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqReceita := inteiro64OuZero(registro["SQ_RECEITA"])
		receita := &types.ReceitaOrgaoPartidario{
			PrestacaoContasID:        prestacao.ID,
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
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (p *ProcessadorLeitorCSV) processarReceitaCandidatoDoadorOriginario(ctx context.Context, caminho string) (int, error) {
	log := logger.New("LeitorCSV: Service: processarReceitaCandidatoDoadorOriginario")
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoPorTipoESQ(ctx, tipos.TipoPrestadorCandidato, registro["SQ_PRESTADOR_CONTAS"])
		if err != nil {
			log.Info("registro de receita de candidato doador originario ignorado",
				"caminho", caminho, "linha", numeroLinha, "erro", err)
			return nil
		}

		sqReceita := inteiro64OuZero(registro["SQ_RECEITA"])
		receitaID := receitaCandidatoIDPorSQ(p.dados, sqReceita)
		origem := &types.ReceitaDoadorOriginarioCandidato{
			PrestacaoContasID:  prestacao.ID,
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
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (p *ProcessadorLeitorCSV) processarReceitaOrgaoPartidarioDoadorOriginario(ctx context.Context, caminho string) (int, error) {
	log := logger.New("LeitorCSV: Service: processarReceitaOrgaoPartidarioDoadorOriginario")
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoPorTipoESQ(ctx, tipos.TipoPrestadorOrgaoPartidario, registro["SQ_PRESTADOR_CONTAS"])
		if err != nil {
			log.Info("registro de receita de orgao partidario doador originario ignorado",
				"caminho", caminho, "linha", numeroLinha, "erro", err)
			return nil
		}

		sqReceita := inteiro64OuZero(registro["SQ_RECEITA"])
		receitaID := receitaOrgaoIDPorSQ(p.dados, sqReceita)
		origem := &types.ReceitaDoadorOriginarioOrgaoPartidario{
			PrestacaoContasID:        prestacao.ID,
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
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
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
