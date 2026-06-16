package tcu

type TCUQueryParams struct {
	ParteNome string `json:"parteNome,omitempty"`
	CPF       string `json:"cpf,omitempty"`
	CNPJ      string `json:"cnpj,omitempty"`
	UF        string `json:"uf,omitempty"`
	Municipio string `json:"municipio,omitempty"`
}

type ContasIrregulares struct {
	NumeroProcessoFormatado    string `json:"numeroProcessoFormatado"`
	Nome                       string `json:"nome"`
	TipoRegistro               string `json:"tipoRegistro"`
	NumeroRegistro             string `json:"numeroRegistro"`
	Municipio                  string `json:"municipio"`
	UF                         string `json:"uf"`
	DataTransitoEmJulgado      string `json:"dataTransitoEmJulgado"`
	LinkDeliberacoesProcesso   string `json:"linkDeliberacoesProcesso"`
	LinkAcompanhamentoProcesso string `json:"linkAcompanhamentoProcesso"`
}

type FinsEleitorais struct {
	NumeroProcessoFormatado    string `json:"numeroProcessoFormatado"`
	Nome                       string `json:"nome"`
	NumeroRegistro             string `json:"numeroRegistro"`
	Municipio                  string `json:"municipio"`
	UF                         string `json:"uf"`
	DataTransitoEmJulgado      string `json:"dataTransitoEmJulgado"`
	DataFinalFinsEleitorais    string `json:"dataFinalFinsEleitorais"`
	LinkDeliberacoesProcesso   string `json:"linkDeliberacoesProcesso"`
	LinkAcompanhamentoProcesso string `json:"linkAcompanhamentoProcesso"`
}

type Sancoes struct {
	NumeroProcessoFormatado    string  `json:"numeroProcessoFormatado"`
	Nome                       string  `json:"nome"`
	TipoRegistro               string  `json:"tipoRegistro"`
	NumeroRegistro             string  `json:"numeroRegistro"`
	Municipio                  *string `json:"municipio"`
	UF                         *string `json:"uf"`
	NumeroAcordaoFormatado     string  `json:"numeroAcordaoFormatado"`
	DataAcordao                string  `json:"dataAcordao"`
	DataTransitoEmJulgado      string  `json:"dataTransitoEmJulgado"`
	DataFinalSancao            string  `json:"dataFinalSancao"`
	LinkDeliberacoesProcesso   string  `json:"linkDeliberacoesProcesso"`
	LinkAcompanhamentoProcesso string  `json:"linkAcompanhamentoProcesso"`
}
