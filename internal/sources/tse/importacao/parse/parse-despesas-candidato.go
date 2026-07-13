package parse

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
	tipos "github.com/danyele/podp/internal/sources/tse/importacao/types"
)

func (p *ProcessadorLeitorCSV) processarDespesaContratadaCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoCandidatoContratada(ctx, registro)
		if err != nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado", "caminho", caminho, "linha", numeroLinha,
				"sq_despesa", registro["SQ_DESPESA"], "sq_candidato", registro["SQ_CANDIDATO"], "erro", err)
			return nil
		}
		candidatoID := uuidOpcional(prestacao.CandidatoID, "prestacao_contas.candidato_id")
		if candidatoID == nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado: candidato_id vazio",
				"caminho", caminho, "linha", numeroLinha,
				"sq_despesa", registro["SQ_DESPESA"], "sq_candidato", registro["SQ_CANDIDATO"])
			return nil
		}

		fornecedorID, err := p.garantirFornecedorOpcional(ctx, registro)
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqDespesa := inteiro64OuZero(registro["SQ_DESPESA"])

		valor := decimalOuZero(registro["VR_DESPESA_CONTRATADA"])

		d := &types.DespesaCandidato{
			PrestacaoContasID:      prestacao.ID,
			SQPrestadorContas:      prestacao.SQPrestadorContas,
			CandidatoID:            *candidatoID,
			FornecedorID:           fornecedorID,
			SQDespesa:              sqDespesa,
			TipoRegistro:           tipos.TipoRegistroDespesaContratada,
			TipoDocumento:          textoOpcional(registro["DS_TIPO_DOCUMENTO"]),
			NumeroDocumento:        textoOpcional(registro["NR_DOCUMENTO"]),
			OrigemDespesaCodigo:    inteiroOpcional(registro["CD_ORIGEM_DESPESA"]),
			OrigemDespesaDescricao: textoOpcional(registro["DS_ORIGEM_DESPESA"]),
			DataDespesa:            dataOpcional(registro["DT_DESPESA"]),
			Descricao:              textoComFallback(registro["DS_DESPESA"], "DESPESA SEM DESCRICAO"),
			Valor:                  valor,
		}
		d.ID = uuid.Must(uuid.NewV7())
		p.dados.DespesasCandidato = append(p.dados.DespesasCandidato, d)
		return nil
	})
}

func (p *ProcessadorLeitorCSV) processarDespesaPagaCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoPorTipoESQ(ctx, tipos.TipoPrestadorCandidato, registro["SQ_PRESTADOR_CONTAS"], registro["CD_ELEICAO"])
		if err != nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado", "caminho", caminho, "linha", numeroLinha,
				"sq_despesa", registro["SQ_DESPESA"], "sq_candidato", registro["SQ_CANDIDATO"], "erro", err)
			return nil
		}
		candidatoID := uuidOpcional(prestacao.CandidatoID, "prestacao_contas.candidato_id")
		if candidatoID == nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado: candidato_id vazio",
				"caminho", caminho, "linha", numeroLinha,
				"sq_despesa", registro["SQ_DESPESA"], "sq_candidato", registro["SQ_CANDIDATO"])
			return nil
		}

		fornecedorID, err := p.garantirFornecedorOpcional(ctx, registro)
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqDespesa := inteiro64OuZero(registro["SQ_DESPESA"])

		valor := decimalOuZero(registro["VR_PAGTO_DESPESA"])

		d := &types.DespesaCandidato{
			PrestacaoContasID:        prestacao.ID,
			SQPrestadorContas:        prestacao.SQPrestadorContas,
			CandidatoID:              *candidatoID,
			FornecedorID:             fornecedorID,
			SQDespesa:                sqDespesa,
			TipoRegistro:             tipos.TipoRegistroDespesaPaga,
			TipoDocumento:            textoOpcional(registro["DS_TIPO_DOCUMENTO"]),
			NumeroDocumento:          textoOpcional(registro["NR_DOCUMENTO"]),
			OrigemDespesaCodigo:      inteiroOpcional(registro["CD_ORIGEM_DESPESA"]),
			OrigemDespesaDescricao:   textoOpcional(registro["DS_ORIGEM_DESPESA"]),
			FonteDespesaCodigo:       inteiroOpcional(registro["CD_FONTE_DESPESA"]),
			FonteDespesaDescricao:    textoOpcional(registro["DS_FONTE_DESPESA"]),
			NaturezaDespesaCodigo:    inteiroOpcional(registro["CD_NATUREZA_DESPESA"]),
			NaturezaDespesaDescricao: textoOpcional(registro["DS_NATUREZA_DESPESA"]),
			EspecieRecursoCodigo:     inteiroOpcional(registro["CD_ESPECIE_RECURSO"]),
			EspecieRecursoDescricao:  textoOpcional(registro["DS_ESPECIE_RECURSO"]),
			SQPlanoParcelamento:      inteiro64Opcional(registro["SQ_PARCELAMENTO_DESPESA"]),
			DataDespesa:              dataOpcional(registro["DT_PAGTO_DESPESA"]),
			Descricao:                textoComFallback(registro["DS_DESPESA"], "DESPESA SEM DESCRICAO"),
			Valor:                    valor,
		}
		d.ID = uuid.Must(uuid.NewV7())
		p.dados.DespesasCandidato = append(p.dados.DespesasCandidato, d)
		return nil
	})
}

func (p *ProcessadorLeitorCSV) garantirPrestacaoCandidatoContratada(ctx context.Context, registro map[string]string) (*types.PrestacaoContas, error) {
	eleicaoID, err := p.garantirEleicao(ctx, registro["CD_ELEICAO"], registro["AA_ELEICAO"], registro["CD_TIPO_ELEICAO"], registro["NM_TIPO_ELEICAO"], registro["DS_ELEICAO"], registro["DT_ELEICAO"])
	if err != nil {
		return nil, err
	}

	ufSigla, err := p.garantirUF(ctx, registro["SG_UF"])
	if err != nil {
		return nil, err
	}

	unidadeID, err := p.garantirUnidadeEleitoral(ctx, ufSigla, registro["SG_UE"], registro["NM_UE"])
	if err != nil {
		return nil, err
	}

	_, err = p.garantirPartidoOpcional(ctx, registro["NR_PARTIDO"], registro["SG_PARTIDO"], registro["NM_PARTIDO"],
		"", "", "", "", "", "")
	if err != nil {
		return nil, err
	}

	sqCandidato := inteiro64Opcional(registro["SQ_CANDIDATO"])
	if sqCandidato == nil {
		return nil, fmt.Errorf("SQ_CANDIDATO nao informado ou invalido")
	}

	candidato, existe := p.dados.Candidatos[*sqCandidato]
	if !existe {
		candidato = &types.Candidato{
			SQCandidato:  *sqCandidato,
			EleicaoID:    eleicaoID,
			UFSigla:      ufSigla,
			NomeCompleto: "__PENDENTE_SQ_" + strconv.FormatInt(*sqCandidato, 10),
			NomeUrna:     "__PENDENTE",
		}
		candidato.ID = uuid.Must(uuid.NewV7())
		p.dados.Candidatos[*sqCandidato] = candidato
		p.dados.CandidatosPorID[candidato.ID] = candidato
	}

	sqPrestador := inteiro64Opcional(registro["SQ_PRESTADOR_CONTAS"])

	chavePrest := chavePrestacaoNatural(tipos.TipoPrestadorCandidato, eleicaoID, *sqPrestador)
	if existente, ok := p.dados.Prestacoes[chavePrest]; ok {
		return existente, nil
	}

	prestacao := &types.PrestacaoContas{
		SQPrestadorContas:  *sqPrestador,
		EleicaoID:          eleicaoID,
		CandidatoID:        &candidato.ID,
		UFSigla:            &ufSigla,
		UnidadeEleitoralID: &unidadeID,
		TipoPrestador:      tipos.TipoPrestadorCandidato,
		TipoPrestacao:      textoOpcional(registro["TP_PRESTACAO_CONTAS"]),
		DataPrestacao:      dataOpcional(registro["DT_PRESTACAO_CONTAS"]),
		Turno:              inteiro16Opcional(registro["ST_TURNO"]),
		CNPJPrestadorConta: documentoOpcional(registro["NR_CNPJ_PRESTADOR_CONTA"]),
	}
	prestacao.ID = uuid.Must(uuid.NewV7())
	p.cachearPrestacao(prestacao)
	return prestacao, nil
}
