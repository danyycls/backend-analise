package portaltransparencia

type PessoaFisicaQueryParams struct {
	Pagina                    int
	CPF                       string
	Nome                      string
	NIS                       string
	FavorecidoDespesas        string
	Servidor                  string
	BeneficiarioDiarias       string
	Permissionario            string
	Contratado                string
	SancionadoCEIS            string
	SancionadoCNEP            string
	SancionadoCEPIM           string
	SancionadoCEAF            string
	SancionadoAcordoLeniencia string
	Ordenacao                 string
	OrdenacaoDirecao          string
}

type PessoaJuridicaQueryParams struct {
	Pagina                    int
	CNPJ                      string
	RazaoSocial               string
	NomeFantasia              string
	FavorecidoDespesas        string
	PossuiContratacao         string
	Convenios                 string
	FavorecidoTransferencias  string
	SancionadoCEPIM           string
	SancionadoCEIS            string
	SancionadoCNEP            string
	SancionadoCEAF            string
	SancionadoAcordoLeniencia string
	Ordenacao                 string
	OrdenacaoDirecao          string
}
