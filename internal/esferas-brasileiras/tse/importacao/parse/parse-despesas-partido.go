// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo implementa os parsers das planilhas de despesas
// de órgãos partidários (despesas_contratadas_orgaos_partidarios_ e
// despesas_pagas_orgaos_partidarios_).
package parse

import (
	"context"
	"fmt"
	"strconv"

	"github.com/danyele/laceu/internal/shared/logger"

	"github.com/google/uuid"

	tipos "github.com/danyele/laceu/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/laceu/internal/shared/types"
)

// -----------------------------------------------------------------------------
// Processamento de despesas contratadas de órgãos partidários.
// -----------------------------------------------------------------------------

// processarDespesaContratadaOrgaoPartidario percorre o CSV de despesas
// contratadas de órgãos partidários. Para cada linha, garante a prestação
// de contas e monta a DespesaOrgaoPartidario com tipo CONTRATADA.
func (p *ProcessadorLeitorCSV) processarDespesaContratadaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	log := logger.New("LeitorCSV: Service: processarDespesaContratadaOrgaoPartidario")
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoOrgaoContratada(ctx, registro)
		if err != nil {
			log.Info("registro de despesa contratada de orgao partidario ignorado",
				"caminho", caminho, "linha", numeroLinha, "erro", err)
			return nil
		}
		partidoID := uuidOpcional(prestacao.PartidoID, "prestacao_contas.partido_id")
		if partidoID == nil {
			log.Info("registro de despesa contratada de orgao partidario ignorado - partido_id vazio",
				"caminho", caminho, "linha", numeroLinha)
			return nil
		}

		fornecedorID, err := p.garantirFornecedorOpcional(ctx, registro)
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqDespesa := inteiro64OuZero(registro["SQ_DESPESA"])

		valor := decimalOuZero(registro["VR_DESPESA_CONTRATADA"])

		d := &types.DespesaOrgaoPartidario{
			PrestacaoContasID:      prestacao.ID,
			PartidoID:              *partidoID,
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
		p.dados.DespesasOrgaoPartidario = append(p.dados.DespesasOrgaoPartidario, d)
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// -----------------------------------------------------------------------------
// Processamento de despesas pagas de órgãos partidários.
// -----------------------------------------------------------------------------

// processarDespesaPagaOrgaoPartidario percorre o CSV de despesas pagas de
// órgãos partidários. Requer que as despesas contratadas já tenham sido
// importadas. Monta DespesaOrgaoPartidario com tipo PAGA.
func (p *ProcessadorLeitorCSV) processarDespesaPagaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	log := logger.New("LeitorCSV: Service: processarDespesaPagaOrgaoPartidario")
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoPorTipoESQ(ctx, tipos.TipoPrestadorOrgaoPartidario, registro["SQ_PRESTADOR_CONTAS"])
		if err != nil {
			log.Info("registro de despesa paga de orgao partidario ignorado",
				"caminho", caminho, "linha", numeroLinha, "erro", err)
			return nil
		}
		partidoID := uuidOpcional(prestacao.PartidoID, "prestacao_contas.partido_id")
		if partidoID == nil {
			log.Info("registro de despesa paga de orgao partidario ignorado - partido_id vazio",
				"caminho", caminho, "linha", numeroLinha)
			return nil
		}

		fornecedorID, err := p.garantirFornecedorOpcional(ctx, registro)
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		sqDespesa := inteiro64OuZero(registro["SQ_DESPESA"])
		valor := decimalOuZero(registro["VR_PAGTO_DESPESA"])

		d := &types.DespesaOrgaoPartidario{
			PrestacaoContasID:        prestacao.ID,
			PartidoID:                *partidoID,
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
		p.dados.DespesasOrgaoPartidario = append(p.dados.DespesasOrgaoPartidario, d)
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// -----------------------------------------------------------------------------
// Garantia de prestação de contas de órgão partidário para despesas
// contratadas.
// -----------------------------------------------------------------------------

// garantirPrestacaoOrgaoContratada busca ou cria uma prestação de contas
// de órgão partidário. Diferente da prestação de candidato, usa
// CD_ELEICAO com fallback para AA_ELEICAO (pois o código da eleição pode
// vir apenas no ano) e a UF é opcional.
func (p *ProcessadorLeitorCSV) garantirPrestacaoOrgaoContratada(ctx context.Context, registro map[string]string) (*types.PrestacaoContas, error) {
	anoEleicao := inteiroOpcional(registro["AA_ELEICAO"])

	// Para órgãos partidários, o código da eleição pode ser apenas o ano;
	// CD_ELEICAO tem precedência sobre AA_ELEICAO quando presente.
	codigoEleicao := anoEleicao
	if codigo := inteiroOpcional(registro["CD_ELEICAO"]); codigo != nil {
		codigoEleicao = codigo
	}

	eleicaoID, err := p.garantirEleicao(ctx, strconv.Itoa(*codigoEleicao), registro["AA_ELEICAO"], registro["CD_TIPO_ELEICAO"], registro["NM_TIPO_ELEICAO"], textoComFallback(registro["DS_ELEICAO"], fmt.Sprintf("Eleicao %s", registro["AA_ELEICAO"])), registro["DT_ELEICAO"])
	if err != nil {
		return nil, err
	}

	ufSigla, err := p.garantirUFOpcional(ctx, registro["SG_UF"])
	if err != nil {
		return nil, err
	}

	unidadeID, err := p.garantirUnidadeEleitoralOpcional(ctx, ufSigla, registro["SG_UE"], registro["NM_UE"])
	if err != nil {
		return nil, err
	}

	partidoID, err := p.garantirPartidoObrigatorio(ctx, registro["NR_PARTIDO"], registro["SG_PARTIDO"], registro["NM_PARTIDO"],
		"", "", "", "", "", "")
	if err != nil {
		return nil, err
	}

	sqPrestador := inteiro64Opcional(registro["SQ_PRESTADOR_CONTAS"])

	chavePrest := chavePrestacaoNatural(tipos.TipoPrestadorOrgaoPartidario, eleicaoID, *sqPrestador)
	if existente, ok := p.dados.Prestacoes[chavePrest]; ok {
		return existente, nil
	}

	prestacao := &types.PrestacaoContas{
		SQPrestadorContas:         *sqPrestador,
		EleicaoID:                 eleicaoID,
		PartidoID:                 &partidoID,
		UFSigla:                   ufSigla,
		UnidadeEleitoralID:        unidadeID,
		TipoPrestador:             tipos.TipoPrestadorOrgaoPartidario,
		TipoPrestacao:             textoOpcional(registro["TP_PRESTACAO_CONTAS"]),
		DataPrestacao:             dataOpcional(registro["DT_PRESTACAO_CONTAS"]),
		CNPJPrestadorConta:        documentoOpcional(registro["NR_CNPJ_PRESTADOR_CONTA"]),
		EsferaPartidariaCodigo:    textoOpcional(registro["CD_ESFERA_PARTIDARIA"]),
		EsferaPartidariaDescricao: textoOpcional(registro["DS_ESFERA_PARTIDARIA"]),
	}
	prestacao.ID = uuid.Must(uuid.NewV7())
	p.cachearPrestacao(prestacao)
	return prestacao, nil
}
