package repositorios

import "time"

type ContratoPersistido struct {
	ID                   string
	NumeroControlePNCP   string
	CNPJOrgao            string
	UGUFSigla            string
	UGCodigoIbge         string
	DataPublicacaoPncp   *time.Time
	DataAssinatura       *time.Time
	DataInicioVigencia   *time.Time
	DataTerminoVigencia  *time.Time
	ValorGlobal          *float64
	ValorInicial         *float64
	ValorTotalEstimado   *float64
	ValorTotalHomologado *float64
	NIFornecedor         *string
	CodigoAmparoLegal    *int
	NumeroContrato       *string
	CodigoContrato       *string
	CodigoTipoContrato   *int
	TipoContratoNome     *string
	CodigoUG             *string
	NomeUG               *string
	UGMunicipioNome      *string
	UGUFNome             *string
	ModalidadeNome       *string
	CodigoOrgao          *string
	NomeOrgao            *string
	NomeOrgaoSub         *string
	ObjetoContrato       *string
	NumeroLicitacao      *string
	OrigemLicitacao      *string
	Produto              *string
	SubtipoContrato      *string
	AnoContrato          *int
	NomeRazaoSocialFornecedor *string
	DadosCompletos       []byte
}

type FornecedorPersistido struct {
	CNPJ          string
	RazaoSocial   string
	DadosCompletos []byte
}

type SocioPersistido struct {
	ID           string
	CNPJCPFSocio string
	NomeSocio    *string
}

type FornecedorSocioPersistido struct {
	CNPJFornecedor         string
	SocioID                string
	DataEntradaSociedade   *string
	IdentificadorSocio     *string
	NomeSocio              *string
	QualificacaoSocio      *string
	NomeRepresentante      *string
	QualificacaoRepresentante *string
	RepresentanteLegal     *string
	FaixaEtaria            *string
	PaisCodigo             *string
	PaisDescricao          *string
	CNPJCPFSocio           string
	NomeSocioGlobal        string
}

type AmparoLegalPersistido struct {
	Codigo    int
	Nome      string
	Descricao *string
}

type BuscaControlePersistido struct {
	TipoBusca                string
	ValorBusca               string
	Ano                      int
	Mes                      int
	DataInicial              time.Time
	DataFinal                time.Time
	TotalContratosEncontrados int
	UltimaAtualizacao        time.Time
}
