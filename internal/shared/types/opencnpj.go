package types

type FornecedorOpenCNPJ struct {
	CapitalSocial     *string                `json:"capitalSocial"`
	CNPJ              *string                `json:"cnpj"`
	NomeFantasia      *string                `json:"nomeFantasia"`
	Socios            []Socio                `json:"qsa"`
	RazaoSocial       *string                `json:"razaoSocial"`
	SituacaoCadastral *string                `json:"situacaoCadastral"`
	Extra             map[string]interface{} `json:"-"`
}

type Socio struct {
	CNPJCPFSocio              *string                `json:"cnpj_cpf_socio"`
	CodigoPais                *string                `json:"codigo_pais"`
	DataEntradaSociedade      *string                `json:"data_entrada_sociedade"`
	FaixaEtaria               *string                `json:"faixa_etaria"`
	IdentificadorSocio        *string                `json:"identificador_socio"`
	NomeRepresentante         *string                `json:"nome_representante"`
	NomeSocio                 *string                `json:"nome_socio"`
	Pais                      *PaisInfo              `json:"pais"`
	QualificacaoRepresentante *QualificacaoInfo      `json:"qualificacao_representante"`
	QualificacaoSocio         *string                `json:"qualificacao_socio"`
	RepresentanteLegal        *string                `json:"representante_legal"`
	Extra                     map[string]interface{} `json:"-"`
}

type PaisInfo struct {
	Codigo    *string                `json:"codigo"`
	Descricao *string                `json:"descricao"`
	Extra     map[string]interface{} `json:"-"`
}

type QualificacaoInfo struct {
	Codigo    *string                `json:"codigo"`
	Descricao *string                `json:"descricao"`
	Extra     map[string]interface{} `json:"-"`
}

type OpenCNPJResponse struct {
	CNPJ              string  `json:"cnpj"`
	RazaoSocial       string  `json:"razao_social"`
	NomeFantasia      string  `json:"nome_fantasia"`
	SituacaoCadastral string  `json:"situacao_cadastral"`
	CapitalSocial     string  `json:"capital_social"`
	Socios            []Socio `json:"qsa"`
}
