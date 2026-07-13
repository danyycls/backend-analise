package deputados

type Partido struct {
	ID    int    `json:"id"`
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
	URI   string `json:"uri"`
}

type PartidoStatus struct {
	Data          string `json:"data"`
	IDLegislatura string `json:"idLegislatura"`
	Situacao      string `json:"situacao"`
	TotalPosse    int    `json:"totalPosse"`
	TotalMembros  int    `json:"totalMembros"`
	URIMembros    string `json:"uriMembros"`
	Lider         struct {
		URI          string `json:"uri"`
		Nome         string `json:"nome"`
		SiglaPartido string `json:"siglaPartido"`
		UF           string `json:"uf"`
		URLFoto      string `json:"urlFoto"`
	} `json:"lider"`
}

type PartidoDetalhe struct {
	ID              int           `json:"id"`
	Sigla           string        `json:"sigla"`
	Nome            string        `json:"nome"`
	URI             string        `json:"uri"`
	Status          PartidoStatus `json:"status"`
	NumeroEleitoral string        `json:"numeroEleitoral"`
	URLLogo         string        `json:"urlLogo"`
	URLWebSite      string        `json:"urlWebSite"`
	URLFacebook     string        `json:"urlFacebook"`
}
