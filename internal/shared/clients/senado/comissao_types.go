package senado

type ComissaoResumo struct {
	Codigo string `json:"Codigo"`
	Sigla  string `json:"Sigla"`
	Nome   string `json:"Nome"`
}

type ColegiadoListaWrapper struct {
	ListaColegiados struct {
		Colegiados struct {
			Colegiado []ComissaoResumo `json:"Colegiado"`
		} `json:"Colegiados"`
	} `json:"ListaColegiados"`
}

type MembroComissao struct {
	CodigoParlamentar       string `json:"CodigoParlamentar"`
	NomeParlamentar         string `json:"NomeParlamentar"`
	SiglaPartidoParlamentar string `json:"SiglaPartidoParlamentar"`
	UfParlamentar           string `json:"UfParlamentar"`
	DescricaoParticipacao   string `json:"DescricaoParticipacao"`
	DataInicio              string `json:"DataInicio"`
	DataFim                 string `json:"DataFim"`
}

type ComissaoDetalheColegiadoRaw struct {
	CodigoColegiado string `json:"CodigoColegiado"`
	SiglaColegiado  string `json:"SiglaColegiado"`
	NomeColegiado   string `json:"NomeColegiado"`
	MembrosBlocoSF  *struct {
		PartidoBloco []struct {
			MembrosSF *struct {
				Membro []MembroBlocoRaw `json:"Membro"`
			} `json:"MembrosSF"`
		} `json:"PartidoBloco"`
	} `json:"MembrosBlocoSF"`
}

type MembroBlocoRaw struct {
	CodigoParlamentar string `json:"CodigoParlamentar"`
	NomeParlamentar   string `json:"NomeParlamentar"`
	Partido           string `json:"Partido"`
	SiglaUf           string `json:"SiglaUf"`
	TipoVaga          string `json:"TipoVaga"`
}

type ComissaoDetalheWrapper struct {
	ComissoesCongressoNacional struct {
		Colegiados struct {
			Colegiado []ComissaoDetalheColegiadoRaw `json:"Colegiado"`
		} `json:"Colegiados"`
	} `json:"ComissoesCongressoNacional"`
}

type ComissaoDetalhe struct {
	Codigo  string           `json:"codigo"`
	Sigla   string           `json:"sigla"`
	Nome    string           `json:"nome"`
	Membros []MembroComissao `json:"membros"`
}
