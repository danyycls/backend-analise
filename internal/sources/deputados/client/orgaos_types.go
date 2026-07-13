package deputados

type Orgao struct {
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

type OrgaoDetalhe struct {
	ID              int    `json:"id"`
	URI             string `json:"uri"`
	Sigla           string `json:"sigla"`
	Nome            string `json:"nome"`
	Apelido         string `json:"apelido"`
	CodTipoOrgao    int    `json:"codTipoOrgao"`
	TipoOrgao       string `json:"tipoOrgao"`
	NomePublicacao  string `json:"nomePublicacao"`
	NomeResumido    string `json:"nomeResumido"`
	DataInicio      string `json:"dataInicio"`
	DataInstalacao  string `json:"dataInstalacao"`
	DataFim         string `json:"dataFim"`
	DataFimOriginal string `json:"dataFimOriginal"`
	Casa            string `json:"casa"`
	Sala            string `json:"sala"`
	URLWebsite      string `json:"urlWebsite"`
}

type MembroOrgao struct {
	ID            int    `json:"id"`
	URI           string `json:"uri"`
	Nome          string `json:"nome"`
	SiglaPartido  string `json:"siglaPartido"`
	URIPartido    string `json:"uriPartido"`
	SiglaUF       string `json:"siglaUf"`
	IDLegislatura int    `json:"idLegislatura"`
	URLFoto       string `json:"urlFoto"`
	Email         string `json:"email"`
	Titulo        string `json:"titulo"`
	CodTitulo     string `json:"codTitulo"`
	DataInicio    string `json:"dataInicio"`
	DataFim       string `json:"dataFim"`
}
