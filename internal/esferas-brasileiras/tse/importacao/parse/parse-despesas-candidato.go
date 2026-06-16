// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo implementa os parsers das planilhas de despesas
// de candidatos (despesas_contratadas_candidatos_ e
// despesas_pagas_candidatos_).
package parse

import (
	"context"
	"fmt"

	"github.com/danyele/laceu/internal/shared/logger"

	"github.com/google/uuid"

	tipos "github.com/danyele/laceu/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/laceu/internal/shared/types"
)

// -----------------------------------------------------------------------------
// Processamento de despesas contratadas de candidatos.
// -----------------------------------------------------------------------------

// processarDespesaContratadaCandidato percorre o CSV de despesas contratadas
// de candidatos. Para cada linha, garante a prestação de contas (criando o
// candidato sintético se necessário), o fornecedor, e monta a
// DespesaCandidato com tipo CONTRATADA.
func (p *ProcessadorLeitorCSV) processarDespesaContratadaCandidato(ctx context.Context, caminho string) (int, error) {
	log := logger.New("LeitorCSV: Service: processarDespesaContratadaCandidato")
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoCandidatoContratada(ctx, registro)
		if err != nil {
			log.Info("registro de despesa contratada de candidato ignorado",
				"caminho", caminho, "linha", numeroLinha, "erro", err)
			return nil
		}
		candidatoID := uuidOpcional(prestacao.CandidatoID, "prestacao_contas.candidato_id")
		if candidatoID == nil {
			log.Info("registro de despesa contratada de candidato ignorado - candidato_id vazio",
				"caminho", caminho, "linha", numeroLinha)
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
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// -----------------------------------------------------------------------------
// Processamento de despesas pagas de candidatos.
// -----------------------------------------------------------------------------

// processarDespesaPagaCandidato percorre o CSV de despesas pagas de
// candidatos. Requer que as despesas contratadas já tenham sido importadas
// (para existir a prestação de contas). Monta DespesaCandidato com tipo PAGA.
func (p *ProcessadorLeitorCSV) processarDespesaPagaCandidato(ctx context.Context, caminho string) (int, error) {
	log := logger.New("LeitorCSV: Service: processarDespesaPagaCandidato")
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoPorTipoESQ(ctx, tipos.TipoPrestadorCandidato, registro["SQ_PRESTADOR_CONTAS"])
		if err != nil {
			log.Info("registro de despesa paga de candidato ignorado",
				"caminho", caminho, "linha", numeroLinha, "erro", err)
			return nil // Pula a linha, não falha o arquivo
		}
		candidatoID := uuidOpcional(prestacao.CandidatoID, "prestacao_contas.candidato_id")
		if candidatoID == nil {
			log.Info("registro de despesa paga de candidato ignorado - candidato_id vazio",
				"caminho", caminho, "linha", numeroLinha)
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
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// -----------------------------------------------------------------------------
// Garantia de prestação de contas de candidato para despesas contratadas.
// -----------------------------------------------------------------------------

// garantirPrestacaoCandidatoContratada busca ou cria uma prestação de contas
// de candidato a partir dos dados da linha. Diferente do processo normal de
// consulta, aqui o candidato pode ser criado sinteticamente se ainda não
// existir no mapa (para arquivos de despesa que precedem a consulta).
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
		const sqCandidatoNaoEncontrado int64 = 1
		dummy, ok := p.dados.Candidatos[sqCandidatoNaoEncontrado]
		if !ok {
			dummy = &types.Candidato{
				SQCandidato:  sqCandidatoNaoEncontrado,
				EleicaoID:    eleicaoID,
				UFSigla:      ufSigla,
				NomeCompleto: "CANDIDATO NAO ENCONTRADO",
			}
			dummy.ID = uuid.Must(uuid.NewV7())
			p.dados.Candidatos[sqCandidatoNaoEncontrado] = dummy
		}
		// candidato nao encontrado: usamos dummy SQ=1 sem log para reduzir ruido
		candidato = dummy
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
