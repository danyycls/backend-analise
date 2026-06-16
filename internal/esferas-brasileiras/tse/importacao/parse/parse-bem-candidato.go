// Pacote service contém funções auxiliares para a importação de dados
// eleitorais. Este arquivo implementa o parser da planilha de bens de
// candidatos (bem_candidato_).
package parse

import (
	"context"

	"github.com/google/uuid"

	"github.com/danyele/laceu/internal/shared/types"
)

// -----------------------------------------------------------------------------
// Processamento da planilha de bens de candidatos.
// -----------------------------------------------------------------------------

// processarBemCandidato percorre o CSV de bens de candidatos linha a linha.
// Para cada registro, localiza o candidato pelo SQ_CANDIDATO (deve ter sido
// importado antes), garante o tipo de bem e monta o BemCandidato.
func (p *ProcessadorLeitorCSV) processarBemCandidato(ctx context.Context, caminho string) (int, error) {
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		candidatoID, err := p.garantirIDCandidato(ctx, registro["SQ_CANDIDATO"])
		if err != nil {
			return erroLinha(numeroLinha, err)
		}

		var numeroOrdem int
		if v := inteiroOpcional(registro["NR_ORDEM_BEM_CANDIDATO"]); v != nil {
			numeroOrdem = *v
		}

		var valor float64
		if v := decimalOpcional(registro["VR_BEM_CANDIDATO"]); v != nil {
			valor = *v
		}

		descricao := textoOpcional(registro["DS_BEM_CANDIDATO"])

		bem := &types.BemCandidato{
			CandidatoID:           candidatoID,
			TipoBemCodigo:         inteiroOpcional(registro["CD_TIPO_BEM_CANDIDATO"]),
			TipoBemNome:           textoOpcional(registro["DS_TIPO_BEM_CANDIDATO"]),
			NumeroOrdem:           numeroOrdem,
			Descricao:             descricao,
			Valor:                 valor,
			DataUltimaAtualizacao: dataOpcional(registro["DT_ULT_ATUAL_BEM_CANDIDATO"]),
			HoraUltimaAtualizacao: textoOpcional(registro["HH_ULT_ATUAL_BEM_CANDIDATO"]),
		}
		bem.ID = uuid.Must(uuid.NewV7())
		p.dados.BensCandidato = append(p.dados.BensCandidato, bem)
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}
