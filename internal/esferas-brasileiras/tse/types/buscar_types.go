package tsetypes

type CandidatoDetalhado struct {
	Candidato          *Candidato        `json:"candidato"`
	Eleicao            *Eleicao          `json:"eleicao,omitempty"`
	UnidadeEleitoral   *UnidadeEleitoral `json:"unidade_eleitoral,omitempty"`
	Partido            *Partido          `json:"partido,omitempty"`
	DescricaoDeVinculo string            `json:"descricao_de_vinculo,omitempty"`
}

type PrestacaoDetalhada struct {
	PrestacaoContas         *PrestacaoContas          `json:"prestacao_contas"`
	DespesasCandidato       []*DespesaCandidato       `json:"despesas_candidato,omitempty"`
	DespesasOrgaoPartidario []*DespesaOrgaoPartidario `json:"despesas_orgao_partidario,omitempty"`
	Eleicao                 *Eleicao                  `json:"eleicao,omitempty"`
	Partido                 *Partido                  `json:"partido,omitempty"`
	UnidadeEleitoral        *UnidadeEleitoral         `json:"unidade_eleitoral,omitempty"`
	DescricaoDeVinculo      string                    `json:"descricao_de_vinculo,omitempty"`
}

type FornecedorEnriquecido struct {
	CPFCNPJ                    string              `bson:"cpfcnpj" json:"cpf_cnpj"`
	Nome                       string              `bson:"nome" json:"nome"`
	NomeRFB                    string              `bson:"nomerfb,omitempty" json:"nome_rfb,omitempty"`
	TipoFornecedorCodigo       *int                `bson:"tipofornecedorcodigo,omitempty" json:"tipo_fornecedor_codigo,omitempty"`
	TipoFornecedorDescricao    string              `bson:"tipofornecedordescricao,omitempty" json:"tipo_fornecedor_descricao,omitempty"`
	CNAECodigo                 string              `bson:"cnaecodigo,omitempty" json:"cnae_codigo,omitempty"`
	CNAEDescricao              string              `bson:"cnaedescricao,omitempty" json:"cnae_descricao,omitempty"`
	EsferaPartidariaCodigo     string              `bson:"esferapartidariacodigo,omitempty" json:"esfera_partidaria_codigo,omitempty"`
	EsferaPartidariaDescricao  string              `bson:"esferapartidariadescricao,omitempty" json:"esfera_partidaria_descricao,omitempty"`
	UFSigla                    *string             `bson:"ufsigla,omitempty" json:"sg_uf,omitempty"`
	MunicipioNome              string              `bson:"municipionome,omitempty" json:"municipio_nome,omitempty"`
	SQCandidatoRelacionado     *int64              `bson:"sqcandidatorelacionado,omitempty" json:"sq_candidato_relacionado,omitempty"`
	NumeroCandidatoRelacionado *int                `bson:"numerocandidatorelacionado,omitempty" json:"numero_candidato_relacionado,omitempty"`
	CargoCodigoRelacionado     *int                `bson:"cargocodigorelacionado,omitempty" json:"cargo_codigo_relacionado,omitempty"`
	CargoDescricaoRelacionada  string              `bson:"cargodescricaorelacionada,omitempty" json:"cargo_descricao_relacionada,omitempty"`
	PartidoNumeroRelacionado   *int16              `bson:"partidonumerorelacionado,omitempty" json:"partido_numero_relacionado,omitempty"`
	PartidoSiglaRelacionado    string              `bson:"partidosiglarelacionado,omitempty" json:"partido_sigla_relacionado,omitempty"`
	PartidoNomeRelacionado     string              `bson:"partidonomerelacionado,omitempty" json:"partido_nome_relacionado,omitempty"`
	Enriquecimento             *FornecedorOpenCNPJ `bson:"enriquecimento,omitempty" json:"enriquecimento,omitempty"`
}

type FornecedorDetalhado struct {
	Fornecedor              FornecedorEnriquecido             `bson:"fornecedor" json:"fornecedor"`
	DespesasCandidato       []DespesaCandidatoDetalhada       `bson:"despesascandidato,omitempty" json:"despesas_candidato,omitempty"`
	DespesasOrgaoPartidario []DespesaOrgaoPartidarioDetalhada `bson:"despesasorgaopartidario,omitempty" json:"despesas_orgao_partidario,omitempty"`
	DescricaoDeVinculo      string                            `bson:"descricaodevinculo,omitempty" json:"descricao_de_vinculo,omitempty"`
}

type DespesaCandidatoDetalhada struct {
	Despesa            *DespesaCandidato `bson:"despesa" json:"despesa"`
	SQCandidato        int64             `bson:"sqcandidato,omitempty" json:"sq_candidato,omitempty"`
	SQPrestacao        int64             `bson:"sqprestacao,omitempty" json:"sq_prestacao,omitempty"`
	DescricaoDeVinculo string            `bson:"descricaodevinculo,omitempty" json:"descricao_de_vinculo,omitempty"`
	Candidato          *Candidato        `bson:"candidato,omitempty" json:"candidato,omitempty"`
}

type DespesaOrgaoPartidarioDetalhada struct {
	Despesa            *DespesaOrgaoPartidario `bson:"despesa" json:"despesa"`
	PartidoNumero      int16                   `bson:"partidonumero,omitempty" json:"partido_numero,omitempty"`
	PartidoNome        string                  `bson:"partidonome,omitempty" json:"partido_nome,omitempty"`
	SQPrestacao        int64                   `bson:"sqprestacao,omitempty" json:"sq_prestacao,omitempty"`
	DescricaoDeVinculo string                  `bson:"descricaodevinculo,omitempty" json:"descricao_de_vinculo,omitempty"`
	Partido            *Partido                `bson:"partido,omitempty" json:"partido,omitempty"`
}

type ReceitaCandidatoDetalhada struct {
	Receita            *ReceitaCandidato `bson:"receita" json:"receita"`
	SQCandidato        int64             `bson:"sqcandidato,omitempty" json:"sq_candidato,omitempty"`
	NumeroCandidato    *int              `bson:"numerocandidato,omitempty" json:"numero_candidato,omitempty"`
	NomeCandidato      string            `bson:"nomecandidato,omitempty" json:"nome_candidato,omitempty"`
	NomeUrnaCandidato  string            `bson:"nomeurnacandidato,omitempty" json:"nome_urna_candidato,omitempty"`
	CargoCandidato     string            `bson:"cargocandidato,omitempty" json:"cargo_candidato,omitempty"`
	UFCandidato        string            `bson:"ufcandidato,omitempty" json:"uf_candidato,omitempty"`
	PartidoSigla       string            `bson:"partidosigla,omitempty" json:"partido_sigla,omitempty"`
	PartidoNome        string            `bson:"partidonome,omitempty" json:"partido_nome,omitempty"`
	DescricaoDeVinculo string            `bson:"descricaodevinculo,omitempty" json:"descricao_de_vinculo,omitempty"`
	Candidato          *Candidato        `bson:"candidato,omitempty" json:"candidato,omitempty"`
}

type ReceitaOrgaoPartidarioDetalhada struct {
	Receita            *ReceitaOrgaoPartidario `bson:"receita" json:"receita"`
	PartidoNumero      int16                   `bson:"partidonumero,omitempty" json:"partido_numero,omitempty"`
	PartidoNome        string                  `bson:"partidonome,omitempty" json:"partido_nome,omitempty"`
	DescricaoDeVinculo string                  `bson:"descricaodevinculo,omitempty" json:"descricao_de_vinculo,omitempty"`
	Partido            *Partido                `bson:"partido,omitempty" json:"partido,omitempty"`
}
