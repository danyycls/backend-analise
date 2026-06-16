package types

import (
	"time"

	"github.com/google/uuid"
)

type Candidato struct {
	ModeloBase
	SQCandidato                  int64          `json:"sq_candidato"`
	EleicaoID                    uuid.UUID      `json:"eleicao_id"`
	UFSigla                      string         `json:"sg_uf"`
	PartidoID                    *uuid.UUID     `json:"partido_id,omitempty"`
	CargoCodigo                  *int           `json:"cargo_codigo,omitempty"`
	CargoNome                    string         `json:"cargo_nome,omitempty"`
	GeneroDescricao              string         `json:"genero_descricao,omitempty"`
	CorRacaDescricao             string         `json:"cor_raca_descricao,omitempty"`
	EstadoCivilNome              string         `json:"estado_civil_nome,omitempty"`
	GrauInstrucaoNome            string         `json:"grau_instrucao_nome,omitempty"`
	OcupacaoCodigo               *int           `json:"ocupacao_codigo,omitempty"`
	OcupacaoNome                 string         `json:"ocupacao_nome,omitempty"`
	NumeroCandidato              *int           `json:"numero_candidato,omitempty"`
	CPF                          string         `json:"cpf,omitempty"`
	CPFVice                      string         `json:"cpf_vice,omitempty"`
	NomeCompleto                 string         `json:"nome_completo"`
	NomeUrna                     string         `json:"nome_urna,omitempty"`
	NomeSocial                   string         `json:"nome_social,omitempty"`
	DataNascimento               *time.Time     `json:"data_nascimento,omitempty"`
	SituacaoTotalizacaoDescricao string         `json:"situacao_totalizacao_descricao,omitempty"`
	Bens                         []BemCandidato `json:"bens,omitempty"`
}

type BemCandidato struct {
	ModeloBase
	CandidatoID           uuid.UUID  `json:"candidato_id"`
	TipoBemCodigo         *int       `json:"tipo_bem_codigo,omitempty"`
	TipoBemNome           string     `json:"tipo_bem_nome,omitempty"`
	NumeroOrdem           int        `json:"numero_ordem"`
	Descricao             string     `json:"descricao"`
	Valor                 float64    `json:"valor"`
	DataUltimaAtualizacao *time.Time `json:"data_ultima_atualizacao,omitempty"`
	HoraUltimaAtualizacao string     `json:"hora_ultima_atualizacao,omitempty"`
}

type Eleicao struct {
	ModeloBase
	Ano               int16      `json:"ano"`
	CodigoTSE         int        `json:"codigo_tse"`
	CodigoTipoEleicao *int       `json:"codigo_tipo_eleicao,omitempty"`
	NomeTipoEleicao   string     `json:"nome_tipo_eleicao,omitempty"`
	Descricao         string     `json:"descricao"`
	DataEleicao       *time.Time `json:"data_eleicao,omitempty"`
}

type UnidadeEleitoral struct {
	ModeloBase
	UFSigla   string `json:"sg_uf"`
	CodigoTSE string `json:"codigo_tse"`
	Nome      string `json:"nome"`
}

type Partido struct {
	ModeloBase
	Numero              int16  `json:"numero"`
	Sigla               string `json:"sigla"`
	Nome                string `json:"nome"`
	FederacaoCodigoTSE  *int64 `json:"federacao_codigo_tse,omitempty"`
	FederacaoSigla      string `json:"federacao_sigla,omitempty"`
	FederacaoNome       string `json:"federacao_nome,omitempty"`
	ColigacaoCodigoTSE  *int64 `json:"coligacao_codigo_tse,omitempty"`
	ColigacaoNome       string `json:"coligacao_nome,omitempty"`
	ColigacaoComposicao string `json:"coligacao_composicao,omitempty"`
}

type Doador struct {
	ModeloBase
	CPFCNPJ                    string  `json:"cpf_cnpj"`
	Nome                       string  `json:"nome"`
	NomeRFB                    string  `json:"nome_rfb,omitempty"`
	CNAECodigo                 string  `json:"cnae_codigo,omitempty"`
	CNAEDescricao              string  `json:"cnae_descricao,omitempty"`
	EsferaPartidariaCodigo     string  `json:"esfera_partidaria_codigo,omitempty"`
	EsferaPartidariaDescricao  string  `json:"esfera_partidaria_descricao,omitempty"`
	UFSigla                    *string `json:"sg_uf,omitempty"`
	MunicipioNome              string  `json:"municipio_nome,omitempty"`
	SQCandidatoRelacionado     *int64  `json:"sq_candidato_relacionado,omitempty"`
	NumeroCandidatoRelacionado *int    `json:"numero_candidato_relacionado,omitempty"`
	CargoCodigoRelacionado     *int    `json:"cargo_codigo_relacionado,omitempty"`
	CargoDescricaoRelacionada  string  `json:"cargo_descricao_relacionada,omitempty"`
	PartidoNumeroRelacionado   *int16  `json:"partido_numero_relacionado,omitempty"`
	PartidoSiglaRelacionado    string  `json:"partido_sigla_relacionado,omitempty"`
	PartidoNomeRelacionado     string  `json:"partido_nome_relacionado,omitempty"`
}

type Fornecedor struct {
	ModeloBase
	CPFCNPJ                    string  `json:"cpf_cnpj"`
	Nome                       string  `json:"nome"`
	NomeRFB                    string  `json:"nome_rfb,omitempty"`
	TipoFornecedorCodigo       *int    `json:"tipo_fornecedor_codigo,omitempty"`
	TipoFornecedorDescricao    string  `json:"tipo_fornecedor_descricao,omitempty"`
	CNAECodigo                 string  `json:"cnae_codigo,omitempty"`
	CNAEDescricao              string  `json:"cnae_descricao,omitempty"`
	EsferaPartidariaCodigo     string  `json:"esfera_partidaria_codigo,omitempty"`
	EsferaPartidariaDescricao  string  `json:"esfera_partidaria_descricao,omitempty"`
	UFSigla                    *string `json:"sg_uf,omitempty"`
	MunicipioNome              string  `json:"municipio_nome,omitempty"`
	SQCandidatoRelacionado     *int64  `json:"sq_candidato_relacionado,omitempty"`
	NumeroCandidatoRelacionado *int    `json:"numero_candidato_relacionado,omitempty"`
	CargoCodigoRelacionado     *int    `json:"cargo_codigo_relacionado,omitempty"`
	CargoDescricaoRelacionada  string  `json:"cargo_descricao_relacionada,omitempty"`
	PartidoNumeroRelacionado   *int16  `json:"partido_numero_relacionado,omitempty"`
	PartidoSiglaRelacionado    string  `json:"partido_sigla_relacionado,omitempty"`
	PartidoNomeRelacionado     string  `json:"partido_nome_relacionado,omitempty"`
}

type PrestacaoContas struct {
	ModeloBase
	SQPrestadorContas         int64      `json:"sq_prestador_contas"`
	EleicaoID                 uuid.UUID  `json:"eleicao_id"`
	CandidatoID               *uuid.UUID `json:"candidato_id,omitempty"`
	PartidoID                 *uuid.UUID `json:"partido_id,omitempty"`
	UFSigla                   *string    `json:"sg_uf,omitempty"`
	UnidadeEleitoralID        *uuid.UUID `json:"unidade_eleitoral_id,omitempty"`
	TipoPrestador             string     `json:"tipo_prestador"`
	TipoPrestacao             string     `json:"tipo_prestacao,omitempty"`
	DataPrestacao             *time.Time `json:"data_prestacao,omitempty"`
	Turno                     *int16     `json:"turno,omitempty"`
	CNPJPrestadorConta        string     `json:"cnpj_prestador_conta,omitempty"`
	EsferaPartidariaCodigo    string     `json:"esfera_partidaria_codigo,omitempty"`
	EsferaPartidariaDescricao string     `json:"esfera_partidaria_descricao,omitempty"`
}

type DespesaCandidato struct {
	ModeloBase
	PrestacaoContasID        uuid.UUID  `json:"prestacao_contas_id"`
	CandidatoID              uuid.UUID  `json:"candidato_id"`
	FornecedorID             *uuid.UUID `json:"fornecedor_id,omitempty"`
	SQDespesa                int64      `json:"sq_despesa"`
	TipoRegistro             string     `json:"tipo_registro"`
	TipoDocumento            string     `json:"tipo_documento,omitempty"`
	NumeroDocumento          string     `json:"numero_documento,omitempty"`
	OrigemDespesaCodigo      *int       `json:"origem_despesa_codigo,omitempty"`
	OrigemDespesaDescricao   string     `json:"origem_despesa_descricao,omitempty"`
	FonteDespesaCodigo       *int       `json:"fonte_despesa_codigo,omitempty"`
	FonteDespesaDescricao    string     `json:"fonte_despesa_descricao,omitempty"`
	NaturezaDespesaCodigo    *int       `json:"natureza_despesa_codigo,omitempty"`
	NaturezaDespesaDescricao string     `json:"natureza_despesa_descricao,omitempty"`
	EspecieRecursoCodigo     *int       `json:"especie_recurso_codigo,omitempty"`
	EspecieRecursoDescricao  string     `json:"especie_recurso_descricao,omitempty"`
	SQPlanoParcelamento      *int64     `json:"sq_parcelamento_despesa,omitempty"`
	DataDespesa              *time.Time `json:"data_despesa,omitempty"`
	Descricao                string     `json:"descricao"`
	Valor                    float64    `json:"valor"`
}

type DespesaOrgaoPartidario struct {
	ModeloBase
	PrestacaoContasID        uuid.UUID  `json:"prestacao_contas_id"`
	PartidoID                uuid.UUID  `json:"partido_id"`
	FornecedorID             *uuid.UUID `json:"fornecedor_id,omitempty"`
	SQDespesa                int64      `json:"sq_despesa"`
	TipoRegistro             string     `json:"tipo_registro"`
	TipoDocumento            string     `json:"tipo_documento,omitempty"`
	NumeroDocumento          string     `json:"numero_documento,omitempty"`
	OrigemDespesaCodigo      *int       `json:"origem_despesa_codigo,omitempty"`
	OrigemDespesaDescricao   string     `json:"origem_despesa_descricao,omitempty"`
	FonteDespesaCodigo       *int       `json:"fonte_despesa_codigo,omitempty"`
	FonteDespesaDescricao    string     `json:"fonte_despesa_descricao,omitempty"`
	NaturezaDespesaCodigo    *int       `json:"natureza_despesa_codigo,omitempty"`
	NaturezaDespesaDescricao string     `json:"natureza_despesa_descricao,omitempty"`
	EspecieRecursoCodigo     *int       `json:"especie_recurso_codigo,omitempty"`
	EspecieRecursoDescricao  string     `json:"especie_recurso_descricao,omitempty"`
	SQPlanoParcelamento      *int64     `json:"sq_parcelamento_despesa,omitempty"`
	DataDespesa              *time.Time `json:"data_despesa,omitempty"`
	Descricao                string     `json:"descricao"`
	Valor                    float64    `json:"valor"`
}

type ReceitaCandidato struct {
	ModeloBase
	PrestacaoContasID        uuid.UUID  `json:"prestacao_contas_id"`
	CandidatoID              uuid.UUID  `json:"candidato_id"`
	DoadorID                 *uuid.UUID `json:"doador_id,omitempty"`
	SQReceita                int64      `json:"sq_receita"`
	FonteReceitaCodigo       *int       `json:"fonte_receita_codigo,omitempty"`
	FonteReceitaDescricao    string     `json:"fonte_receita_descricao,omitempty"`
	OrigemReceitaCodigo      *int       `json:"origem_receita_codigo,omitempty"`
	OrigemReceitaDescricao   string     `json:"origem_receita_descricao,omitempty"`
	NaturezaReceitaCodigo    *int       `json:"natureza_receita_codigo,omitempty"`
	NaturezaReceitaDescricao string     `json:"natureza_receita_descricao,omitempty"`
	EspecieReceitaCodigo     *int       `json:"especie_receita_codigo,omitempty"`
	EspecieReceitaDescricao  string     `json:"especie_receita_descricao,omitempty"`
	NumeroReciboDoacao       string     `json:"numero_recibo_doacao,omitempty"`
	NumeroDocumentoDoacao    string     `json:"numero_documento_doacao,omitempty"`
	DataReceita              *time.Time `json:"data_receita,omitempty"`
	Descricao                string     `json:"descricao"`
	Valor                    float64    `json:"valor"`
	NaturezaRecursoEstimavel string     `json:"natureza_recurso_estimavel,omitempty"`
	Genero                   string     `json:"genero,omitempty"`
	CorRaca                  string     `json:"cor_raca,omitempty"`
}

type ReceitaOrgaoPartidario struct {
	ModeloBase
	PrestacaoContasID        uuid.UUID  `json:"prestacao_contas_id"`
	PartidoID                uuid.UUID  `json:"partido_id"`
	DoadorID                 *uuid.UUID `json:"doador_id,omitempty"`
	SQReceita                int64      `json:"sq_receita"`
	FonteReceitaCodigo       *int       `json:"fonte_receita_codigo,omitempty"`
	FonteReceitaDescricao    string     `json:"fonte_receita_descricao,omitempty"`
	OrigemReceitaCodigo      *int       `json:"origem_receita_codigo,omitempty"`
	OrigemReceitaDescricao   string     `json:"origem_receita_descricao,omitempty"`
	NaturezaReceitaCodigo    *int       `json:"natureza_receita_codigo,omitempty"`
	NaturezaReceitaDescricao string     `json:"natureza_receita_descricao,omitempty"`
	EspecieReceitaCodigo     *int       `json:"especie_receita_codigo,omitempty"`
	EspecieReceitaDescricao  string     `json:"especie_receita_descricao,omitempty"`
	NumeroReciboDoacao       string     `json:"numero_recibo_doacao,omitempty"`
	NumeroDocumentoDoacao    string     `json:"numero_documento_doacao,omitempty"`
	DataReceita              *time.Time `json:"data_receita,omitempty"`
	Descricao                string     `json:"descricao"`
	Valor                    float64    `json:"valor"`
}

type ReceitaDoadorOriginarioCandidato struct {
	ModeloBase
	PrestacaoContasID  uuid.UUID  `json:"prestacao_contas_id"`
	ReceitaCandidatoID *uuid.UUID `json:"receita_candidato_id,omitempty"`
	SQReceita          int64      `json:"sq_receita"`
	DocumentoDoador    string     `json:"documento_doador,omitempty"`
	NomeDoador         string     `json:"nome_doador"`
	NomeDoadorRFB      string     `json:"nome_doador_rfb,omitempty"`
	TipoDoador         string     `json:"tipo_doador,omitempty"`
	CNAECodigo         string     `json:"cnae_codigo,omitempty"`
	CNAEDescricao      string     `json:"cnae_descricao,omitempty"`
	DataReceita        *time.Time `json:"data_receita,omitempty"`
	Descricao          string     `json:"descricao"`
	Valor              float64    `json:"valor"`
}

type ReceitaDoadorOriginarioOrgaoPartidario struct {
	ModeloBase
	PrestacaoContasID        uuid.UUID  `json:"prestacao_contas_id"`
	ReceitaOrgaoPartidarioID *uuid.UUID `json:"receita_orgao_partidario_id,omitempty"`
	SQReceita                int64      `json:"sq_receita"`
	DocumentoDoador          string     `json:"documento_doador,omitempty"`
	NomeDoador               string     `json:"nome_doador"`
	NomeDoadorRFB            string     `json:"nome_doador_rfb,omitempty"`
	TipoDoador               string     `json:"tipo_doador,omitempty"`
	CNAECodigo               string     `json:"cnae_codigo,omitempty"`
	CNAEDescricao            string     `json:"cnae_descricao,omitempty"`
	DataReceita              *time.Time `json:"data_receita,omitempty"`
	Descricao                string     `json:"descricao"`
	Valor                    float64    `json:"valor"`
}
