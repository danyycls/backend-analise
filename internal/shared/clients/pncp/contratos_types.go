package pncp

import "github.com/danyele/podp/internal/shared/types"

type CategoriaProcesso struct {
	ID   *int    `json:"id"`
	Nome *string `json:"nome"`
}

type TipoContrato struct {
	ID   *int    `json:"id"`
	Nome *string `json:"nome"`
}

type OrgaoEntidade struct {
	CNPJ        *string `json:"cnpj"`
	EsferaID    *string `json:"esferaId"`
	PoderID     *string `json:"poderId"`
	RazaoSocial *string `json:"razaoSocial"`
}

type UnidadeOrgao struct {
	CodigoIbge    *string `json:"codigoIbge"`
	CodigoUnidade *string `json:"codigoUnidade"`
	MunicipioNome *string `json:"municipioNome"`
	NomeUnidade   *string `json:"nomeUnidade"`
	UFNome        *string `json:"ufNome"`
	UFSigla       *string `json:"ufSigla"`
}

type Contrato struct {
	AnoContrato               *int                      `json:"anoContrato"`
	CategoriaProcesso         *CategoriaProcesso        `json:"categoriaProcesso"`
	CNPJOrgao                 *string                   `json:"cnpjOrgao"`
	CNPJOrgaoSub              *string                   `json:"cnpjOrgaoSub"`
	CodigoContrato            *string                   `json:"codigoContrato"`
	CodigoFonteOrcamentaria   *int                      `json:"codigoFonteOrcamentaria"`
	CodigoOrgao               *string                   `json:"codigoOrgao"`
	CodigoTipoContrato        *int                      `json:"codigoTipoContrato"`
	CodigoUG                  *string                   `json:"codigoUg"`
	DataAssinatura            *string                   `json:"dataAssinatura"`
	DataInicioVigencia        *string                   `json:"dataVigenciaInicio"`
	DataPublicacao            *string                   `json:"dataPublicacaoPncp"`
	DataTerminoVigencia       *string                   `json:"dataVigenciaFim"`
	FonteOrcamentaria         *FonteOrcamentaria        `json:"fonteOrcamentaria"`
	Fornecedor                *types.FornecedorOpenCNPJ `json:"fornecedor"`
	ModalidadeNome            *string                   `json:"modalidadeNome"`
	NIFornecedor              *string                   `json:"niFornecedor"`
	NomeRazaoSocialFornecedor *string                   `json:"nomeRazaoSocialFornecedor"`
	NomeOrgao                 *string                   `json:"nomeOrgao"`
	NomeOrgaoSub              *string                   `json:"nomeOrgaoSub"`
	NumeroControlePNCP        *string                   `json:"numeroControlePNCP"`
	NumeroContrato            *string                   `json:"numeroContrato"`
	NumeroCNPJ                *string                   `json:"numeroCNPJ"`
	NumeroCPF                 *string                   `json:"numeroCPF"`
	NumeroLicitação           *string                   `json:"numeroLicitacao"`
	ObjetoCompra              *string                   `json:"objetoContrato"`
	OrgaoEntidade             *OrgaoEntidade            `json:"orgaoEntidade"`
	OrgaoSub                  *OrgaoEntidade            `json:"orgaoSub"`
	OrgaoVinculado            *UnidadeOrgao             `json:"orgaoVinculado"`
	OrigemLicitação           *string                   `json:"origemLicitacao"`
	PrazoInicioVigencia       *string                   `json:"prazoInicioVigencia"`
	PrazoTerminoVigencia      *string                   `json:"prazoTerminoVigencia"`
	Produto                   *string                   `json:"produto"`
	Srp                       interface{}               `json:"srp,omitempty"`
	SubtipoContrato           *string                   `json:"subtipoContrato"`
	TipoContrato              *TipoContrato             `json:"tipoContrato"`
	UG                        *UnidadeOrgao             `json:"unidadeOrgao"`
	UnidadeSub                *UnidadeOrgao             `json:"unidadeSub"`
	ValorGlobal               *float64                  `json:"valorGlobal"`
	ValorInicial              *float64                  `json:"valorInicial"`
	ValorParcela              *float64                  `json:"valorParcela"`
	ValorTotalEstimado        *float64                  `json:"valorTotalEstimado"`
	ValorTotalHomologado      *float64                  `json:"valorTotalHomologado"`
	AmparoLegal               *AmparoLegal              `json:"amparoLegal"`
}

type OrgaoInfo struct {
	CNPJ        *string `json:"cnpj"`
	RazaoSocial *string `json:"razaoSocial"`
}

type Periodo struct {
	DataInicial *string `json:"dataInicial"`
	DataFinal   *string `json:"dataFinal"`
}

type Resumo struct {
	TotalContratos      *int     `json:"totalContratos"`
	TotalEmpresas       *int     `json:"totalEmpresas"`
	ValorTotalContratos *float64 `json:"valorTotalContratos"`
}

type AnaliseResultado struct {
	Orgao     *OrgaoInfo `json:"orgao,omitempty"`
	Periodo   *Periodo   `json:"periodo,omitempty"`
	Resumo    *Resumo    `json:"resumo,omitempty"`
	Contratos []Contrato `json:"contratos,omitempty"`
}

type ContratoResponse struct {
	TotalPaginas     int        `json:"totalPaginas"`
	TotalRegistros   int        `json:"totalRegistros"`
	NumeroPagina     int        `json:"numeroPagina"`
	PaginasRestantes int        `json:"paginasRestantes"`
	Empty            bool       `json:"empty"`
	Data             []Contrato `json:"data"`
}

type AnaliseContratoRequest struct {
	Tipo                string `json:"tipo"`
	UF                  string `json:"uf"`
	CodigoMunicipioIbge string `json:"codigo_municipio_ibge"`
	DataInicial         string `json:"data_inicial"`
	DataFinal           string `json:"data_final"`
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
