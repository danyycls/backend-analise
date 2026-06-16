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
	CPFCNPJ                    string              `json:"cpf_cnpj"`
	Nome                       string              `json:"nome"`
	NomeRFB                    string              `json:"nome_rfb,omitempty"`
	TipoFornecedorCodigo       *int                `json:"tipo_fornecedor_codigo,omitempty"`
	TipoFornecedorDescricao    string              `json:"tipo_fornecedor_descricao,omitempty"`
	CNAECodigo                 string              `json:"cnae_codigo,omitempty"`
	CNAEDescricao              string              `json:"cnae_descricao,omitempty"`
	EsferaPartidariaCodigo     string              `json:"esfera_partidaria_codigo,omitempty"`
	EsferaPartidariaDescricao  string              `json:"esfera_partidaria_descricao,omitempty"`
	UFSigla                    *string             `json:"sg_uf,omitempty"`
	MunicipioNome              string              `json:"municipio_nome,omitempty"`
	SQCandidatoRelacionado     *int64              `json:"sq_candidato_relacionado,omitempty"`
	NumeroCandidatoRelacionado *int                `json:"numero_candidato_relacionado,omitempty"`
	CargoCodigoRelacionado     *int                `json:"cargo_codigo_relacionado,omitempty"`
	CargoDescricaoRelacionada  string              `json:"cargo_descricao_relacionada,omitempty"`
	PartidoNumeroRelacionado   *int16              `json:"partido_numero_relacionado,omitempty"`
	PartidoSiglaRelacionado    string              `json:"partido_sigla_relacionado,omitempty"`
	PartidoNomeRelacionado     string              `json:"partido_nome_relacionado,omitempty"`
	Enriquecimento             *FornecedorOpenCNPJ `json:"enriquecimento,omitempty"`
}

type FornecedorDetalhado struct {
	Fornecedor              FornecedorEnriquecido             `json:"fornecedor"`
	DespesasCandidato       []DespesaCandidatoDetalhada       `json:"despesas_candidato,omitempty"`
	DespesasOrgaoPartidario []DespesaOrgaoPartidarioDetalhada `json:"despesas_orgao_partidario,omitempty"`
	DescricaoDeVinculo      string                            `json:"descricao_de_vinculo,omitempty"`
}

type DespesaCandidatoDetalhada struct {
	Despesa            *DespesaCandidato `json:"despesa"`
	SQCandidato        int64             `json:"sq_candidato,omitempty"`
	SQPrestacao        int64             `json:"sq_prestacao,omitempty"`
	DescricaoDeVinculo string            `json:"descricao_de_vinculo,omitempty"`
}

type DespesaOrgaoPartidarioDetalhada struct {
	Despesa            *DespesaOrgaoPartidario `json:"despesa"`
	PartidoNumero      int16                   `json:"partido_numero,omitempty"`
	PartidoNome        string                  `json:"partido_nome,omitempty"`
	SQPrestacao        int64                   `json:"sq_prestacao,omitempty"`
	DescricaoDeVinculo string                  `json:"descricao_de_vinculo,omitempty"`
}
