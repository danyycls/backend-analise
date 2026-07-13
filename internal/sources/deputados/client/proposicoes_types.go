package deputados

type Proposicao struct {
	ID               int    `json:"id"`
	URI              string `json:"uri"`
	SiglaTipo        string `json:"siglaTipo"`
	CodTipo          int    `json:"codTipo"`
	Numero           int    `json:"numero"`
	Ano              int    `json:"ano"`
	Ementa           string `json:"ementa"`
	DataApresentacao string `json:"dataApresentacao"`
}

type StatusProposicao struct {
	DataHora            string `json:"dataHora"`
	Sequencia           int    `json:"sequencia"`
	SiglaOrgao          string `json:"siglaOrgao"`
	URIOrgao            string `json:"uriOrgao"`
	URIUltimoRelator    string `json:"uriUltimoRelator"`
	Regime              string `json:"regime"`
	DescricaoTramitacao string `json:"descricaoTramitacao"`
	CodTipoTramitacao   int    `json:"codTipoTramitacao"`
	DescricaoSituacao   string `json:"descricaoSituacao"`
	CodSituacao         int    `json:"codSituacao"`
	Despacho            string `json:"despacho"`
	URL                 string `json:"url"`
	Ambito              string `json:"ambito"`
	Apreciacao          string `json:"apreciacao"`
}

type ProposicaoDetalhe struct {
	ID                int              `json:"id"`
	URI               string           `json:"uri"`
	SiglaTipo         string           `json:"siglaTipo"`
	CodTipo           int              `json:"codTipo"`
	Numero            int              `json:"numero"`
	Ano               int              `json:"ano"`
	Ementa            string           `json:"ementa"`
	DataApresentacao  string           `json:"dataApresentacao"`
	DescricaoTipo     string           `json:"descricaoTipo"`
	EmentaDetalhada   string           `json:"ementaDetalhada"`
	Keywords          string           `json:"keywords"`
	URIOrgaoNumerador string           `json:"uriOrgaoNumerador"`
	StatusProposicao  StatusProposicao `json:"statusProposicao"`
	URIAutores        string           `json:"uriAutores"`
	URIPropPrincipal  string           `json:"uriPropPrincipal"`
	URIPropAnterior   string           `json:"uriPropAnterior"`
	URIPropPosterior  string           `json:"uriPropPosterior"`
	URLInteiroTeor    string           `json:"urlInteiroTeor"`
	URNFinal          string           `json:"urnFinal"`
	Texto             string           `json:"texto"`
	Justificativa     string           `json:"justificativa"`
}

type Tramitacao struct {
	DataHora            string `json:"dataHora"`
	Sequencia           int    `json:"sequencia"`
	SiglaOrgao          string `json:"siglaOrgao"`
	URIOrgao            string `json:"uriOrgao"`
	URIUltimoRelator    string `json:"uriUltimoRelator"`
	Regime              string `json:"regime"`
	DescricaoTramitacao string `json:"descricaoTramitacao"`
	CodTipoTramitacao   int    `json:"codTipoTramitacao"`
	DescricaoSituacao   string `json:"descricaoSituacao"`
	CodSituacao         int    `json:"codSituacao"`
	Despacho            string `json:"despacho"`
	URL                 string `json:"url"`
	Ambito              string `json:"ambito"`
	Apreciacao          string `json:"apreciacao"`
}

type Author struct {
	URI             string `json:"uri"`
	Nome            string `json:"nome"`
	CodTipo         int    `json:"codTipo"`
	Tipo            string `json:"tipo"`
	OrdemAssinatura int    `json:"ordemAssinatura"`
	Proponente      int    `json:"proponente"`
}

type Tema struct {
	CodTema    int    `json:"codTema"`
	Tema       string `json:"tema"`
	Relevancia int    `json:"relevancia"`
}
