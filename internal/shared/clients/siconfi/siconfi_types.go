package siconfi

import (
	"encoding/json"
	"fmt"
)

// Response[T] é o envelope de paginação comum a todos os endpoints da API SICONFI.
// Todas as respostas seguem o formato { items, hasMore, count, limit, offset, links }.
type Response[T any] struct {
	Items   []T    `json:"items"`
	HasMore bool   `json:"hasMore"`
	Count   int    `json:"count"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Links   []Link `json:"links"`
}

// Link representa um link de navegação retornado pela API (self, describedby, first, next).
type Link struct {
	Rel  string `json:"rel"`
	HRef string `json:"href"`
}

// Decode extrai a lista de itens do envelope de paginação genérico da SICONFI.
// T é o tipo do item (Ente, DCAItem, RGFItem, etc.).
func Decode[T any](data []byte) ([]T, error) {
	var resp Response[T]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("siconfi: erro ao decodificar resposta: %w", err)
	}
	return resp.Items, nil
}

// Ente representa as informações básicas de cadastro de um ente da federação (estado ou município).
// Retornado pelo endpoint /entes.
type Ente struct {
	CodIBGE   int    `json:"cod_ibge"`
	Ente      string `json:"ente"`
	Capital   string `json:"capital"`
	Regiao    string `json:"regiao"`
	UF        string `json:"uf"`
	Esfera    string `json:"esfera"`
	Exercicio int    `json:"exercicio"`
	Populacao int    `json:"populacao"`
	CNPJ      string `json:"cnpj"`
}

// DCAItem representa uma linha dos quadros da Declaração das Contas Anuais (DCA)
// ou do antigo Quadro de Detalhamento das Contas Contábeis (QDCC).
// Cada item é uma conta contábil com seu valor para um ente e exercício específicos.
// Retornado pelo endpoint /dca.
type DCAItem struct {
	Exercicio   int     `json:"exercicio"`
	Instituicao string  `json:"instituicao"`
	CodIBGE     int     `json:"cod_ibge"`
	UF          string  `json:"uf"`
	Anexo       string  `json:"anexo"`
	Rotulo      string  `json:"rotulo"`
	Coluna      string  `json:"coluna"`
	CodConta    string  `json:"cod_conta"`
	Conta       string  `json:"conta"`
	Valor       float64 `json:"valor"`
	Populacao   int     `json:"populacao"`
}

// RGFItem representa uma linha dos anexos do Relatório de Gestão Fiscal (RGF).
// Contém os indicadores fiscais exigidos pela Lei de Responsabilidade Fiscal (LRF)
// como despesa com pessoal, dívida consolidada, operações de crédito, etc.
// Retornado pelo endpoint /rgf.
type RGFItem struct {
	Exercicio     int     `json:"exercicio"`
	Periodo       int     `json:"periodo"`
	Periodicidade string  `json:"periodicidade"`
	Instituicao   string  `json:"instituicao"`
	CodIBGE       int     `json:"cod_ibge"`
	UF            string  `json:"uf"`
	CoPoder       string  `json:"co_poder"`
	Populacao     int     `json:"populacao"`
	Anexo         string  `json:"anexo"`
	Esfera        string  `json:"esfera"`
	Rotulo        string  `json:"rotulo"`
	Coluna        string  `json:"coluna"`
	CodConta      string  `json:"cod_conta"`
	Conta         string  `json:"conta"`
	Valor         float64 `json:"valor"`
}

// RREOItem representa uma linha dos anexos do Relatório Resumido da Execução
// Orçamentária (RREO). Mostra a execução do orçamento por bimestre: receita
// prevista vs realizada, despesa fixada vs executada, restos a pagar, etc.
// Retornado pelo endpoint /rreo.
type RREOItem struct {
	Exercicio     int     `json:"exercicio"`
	Demonstration string  `json:"demonstration"`
	Periodo       int     `json:"periodo"`
	Periodicidade string  `json:"periodicidade"`
	Instituicao   string  `json:"instituicao"`
	CodIBGE       int     `json:"cod_ibge"`
	UF            string  `json:"uf"`
	Populacao     int     `json:"populacao"`
	Anexo         string  `json:"anexo"`
	Esfera        string  `json:"esfera"`
	Rotulo        string  `json:"rotulo"`
	Coluna        string  `json:"coluna"`
	CodConta      string  `json:"cod_conta"`
	Conta         string  `json:"conta"`
	Valor         float64 `json:"valor"`
}

// MSCItem representa uma linha da Matriz de Saldos Contábeis (MSC).
// Detalha os lançamentos contábeis no nível de conta contábil (rubrica) para um ente,
// exercício e mês específicos. Abrange as classes 1 a 8 do PCASP.
// Usado pelos endpoints /msc_patrimonial, /msc_orcamentaria e /msc_controle.
type MSCItem struct {
	TipoMatriz       string  `json:"tipo_matriz"`
	CodIBGE          int     `json:"cod_ibge"`
	ClasseConta      int     `json:"classe_conta"`
	ContaContabil    string  `json:"conta_contabil"`
	PoderOrgao       string  `json:"poder_orgao"`
	AnoFonteRecursos *int    `json:"ano_fonte_recursos"`
	FonteRecursos    *string `json:"fonte_recursos"`
	Funcao           *string `json:"funcao"`
	Subfuncao        *string `json:"subfuncao"`
	Exercicio        int     `json:"exercicio"`
	MesReferencia    int     `json:"mes_referencia"`
	EducacaoSaude    *int    `json:"educacao_saude"`
	DataReferencia   string  `json:"data_referencia"`
	EntradaMSC       int     `json:"entrada_msc"`
	NaturezaDespesa  *string `json:"natureza_despesa"`
	AnoInscricao     *int    `json:"ano_inscricao"`
	NaturezaReceita  *string `json:"natureza_receita"`
	Valor            float64 `json:"valor"`
	NaturezaConta    string  `json:"natureza_conta"`
	TipoValor        string  `json:"tipo_valor"`
	ComplementoFonte *string `json:"complemento_fonte"`
}

// ExtratoEntregasItem representa uma linha do extrato de entregas do SICONFI.
// Informa o status de envio/homologação/retificação dos relatórios e matrizes
// de um ente em um determinado exercício.
// Retornado pelo endpoint /extrato_entregas.
type ExtratoEntregasItem struct {
	Exercicio       int     `json:"exercicio"`
	CodIBGE         int     `json:"cod_ibge"`
	Populacao       int     `json:"populacao"`
	Instituicao     string  `json:"instituicao"`
	Entregavel      string  `json:"entregavel"`
	Periodo         int     `json:"periodo"`
	Periodicidade   string  `json:"periodicidade"`
	StatusRelatorio *string `json:"status_relatorio"`
	DataStatus      string  `json:"data_status"`
	FormaEnvio      string  `json:"forma_envio"`
	TipoRelatorio   *string `json:"tipo_relatorio"`
}

// AnexoRelatorio representa um anexo de relatório/demonstration agrupado por
// esfera de governo (E=Estados, M=Municípios, U=União, C=Consórcio).
// Serve como tabela de apoio para mapear quais anexos pertencem a cada demonstration.
// Retornado pelo endpoint /anexos-relatorios.
type AnexoRelatorio struct {
	Esfera        string `json:"esfera"`
	Demonstration string `json:"demonstration"`
	Anexo         string `json:"anexo"`
}
