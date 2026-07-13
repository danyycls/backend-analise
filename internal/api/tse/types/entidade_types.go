package tsetypes

type ConsultaEntidadeRequest struct {
	Tipo  string `json:"tipo"`
	Chave string `json:"chave"`
}

type CandidatoEntidade struct {
	SQCandidato  int64  `json:"sq_candidato"`
	NomeCompleto string `json:"nome_completo"`
	CPF          string `json:"cpf"`
	PartidoSigla string `json:"partido_sigla,omitempty"`
	PartidoNome  string `json:"partido_nome,omitempty"`
	CargoNome    string `json:"cargo_nome"`
	UFSigla      string `json:"sg_uf"`
	Genero       string `json:"genero,omitempty"`
	CorRaca      string `json:"cor_raca,omitempty"`
	OcupacaoNome string `json:"ocupacao_nome,omitempty"`
}

type FornecedorEntidade struct {
	CPFCNPJ string `json:"cpf_cnpj"`
	Nome    string `json:"nome"`
	NomeRFB string `json:"nome_rfb,omitempty"`
	CNAE    string `json:"cnae,omitempty"`
	UFSigla string `json:"sg_uf,omitempty"`
}

type DoadorEntidade struct {
	CPFCNPJ string `json:"cpf_cnpj"`
	Nome    string `json:"nome"`
	NomeRFB string `json:"nome_rfb,omitempty"`
	CNAE    string `json:"cnae,omitempty"`
	UFSigla string `json:"sg_uf,omitempty"`
}

type ReceitaEntidade struct {
	SQReceita   int64              `json:"sq_receita"`
	Tipo        string             `json:"tipo"`
	Descricao   string             `json:"descricao"`
	Valor       float64            `json:"valor"`
	DataReceita *string            `json:"data_receita,omitempty"`
	Candidato   *CandidatoResumido `json:"candidato,omitempty"`
	Partido     *PartidoResumido   `json:"partido,omitempty"`
	DoadorNome  string             `json:"doador_nome,omitempty"`
}

type DespesaEntidade struct {
	SQDespesa      int64              `json:"sq_despesa"`
	Tipo           string             `json:"tipo"`
	Descricao      string             `json:"descricao"`
	Valor          float64            `json:"valor"`
	DataDespesa    *string            `json:"data_despesa,omitempty"`
	Candidato      *CandidatoResumido `json:"candidato,omitempty"`
	Partido        *PartidoResumido   `json:"partido,omitempty"`
	FornecedorNome string             `json:"fornecedor_nome,omitempty"`
}

type ConsultaEntidadeResponse struct {
	Tipo  string      `json:"tipo"`
	Chave string      `json:"chave"`
	Dados interface{} `json:"dados,omitempty"`
	Erro  string      `json:"erro,omitempty"`
}
