package parse

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
)

func (p *ProcessadorLeitorCSV) processarBemCandidato(ctx context.Context, caminho string) (int, error) {
	return p.processarCSV(caminho, func(numeroLinha int, registro map[string]string) error {
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
		p.dados.BensCandidato[fmt.Sprintf("%s|%d", candidatoID, numeroOrdem)] = bem
		return nil
	})
}
