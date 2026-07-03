package types

import (
	"time"

	"github.com/google/uuid"
)

type Convenio struct {
	ModeloBase
	NumeroConvenio        string     `json:"numero_convenio"`
	UF                    string     `json:"uf,omitempty"`
	CodigoSIAFIMunicipio  string     `json:"codigo_siafi_municipio,omitempty"`
	NomeMunicipio         string     `json:"nome_municipio,omitempty"`
	SituacaoConvenio      string     `json:"situacao_convenio,omitempty"`
	NumeroOriginal        string     `json:"numero_original,omitempty"`
	NumeroProcesso        string     `json:"numero_processo,omitempty"`
	ObjetoConvenio        string     `json:"objeto_convenio,omitempty"`
	CodigoOrgaoSuperior   string     `json:"codigo_orgao_superior,omitempty"`
	NomeOrgaoSuperior     string     `json:"nome_orgao_superior,omitempty"`
	CodigoOrgaoConcedente string     `json:"codigo_orgao_concedente,omitempty"`
	NomeOrgaoConcedente   string     `json:"nome_orgao_concedente,omitempty"`
	CodigoUGConcedente    string     `json:"codigo_ug_concedente,omitempty"`
	NomeUGConcedente      string     `json:"nome_ug_concedente,omitempty"`
	CodigoConvenente      string     `json:"codigo_convenente,omitempty"`
	TipoConvenente        string     `json:"tipo_convenente,omitempty"`
	NomeConvenente        string     `json:"nome_convenente,omitempty"`
	TipoEnteConvenente    string     `json:"tipo_ente_convenente,omitempty"`
	TipoInstrumento       string     `json:"tipo_instrumento,omitempty"`
	ValorConvenio         *float64   `json:"valor_convenio,omitempty"`
	ValorLiberado         *float64   `json:"valor_liberado,omitempty"`
	DataPublicacao        *time.Time `json:"data_publicacao,omitempty"`
	DataInicioVigencia    *time.Time `json:"data_inicio_vigencia,omitempty"`
	DataFinalVigencia     *time.Time `json:"data_final_vigencia,omitempty"`
	ValorContrapartida    *float64   `json:"valor_contrapartida,omitempty"`
	DataUltimaLiberacao   *time.Time `json:"data_ultima_liberacao,omitempty"`
	ValorUltimaLiberacao  *float64   `json:"valor_ultima_liberacao,omitempty"`
}

func NovoConvenio() *Convenio {
	return &Convenio{
		ModeloBase: ModeloBase{ID: uuid.Must(uuid.NewV7())},
	}
}
