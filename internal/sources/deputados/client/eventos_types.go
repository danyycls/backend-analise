package deputados

type LocalCamara struct {
	Nome   string `json:"nome"`
	Predio string `json:"predio"`
	Sala   string `json:"sala"`
	Andar  string `json:"andar"`
}

type OrgaoEvento struct {
	ID             int    `json:"id"`
	URI            string `json:"uri"`
	Sigla          string `json:"sigla"`
	Nome           string `json:"nome"`
	Apelido        string `json:"apelido"`
	CodTipoOrgao   int    `json:"codTipoOrgao"`
	TipoOrgao      string `json:"tipoOrgao"`
	NomePublicacao string `json:"nomePublicacao"`
	NomeResumido   string `json:"nomeResumido"`
}

type Evento struct {
	ID             int           `json:"id"`
	URI            string        `json:"uri"`
	DataHoraInicio string        `json:"dataHoraInicio"`
	DataHoraFim    string        `json:"dataHoraFim"`
	Situacao       string        `json:"situacao"`
	DescricaoTipo  string        `json:"descricaoTipo"`
	Descricao      string        `json:"descricao"`
	LocalExterno   string        `json:"localExterno"`
	Orgaos         []OrgaoEvento `json:"orgaos"`
	LocalCamara    LocalCamara   `json:"localCamara"`
	URLRegistro    string        `json:"urlRegistro"`
}

type RequerimentoEvento struct {
	Titulo string `json:"titulo"`
	URI    string `json:"uri"`
}

type FaseEvento struct {
	DataHora  string `json:"dataHora"`
	Descricao string `json:"descricao"`
}

type EventoDetalhe struct {
	ID                int                  `json:"id"`
	URI               string               `json:"uri"`
	DataHoraInicio    string               `json:"dataHoraInicio"`
	DataHoraFim       string               `json:"dataHoraFim"`
	Situacao          string               `json:"situacao"`
	DescricaoTipo     string               `json:"descricaoTipo"`
	Descricao         string               `json:"descricao"`
	LocalExterno      string               `json:"localExterno"`
	Orgaos            []OrgaoEvento        `json:"orgaos"`
	LocalCamara       LocalCamara          `json:"localCamara"`
	URLRegistro       string               `json:"urlRegistro"`
	URLDocumentoPauta string               `json:"urlDocumentoPauta"`
	Requerimentos     []RequerimentoEvento `json:"requerimentos"`
	Fases             []FaseEvento         `json:"fases,omitempty"`
}
