package parse

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/podp/internal/shared/types"
)

func (p *ProcessadorLeitorCSV) processarDespesaContratadaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoOrgaoContratada(ctx, registro)
		if err != nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado", "caminho", caminho, "linha", numeroLinha,
				"sq_despesa", registro["SQ_DESPESA"], "sq_prestador_contas", registro["SQ_PRESTADOR_CONTAS"], "erro", err)
			return nil
		}
		partidoID := uuidOpcional(prestacao.PartidoID, "prestacao_contas.partido_id")
		if partidoID == nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado: partido_id vazio",
				"caminho", caminho, "linha", numeroLinha,
				"sq_despesa", registro["SQ_DESPESA"], "sq_prestador_contas", registro["SQ_PRESTADOR_CONTAS"])
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
			SQPrestadorContas:      prestacao.SQPrestadorContas,
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
		return nil
	})
}

func (p *ProcessadorLeitorCSV) processarDespesaPagaOrgaoPartidario(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		prestacao, err := p.garantirPrestacaoPorTipoESQ(ctx, tipos.TipoPrestadorOrgaoPartidario, registro["SQ_PRESTADOR_CONTAS"], registro["CD_ELEICAO"])
		if err != nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado", "caminho", caminho, "linha", numeroLinha,
				"sq_despesa", registro["SQ_DESPESA"], "sq_prestador_contas", registro["SQ_PRESTADOR_CONTAS"], "erro", err)
			return nil
		}
		partidoID := uuidOpcional(prestacao.PartidoID, "prestacao_contas.partido_id")
		if partidoID == nil {
			p.RegistrosIgnorados++
			p.log.Warn("registro ignorado: partido_id vazio",
				"caminho", caminho, "linha", numeroLinha,
				"sq_despesa", registro["SQ_DESPESA"], "sq_prestador_contas", registro["SQ_PRESTADOR_CONTAS"])
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
			SQPrestadorContas:        prestacao.SQPrestadorContas,
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
		return nil
	})
}

func (p *ProcessadorLeitorCSV) garantirPrestacaoOrgaoContratada(ctx context.Context, registro map[string]string) (*types.PrestacaoContas, error) {
	anoEleicao := inteiroOpcional(registro["AA_ELEICAO"])

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
	if sqPrestador == nil {
		return nil, fmt.Errorf("SQ_PRESTADOR_CONTAS obrigatorio")
	}

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
