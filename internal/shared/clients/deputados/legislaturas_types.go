package deputados

type Legislatura struct {
	ID         int    `json:"id"`
	URI        string `json:"uri"`
	DataInicio string `json:"dataInicio"`
	DataFim    string `json:"dataFim"`
}

type BancadaLideres struct {
	Nome string `json:"nome"`
	Tipo string `json:"tipo"`
	URI  string `json:"uri"`
}

type ParlamentarLider struct {
	ID            int    `json:"id"`
	URI           string `json:"uri"`
	Nome          string `json:"nome"`
	SiglaPartido  string `json:"siglaPartido"`
	URIPartido    string `json:"uriPartido"`
	SiglaUF       string `json:"siglaUf"`
	URLFoto       string `json:"urlFoto"`
	Email         string `json:"email"`
	IDLegislatura int    `json:"idLegislatura"`
}

type Lider struct {
	Parlamentar ParlamentarLider `json:"parlamentar"`
	Bancada     BancadaLideres   `json:"bancada"`
	Titulo      string           `json:"titulo"`
	DataInicio  string           `json:"dataInicio"`
	DataFim     string           `json:"dataFim"`
}

type MembroMesa struct {
	ID            int    `json:"id"`
	URI           string `json:"uri"`
	Nome          string `json:"nome"`
	SiglaPartido  string `json:"siglaPartido"`
	URIPartido    string `json:"uriPartido"`
	SiglaUF       string `json:"siglaUf"`
	URLFoto       string `json:"urlFoto"`
	Email         string `json:"email"`
	Titulo        string `json:"titulo"`
	CodTitulo     string `json:"codTitulo"`
	IDLegislatura int    `json:"idLegislatura"`
	DataInicio    string `json:"dataInicio"`
	DataFim       string `json:"dataFim"`
}
