package tsetypes

type BuscaRelacoesRequest struct {
	Documento string `json:"documento" binding:"required"`
}

type RelacoesResponse struct {
	Fornecedor    *FornecedorEnriquecido `json:"fornecedor,omitempty"`
	Doador        *DoadorRelacoes        `json:"doador,omitempty"`
	Despesas      []DespesaRelacao       `json:"despesas,omitempty"`
	Receitas      []ReceitaRelacao       `json:"receitas,omitempty"`
	TotalDespesas int                    `json:"total_despesas"`
	TotalReceitas int                    `json:"total_receitas"`
}

type DoadorRelacoes struct {
	CPFCNPJ string `json:"cpf_cnpj"`
	Nome    string `json:"nome"`
}

type CandidatoResumido struct {
	SQCandidato  int64  `json:"sq_candidato"`
	NomeCompleto string `json:"nome_completo"`
	PartidoSigla string `json:"partido_sigla,omitempty"`
	PartidoNome  string `json:"partido_nome,omitempty"`
	CargoNome    string `json:"cargo_nome"`
	UFSigla      string `json:"sg_uf"`
}

type PartidoResumido struct {
	Numero int16  `json:"numero"`
	Sigla  string `json:"sigla"`
	Nome   string `json:"nome"`
}

type DespesaRelacao struct {
	SQDespesa              int64              `json:"sq_despesa"`
	Tipo                   string             `json:"tipo"`
	TipoRegistro           string             `json:"tipo_registro,omitempty"`
	DataDespesa            *string            `json:"data_despesa,omitempty"`
	Descricao              string             `json:"descricao"`
	Valor                  float64            `json:"valor"`
	OrigemDespesaDescricao string             `json:"origem_despesa_descricao,omitempty"`
	Candidato              *CandidatoResumido `json:"candidato,omitempty"`
	Partido                *PartidoResumido   `json:"partido,omitempty"`
}

type ReceitaRelacao struct {
	SQReceita              int64              `json:"sq_receita"`
	Tipo                   string             `json:"tipo"`
	DataReceita            *string            `json:"data_receita,omitempty"`
	Descricao              string             `json:"descricao"`
	Valor                  float64            `json:"valor"`
	OrigemReceitaDescricao string             `json:"origem_receita_descricao,omitempty"`
	Candidato              *CandidatoResumido `json:"candidato,omitempty"`
	Partido                *PartidoResumido   `json:"partido,omitempty"`
}

type DoadorRelacoesResponse struct {
	Doador        *DoadorRelacoes  `json:"doador,omitempty"`
	Receitas      []ReceitaRelacao `json:"receitas,omitempty"`
	TotalReceitas int              `json:"total_receitas"`
}

type FornecedorRelacoesResponse struct {
	Fornecedor    *FornecedorEnriquecido `json:"fornecedor,omitempty"`
	Despesas      []DespesaRelacao       `json:"despesas,omitempty"`
	TotalDespesas int                    `json:"total_despesas"`
}
