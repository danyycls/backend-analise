package portaltransparencia

type ServidorQueryParams struct {
	Pagina                 int
	CPF                    string
	Nome                   string
	OrgaoServidorLotacao   string
	OrgaoServidorExercicio string
	SituacaoServidor       string
	TipoServidor           string
	CodigoFuncaoCargo      string
}

type ServidorRemuneracaoQueryParams struct {
	Pagina                          int
	CPF                             string
	IDservidorAposentadoPensionista string
	MesAno                          string
}

type ServidorPorOrgaoQueryParams struct {
	Pagina         int
	OrgaoLotacao   string
	OrgaoExercicio string
	TipoServidor   string
	TipoVinculo    string
	Licenca        string
}

type FuncaoCargoQueryParams struct {
	Pagina               int
	CodigoFuncaoCargo    string
	DescricaoFuncaoCargo string
}

type PEPQueryParams struct {
	Pagina                 int
	CPF                    string
	Nome                   string
	DescricaoFuncao        string
	OrgaoServidorLotacao   string
	DataInicioExercicioDe  string
	DataInicioExercicioAte string
	DataFimExercicioDe     string
	DataFimExercicioAte    string
}
