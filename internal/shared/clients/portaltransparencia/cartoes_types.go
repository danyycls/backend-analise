package portaltransparencia

type CartaoQueryParams struct {
	Pagina              int
	MesExtratoInicio    string
	MesExtratoFim       string
	DataTransacaoInicio string
	DataTransacaoFim    string
	TipoCartao          string
	CodigoOrgao         string
	CPFPortador         string
	CPFCNPJFavorecido   string
	ValorDe             string
	ValorAte            string
}
