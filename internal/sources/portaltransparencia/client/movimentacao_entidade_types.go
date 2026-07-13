package portaltransparencia

type DespesaAnualPorOrgao struct {
	Ano                 int    `json:"ano"`
	Orgao               string `json:"orgao"`
	CodigoOrgao         string `json:"codigoOrgao"`
	OrgaoSuperior       string `json:"orgaoSuperior"`
	CodigoOrgaoSuperior string `json:"codigoOrgaoSuperior"`
	Empenhado           string `json:"empenhado"`
	Liquidado           string `json:"liquidado"`
	Pago                string `json:"pago"`
}

type DespesaAnualPorFuncaoESubfuncao struct {
	Ano             int    `json:"ano"`
	Funcao          string `json:"funcao"`
	CodigoFuncao    string `json:"codigoFuncao"`
	Subfuncao       string `json:"subfuncao"`
	CodigoSubfuncao string `json:"codigoSubfuncao"`
	Programa        string `json:"programa"`
	CodigoPrograma  string `json:"codigoPrograma"`
	Acao            string `json:"acao"`
	CodigoAcao      string `json:"codigoAcao"`
	Empenhado       string `json:"empenhado"`
	Liquidado       string `json:"liquidado"`
	Pago            string `json:"pago"`
}

type DespesaLiquidaAnualPorFuncaoESubfuncao struct {
	Ano                     int    `json:"ano"`
	Funcao                  string `json:"funcao"`
	CodigoFuncao            string `json:"codigoFuncao"`
	Subfuncao               string `json:"subfuncao"`
	CodigoSubfuncao         string `json:"codigoSubfuncao"`
	Programa                string `json:"programa"`
	CodigoPrograma          string `json:"codigoPrograma"`
	Acao                    string `json:"acao"`
	CodigoAcao              string `json:"codigoAcao"`
	PlanoOrcamentario       string `json:"planoOrcamentario"`
	IDPlanoOrcamentario     int    `json:"idPlanoOrcamentario"`
	CodigoPlanoOrcamentario string `json:"codigoPlanoOrcamentario"`
	GrupoDespesa            string `json:"grupoDespesa"`
	CodigoGrupoDespesa      string `json:"codigoGrupoDespesa"`
	ElementoDespesa         string `json:"elementoDespesa"`
	CodigoElementoDespesa   string `json:"codigoElementoDespesa"`
	ModalidadeDespesa       string `json:"modalidadeDespesa"`
	CodigoModalidadeDespesa string `json:"codigoModalidadeDespesa"`
	Empenhado               string `json:"empenhado"`
	Liquidado               string `json:"liquidado"`
	Pago                    string `json:"pago"`
}

type DespesasPorPlanoOrcamentario struct {
	ID                     int    `json:"id"`
	Codigo                 string `json:"codigo"`
	Descricao              string `json:"descricao"`
	CodUnidadeOrcamentaria string `json:"codUnidadeOrcamentaria"`
	CodigoFuncao           string `json:"codigoFuncao"`
	CodigoSubFuncao        string `json:"codigoSubFuncao"`
	CodigoPrograma         string `json:"codigoPrograma"`
	CodigoAcao             string `json:"codigoAcao"`
	CodPOIdAcompanhamento  string `json:"codPOIdAcompanhamento"`
	DescPOIdAcompanhamento string `json:"descPOIdAcompanhamento"`
	NumAno                 int    `json:"numAno"`
}

type DetalhamentoDoGasto struct {
	CodigoItemEmpenho    string `json:"codigoItemEmpenho"`
	Descricao            string `json:"descricao"`
	CodigoSubelemento    string `json:"codigoSubelemento"`
	DescricaoSubelemento string `json:"descricaoSubelemento"`
	ValorAtual           string `json:"valorAtual"`
	Sequencial           int    `json:"sequencial"`
}

type HistoricoSubItemEmpenho struct {
	Data          string `json:"data"`
	Operacao      string `json:"operacao"`
	Quantidade    string `json:"quantidade"`
	ValorUnitario string `json:"valorUnitario"`
	ValorTotal    string `json:"valorTotal"`
}

type FuncionalProgramatica struct {
	ID              int    `json:"id"`
	CodigoFuncao    string `json:"codigoFuncao"`
	CodigoSubfuncao string `json:"codigoSubfuncao"`
	CodigoPrograma  string `json:"codigoPrograma"`
	CodigoAcao      string `json:"codigoAcao"`
	Ano             int    `json:"ano"`
}

type Funcao struct {
	CodigoFuncao    string `json:"codigoFuncao"`
	DescricaoFuncao string `json:"descricaoFuncao"`
}

type Subfuncao struct {
	CodigoSubfuncao    string `json:"codigoSubfuncao"`
	DescricaoSubfuncao string `json:"descricaoSubfuncao"`
}

type CodigoDescricao struct {
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
}

type ConsultaFavorecidosFinaisPorDocumento struct {
	SkFatDW                  int    `json:"skFatDW"`
	CodigoPagamento          string `json:"codigoPagamento"`
	CodigoListaCredor        string `json:"codigoListaCredor"`
	ValorFinal               string `json:"valorFinal"`
	TipoOB                   string `json:"tipoOB"`
	TipoDocumento            string `json:"tipoDocumento"`
	DataCarga                string `json:"dataCarga"`
	SkPessoaFinal            int    `json:"skPessoaFinal"`
	CodigoFavorecidoFinal    string `json:"codigoFavorecidoFinal"`
	NomeFavorecidoFinal      string `json:"nomeFavorecidoFinal"`
	TipoFavorecidoFinal      string `json:"tipoFavorecidoFinal"`
	UFFavorecidoFinal        string `json:"ufFavorecidoFinal"`
	MunicipioFavorecidoFinal string `json:"municipioFavorecidoFinal"`
	SkPessoaDespesa          int    `json:"skPessoaDespesa"`
	CodigoFavorecidoDespesa  string `json:"codigoFavorecidoDespesa"`
	NomeFavorecidoDespesa    string `json:"nomeFavorecidoDespesa"`
	TipoFavorecidoDespesa    string `json:"tipoFavorecidoDespesa"`
	CodigoOrgaoSuperior      string `json:"codigoOrgaoSuperior"`
	OrgaoSuperior            string `json:"orgaoSuperior"`
	CodigoOrgaoVinculado     string `json:"codigoOrgaoVinculado"`
	OrgaoVinculado           string `json:"orgaoVinculado"`
	CodigoUnidadeGestora     string `json:"codigoUnidadeGestora"`
	UnidadeGestora           string `json:"unidadeGestora"`
}

type EmpenhoImpactadoBasico struct {
	Empenho             string `json:"empenho"`
	Subitem             string `json:"subitem"`
	EmpenhoResumido     string `json:"empenhoResumido"`
	ValorLiquidado      string `json:"valorLiquidado"`
	ValorPago           string `json:"valorPago"`
	ValorRestoInscrito  string `json:"valorRestoInscrito"`
	ValorRestoCancelado string `json:"valorRestoCancelado"`
	ValorRestoPago      string `json:"valorRestoPago"`
}

type DespesasPorDocumento struct {
	Data                    string `json:"data"`
	Documento               string `json:"documento"`
	DocumentoResumido       string `json:"documentoResumido"`
	Observacao              string `json:"observacao"`
	Funcao                  string `json:"funcao"`
	Subfuncao               string `json:"subfuncao"`
	Programa                string `json:"programa"`
	Acao                    string `json:"acao"`
	SubTitulo               string `json:"subTitulo"`
	LocalizadorGasto        string `json:"localizadorGasto"`
	Fase                    string `json:"fase"`
	Especie                 string `json:"especie"`
	Favorecido              string `json:"favorecido"`
	CodigoFavorecido        string `json:"codigoFavorecido"`
	NomeFavorecido          string `json:"nomeFavorecido"`
	UFFavorecido            string `json:"ufFavorecido"`
	Valor                   string `json:"valor"`
	CodigoUG                string `json:"codigoUg"`
	UG                      string `json:"ug"`
	CodigoUO                string `json:"codigoUo"`
	UO                      string `json:"uo"`
	CodigoOrgao             string `json:"codigoOrgao"`
	Orgao                   string `json:"orgao"`
	CodigoOrgaoSuperior     string `json:"codigoOrgaoSuperior"`
	OrgaoSuperior           string `json:"orgaoSuperior"`
	Categoria               string `json:"categoria"`
	Grupo                   string `json:"grupo"`
	Elemento                string `json:"elemento"`
	Modalidade              string `json:"modalidade"`
	NumeroProcesso          string `json:"numeroProcesso"`
	PlanoOrcamentario       string `json:"planoOrcamentario"`
	Author                  string `json:"author"`
	FavorecidoIntermediario bool   `json:"favorecidoIntermediario"`
	FavorecidoListaFaturas  bool   `json:"favorecidoListaFaturas"`
}

type DocumentoRelacionado struct {
	Data              string `json:"data"`
	Fase              string `json:"fase"`
	Documento         string `json:"documento"`
	DocumentoResumido string `json:"documentoResumido"`
	Especie           string `json:"especie"`
	OrgaoSuperior     string `json:"orgaoSuperior"`
	OrgaoVinculado    string `json:"orgaoVinculado"`
	UnidadeGestora    string `json:"unidadeGestora"`
	ElementoDespesa   string `json:"elementoDespesa"`
	Favorecido        string `json:"favorecido"`
	Valor             string `json:"valor"`
}

type PessoaRecursosRecebidosUGMesDesnormalizada struct {
	AnoMes              int     `json:"anoMes"`
	CodigoPessoa        string  `json:"codigoPessoa"`
	NomePessoa          string  `json:"nomePessoa"`
	TipoPessoa          string  `json:"tipoPessoa"`
	MunicipioPessoa     string  `json:"municipioPessoa"`
	SiglaUFPessoa       string  `json:"siglaUFPessoa"`
	CodigoUG            string  `json:"codigoUG"`
	NomeUG              string  `json:"nomeUG"`
	CodigoOrgao         string  `json:"codigoOrgao"`
	NomeOrgao           string  `json:"nomeOrgao"`
	CodigoOrgaoSuperior string  `json:"codigoOrgaoSuperior"`
	NomeOrgaoSuperior   string  `json:"nomeOrgaoSuperior"`
	Valor               float64 `json:"valor"`
}

type ConsultaEmendas struct {
	CodigoEmenda        string `json:"codigoEmenda"`
	Ano                 int    `json:"ano"`
	TipoEmenda          string `json:"tipoEmenda"`
	Author              string `json:"author"`
	NomeAutor           string `json:"nomeAutor"`
	NumeroEmenda        string `json:"numeroEmenda"`
	LocalidadeDoGasto   string `json:"localidadeDoGasto"`
	Funcao              string `json:"funcao"`
	Subfuncao           string `json:"subfuncao"`
	ValorEmpenhado      string `json:"valorEmpenhado"`
	ValorLiquidado      string `json:"valorLiquidado"`
	ValorPago           string `json:"valorPago"`
	ValorRestoInscrito  string `json:"valorRestoInscrito"`
	ValorRestoCancelado string `json:"valorRestoCancelado"`
	ValorRestoPago      string `json:"valorRestoPago"`
}

type DocumentoRelacionadoEmenda struct {
	ID                      int    `json:"id"`
	Data                    string `json:"data"`
	Fase                    string `json:"fase"`
	CodigoDocumento         string `json:"codigoDocumento"`
	CodigoDocumentoResumido string `json:"codigoDocumentoResumido"`
	EspecieTipo             string `json:"especieTipo"`
	TipoEmenda              string `json:"tipoEmenda"`
}
