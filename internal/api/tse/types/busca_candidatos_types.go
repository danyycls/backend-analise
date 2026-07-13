package tsetypes

type BuscaFiltrosParams struct {
	Tag string `json:"tag"`
}

type OpcaoFiltro struct {
	Valor string `json:"valor"`
	Label string `json:"label"`
}

type OpcoesFiltroResponse struct {
	Opcoes []OpcaoFiltro `json:"opcoes"`
}

type BuscaCandidatosRequest struct {
	CargoNome *string `json:"cargo_nome,omitempty"`
	PartidoID *string `json:"partido_id,omitempty"`
	Eleito    *string `json:"eleito,omitempty"`
	UFSigla   *string `json:"sg_uf,omitempty"`
}

type CandidatoLista struct {
	SQCandidato                  int64            `json:"sq_candidato"`
	NomeCompleto                 string           `json:"nome_completo"`
	NomeUrna                     string           `json:"nome_urna,omitempty"`
	CPF                          string           `json:"cpf"`
	NumeroCandidato              *int             `json:"numero_candidato,omitempty"`
	CargoCodigo                  *int             `json:"cargo_codigo,omitempty"`
	CargoNome                    string           `json:"cargo_nome"`
	UFSigla                      string           `json:"sg_uf"`
	Partido                      *PartidoResumido `json:"partido,omitempty"`
	SituacaoTotalizacaoDescricao string           `json:"situacao_totalizacao_descricao,omitempty"`
	Eleito                       bool             `json:"eleito"`
}

type CandidatosResponse struct {
	Candidatos []CandidatoLista `json:"candidatos"`
	Total      int              `json:"total"`
}

type BuscaDocumentoRequest struct {
	Documento string `json:"documento"`
}
