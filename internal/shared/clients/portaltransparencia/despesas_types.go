package portaltransparencia

type DespesaRecursosRecebidosQueryParams struct {
	Pagina           int
	MesAnoInicio     string
	MesAnoFim        string
	NomeFavorecido   string
	CodigoFavorecido string
	TipoFavorecido   string
	UF               string
	CodigoIBGE       string
	OrgaoSuperior    string
	Orgao            string
	UnidadeGestora   string
}

type DespesaPorOrgaoQueryParams struct {
	Pagina        int
	Ano           string
	OrgaoSuperior string
	Orgao         string
}

type DespesaFuncionalProgramaticaQueryParams struct {
	Pagina    int
	Ano       string
	Funcao    string
	Subfuncao string
	Programa  string
	Acao      string
}

type DespesaMovimentacaoLiquidaQueryParams struct {
	Pagina              int
	Ano                 string
	Funcao              string
	Subfuncao           string
	Programa            string
	Acao                string
	GrupoDespesa        string
	ElementoDespesa     string
	ModalidadeAplicacao string
	IDPlanoOrcamentario string
}

type DespesaPlanoOrcamentarioQueryParams struct {
	Pagina                    int
	Ano                       string
	CodPlanoOrcamentario      string
	DescPlanoOrcamentario     string
	CodPOIdentfAcompanhamento string
}

type ListarFuncionalProgramaticaQueryParams struct {
	AnoInicio int
	Pagina    int
	Codigo    string
}

type DespesaDocumentosQueryParams struct {
	DataEmissao    string
	Fase           string
	Pagina         int
	UnidadeGestora string
	Gestao         string
}

type DespesaDocumentosPorFavorecidoQueryParams struct {
	CodigoPessoa       string
	Fase               string
	Ano                string
	Pagina             int
	UG                 string
	Gestao             string
	OrdenacaoResultado string
}
