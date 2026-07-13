package deputados

type Frente struct {
	ID            int    `json:"id"`
	Titulo        string `json:"titulo"`
	URI           string `json:"uri"`
	IDLegislatura int    `json:"idLegislatura"`
}

type FrenteDetalhe struct {
	ID            int      `json:"id"`
	Titulo        string   `json:"titulo"`
	URI           string   `json:"uri"`
	IDLegislatura int      `json:"idLegislatura"`
	IDSituacao    int      `json:"idSituacao"`
	Situacao      string   `json:"situacao"`
	Keywords      string   `json:"keywords"`
	Email         string   `json:"email"`
	Telefone      string   `json:"telefone"`
	URLWebsite    string   `json:"urlWebsite"`
	URLDocumento  string   `json:"urlDocumento"`
	Coordenador   Deputado `json:"coordenador"`
}

type MembroFrente struct {
	ID            int    `json:"id"`
	URI           string `json:"uri"`
	Nome          string `json:"nome"`
	SiglaPartido  string `json:"siglaPartido"`
	URIPartido    string `json:"uriPartido"`
	SiglaUF       string `json:"siglaUf"`
	URLFoto       string `json:"urlFoto"`
	Email         string `json:"email"`
	Titulo        string `json:"titulo"`
	CodTitulo     int    `json:"codTitulo"`
	IDLegislatura int    `json:"idLegislatura"`
	DataInicio    string `json:"dataInicio"`
	DataFim       string `json:"dataFim"`
}
