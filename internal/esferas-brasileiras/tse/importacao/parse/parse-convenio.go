package parse

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/danyele/podp/internal/shared/types"
)

func (p *ProcessadorLeitorCSV) processarConvenioPortal(ctx context.Context, caminho string) (int, error) {
	total := 0

	err := lerArquivoCSV(caminho, func(numeroLinha int, registro map[string]string) error {
		convenio := types.NovoConvenio()
		convenio.NumeroConvenio = textoOpcional(registro["NÚMERO CONVÊNIO"])
		convenio.UF = textoOpcional(registro["UF"])
		convenio.CodigoSIAFIMunicipio = textoOpcional(registro["CÓDIGO SIAFI MUNICÍPIO"])
		convenio.NomeMunicipio = textoOpcional(registro["NOME MUNICÍPIO"])
		convenio.SituacaoConvenio = textoOpcional(registro["SITUAÇÃO CONVÊNIO"])
		convenio.NumeroOriginal = textoOpcional(registro["NÚMERO ORIGINAL"])
		convenio.NumeroProcesso = textoOpcional(registro["NÚMERO PROCESSO DO CONVÊNIO"])
		convenio.ObjetoConvenio = textoOpcional(registro["OBJETO DO CONVÊNIO"])
		convenio.CodigoOrgaoSuperior = textoOpcional(registro["CÓDIGO ÓRGÃO SUPERIOR"])
		convenio.NomeOrgaoSuperior = textoOpcional(registro["NOME ÓRGÃO SUPERIOR"])
		convenio.CodigoOrgaoConcedente = textoOpcional(registro["CÓDIGO ÓRGÃO CONCEDENTE"])
		convenio.NomeOrgaoConcedente = textoOpcional(registro["NOME ÓRGÃO CONCEDENTE"])
		convenio.CodigoUGConcedente = textoOpcional(registro["CÓDIGO UG CONCEDENTE"])
		convenio.NomeUGConcedente = textoOpcional(registro["NOME UG CONCEDENTE"])
		convenio.CodigoConvenente = textoOpcional(registro["CÓDIGO CONVENENTE"])
		convenio.TipoConvenente = textoOpcional(registro["TIPO CONVENENTE"])
		convenio.NomeConvenente = textoOpcional(registro["NOME CONVENENTE"])
		convenio.TipoEnteConvenente = textoOpcional(registro["TIPO ENTE CONVENENTE"])
		convenio.TipoInstrumento = textoOpcional(registro["TIPO INSTRUMENTO"])
		convenio.ValorConvenio = decimalOpcional(registro["VALOR CONVÊNIO"])
		convenio.ValorLiberado = decimalOpcional(registro["VALOR LIBERADO"])
		convenio.DataPublicacao = dataOpcional(registro["DATA PUBLICAÇÃO"])
		convenio.DataInicioVigencia = dataOpcional(registro["DATA INÍCIO VIGÊNCIA"])
		convenio.DataFinalVigencia = dataOpcional(registro["DATA FINAL VIGÊNCIA"])
		convenio.ValorContrapartida = decimalOpcional(registro["VALOR CONTRAPARTIDA"])
		convenio.DataUltimaLiberacao = dataOpcional(registro["DATA ÚLTIMA LIBERAÇÃO"])
		convenio.ValorUltimaLiberacao = decimalOpcional(registro["VALOR ÚLTIMA LIBERAÇÃO"])

		if convenio.NumeroConvenio == "" {
			return fmt.Errorf("linha %d: NÚMERO CONVÊNIO obrigatorio", numeroLinha)
		}

		convenio.ID = uuid.Must(uuid.NewV7())
		p.dados.Convenios = append(p.dados.Convenios, convenio)
		total++
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}
