package pncp

type AnalisePublicacaoRequest struct {
	Tipo                        string `json:"tipo"`
	UF                          string `json:"uf"`
	CodigoMunicipioIbge         string `json:"codigo_municipio_ibge"`
	DataInicial                 string `json:"data_inicial"`
	DataFinal                   string `json:"data_final"`
	CodigoModalidadeContratacao string `json:"codigo_modalidade_contratacao,omitempty"`
}

type PublicacaoSearchResponse struct {
	Data             []map[string]interface{} `json:"data"`
	TotalRegistros   int                      `json:"totalRegistros"`
	TotalPaginas     int                      `json:"totalPaginas"`
	NumeroPagina     int                      `json:"numeroPagina"`
	PaginasRestantes int                      `json:"paginasRestantes"`
	Empty            bool                     `json:"empty"`
}

type AmparoLegal struct {
	Descricao *string `json:"descricao"`
	Nome      *string `json:"nome"`
	Codigo    *int    `json:"codigo"`
}

type FonteOrcamentaria struct {
	Codigo       *int    `json:"codigo"`
	Nome         *string `json:"nome"`
	Descricao    *string `json:"descricao"`
	DataInclusao *string `json:"dataInclusao"`
}
