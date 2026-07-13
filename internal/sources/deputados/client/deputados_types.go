package deputados

type Deputado struct {
	ID            int    `json:"id"`
	Nome          string `json:"nome"`
	SiglaPartido  string `json:"siglaPartido"`
	SiglaUF       string `json:"siglaUf"`
	URI           string `json:"uri"`
	URLFoto       string `json:"urlFoto"`
	Email         string `json:"email"`
	IDLegislatura int    `json:"idLegislatura"`
	URIPartido    string `json:"uriPartido"`
}

type Gabinete struct {
	Nome     string `json:"nome"`
	Predio   string `json:"predio"`
	Sala     string `json:"sala"`
	Andar    string `json:"andar"`
	Telefone string `json:"telefone"`
	Email    string `json:"email"`
}

type UltimoStatus struct {
	ID                int      `json:"id"`
	Nome              string   `json:"nome"`
	SiglaPartido      string   `json:"siglaPartido"`
	SiglaUF           string   `json:"siglaUf"`
	URI               string   `json:"uri"`
	URIPartido        string   `json:"uriPartido"`
	URLFoto           string   `json:"urlFoto"`
	Email             string   `json:"email"`
	NomeEleitoral     string   `json:"nomeEleitoral"`
	Situacao          string   `json:"situacao"`
	CondicaoEleitoral string   `json:"condicaoEleitoral"`
	DescricaoStatus   string   `json:"descricaoStatus"`
	Data              string   `json:"data"`
	IDLegislatura     int      `json:"idLegislatura"`
	Gabinete          Gabinete `json:"gabinete"`
}

type DeputadoDetalhe struct {
	ID                  int          `json:"id"`
	URI                 string       `json:"uri"`
	CPF                 string       `json:"cpf"`
	NomeCivil           string       `json:"nomeCivil"`
	DataNascimento      string       `json:"dataNascimento"`
	DataFalecimento     string       `json:"dataFalecimento,omitempty"`
	Sexo                string       `json:"sexo"`
	UFNascimento        string       `json:"ufNascimento"`
	MunicipioNascimento string       `json:"municipioNascimento"`
	Escolaridade        string       `json:"escolaridade"`
	UltimoStatus        UltimoStatus `json:"ultimoStatus"`
	RedeSocial          []string     `json:"redeSocial"`
	URLWebsite          string       `json:"urlWebsite"`
}

type DeputadoDespesa struct {
	Ano               int     `json:"ano"`
	Mes               int     `json:"mes"`
	TipoDocumento     string  `json:"tipoDocumento"`
	CodTipoDocumento  int     `json:"codTipoDocumento"`
	TipoDespesa       string  `json:"tipoDespesa"`
	CodDocumento      string  `json:"codDocumento"`
	NumDocumento      string  `json:"numDocumento"`
	CodLote           int     `json:"codLote"`
	Parcela           int     `json:"parcela"`
	ValorDocumento    float64 `json:"valorDocumento"`
	ValorGlosa        float64 `json:"valorGlosa"`
	ValorLiquido      float64 `json:"valorLiquido"`
	NumRessarcimento  string  `json:"numRessarcimento"`
	DataDocumento     string  `json:"dataDocumento"`
	CNPJCPFFornecedor string  `json:"cnpjCpfFornecedor"`
	NomeFornecedor    string  `json:"nomeFornecedor"`
	URLDocumento      string  `json:"urlDocumento"`
}

type DeputadoHistorico struct {
	ID                int    `json:"id"`
	IDLegislatura     int    `json:"idLegislatura"`
	URI               string `json:"uri"`
	Nome              string `json:"nome"`
	NomeEleitoral     string `json:"nomeEleitoral"`
	SiglaPartido      string `json:"siglaPartido"`
	URIPartido        string `json:"uriPartido"`
	SiglaUF           string `json:"siglaUf"`
	URLFoto           string `json:"urlFoto"`
	Email             string `json:"email"`
	Situacao          string `json:"situacao"`
	CondicaoEleitoral string `json:"condicaoEleitoral"`
	DescricaoStatus   string `json:"descricaoStatus"`
	DataHora          string `json:"dataHora"`
}

type DeputadoMandatoExterno struct {
	AnoInicio           string `json:"anoInicio"`
	AnoFim              string `json:"anoFim"`
	Cargo               string `json:"cargo"`
	SiglaPartidoEleicao string `json:"siglaPartidoEleicao"`
	URIPartidoEleicao   string `json:"uriPartidoEleicao"`
	Municipio           string `json:"municipio"`
	SiglaUF             string `json:"siglaUf"`
}

type DeputadoOrgao struct {
	IDOrgao        int    `json:"idOrgao"`
	SiglaOrgao     string `json:"siglaOrgao"`
	NomeOrgao      string `json:"nomeOrgao"`
	NomePublicacao string `json:"nomePublicacao"`
	URIOrgao       string `json:"uriOrgao"`
	Titulo         string `json:"titulo"`
	CodTitulo      string `json:"codTitulo"`
	DataInicio     string `json:"dataInicio"`
	DataFim        string `json:"dataFim"`
}

type DeputadoResponse struct {
	Deputado         *DeputadoDetalhe         `json:"deputado"`
	Frentes          []Frente                 `json:"frentes"`
	Historico        []DeputadoHistorico      `json:"historico"`
	MandatosExternos []DeputadoMandatoExterno `json:"mandatosExternos"`
}
