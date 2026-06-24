package siconfi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
)

var ErrSICONFIIndisponivel = errors.New("API SICONFI temporariamente indisponivel")

type siconfiErrorResp struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// SICONFIClient é o cliente HTTP para a API de dados abertos do SICONFI
// (Sistema de Informações Contábeis e Fiscais do Setor Público Brasileiro)
// da Secretaria do Tesouro Nacional.
//
// Rate limit: 1 requisição por segundo.
// Paginação padrão: 5000 itens por página.
type SICONFIClient struct {
	baseURL string
	client  *http.Client
}

// NovoSICONFIClient cria uma nova instância do cliente SICONFI com timeout de 30 segundos.
func NovoSICONFIClient(baseURL string) *SICONFIClient {
	return &SICONFIClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// RGFParams agrupa os parâmetros obrigatórios e opcionais para consulta do
// Relatório de Gestão Fiscal (RGF) via endpoint /rgf.
type RGFParams struct {
	// AnExercicio é o exercício do relatório (ex: 2023). Obrigatório.
	AnExercicio int64
	// InPeriodicidade é a periodicidade de publicação: Q = quadrimestral, S = semestral.
	// A semestral aplica-se apenas a municípios com menos de 50 mil habitantes.
	// Obrigatório.
	InPeriodicidade string
	// NrPeriodo é o quadrimestre (1-3) ou semester (1-2) de referência. Obrigatório.
	NrPeriodo int
	// CoTipoDemonstrativo é o tipo de demonstration (ex: "RGF"). Obrigatório.
	CoTipoDemonstrativo string
	// CoPoder é o código do poder: E = Executivo, L = Legislation, J = Judiciário,
	// M = Ministério Público, D = Defensoria Pública. Obrigatório.
	CoPoder string
	// IdEnte é o código IBGE de 7 dígitos do ente. Obrigatório.
	IdEnte int
	// NoAnexo filtra por um anexo específico (ex: "RGF-Anexo 01"). Opcional.
	NoAnexo string
	// CoEsfera filtra por esfera: M = Municípios, E = Estados e DF, U = União, C = Consórcio. Opcional.
	CoEsfera string
}

// RREOParams agrupa os parâmetros obrigatórios e opcionais para consulta do
// Relatório Resumido da Execução Orçamentária (RREO) via endpoint /rreo.
type RREOParams struct {
	// AnExercicio é o exercício do relatório (ex: 2023). Obrigatório.
	AnExercicio int64
	// NrPeriodo é o bimestre de referência (1-6). Obrigatório.
	NrPeriodo int
	// CoTipoDemonstrativo é o tipo de demonstration (ex: "RREO"). Obrigatório.
	CoTipoDemonstrativo string
	// IdEnte é o código IBGE de 7 dígitos do ente. Obrigatório.
	IdEnte int
	// NoAnexo filtra por um anexo específico (ex: "RREO-Anexo 01"). Opcional.
	NoAnexo string
	// CoEsfera filtra por esfera: M = Municípios, E = Estados e DF, U = União, C = Consórcio. Opcional.
	CoEsfera string
}

// MSCParams agrupa os parâmetros obrigatórios para consulta das Matrizes de Saldos
// Contábeis (MSC) via endpoints /msc_patrimonial, /msc_orcamentaria e /msc_controle.
// Todos os campos são obrigatórios.
type MSCParams struct {
	// IdEnte é o código IBGE de 7 dígitos do ente.
	IdEnte int
	// AnReferencia é o exercício de referência da matriz (ex: 2023).
	AnReferencia int64
	// MeReferencia é o mês de referência (1-12) ou 13 para encerramento do exercício.
	MeReferencia int64
	// CoTipoMatriz é o tipo de matriz: MSCC (mensal/agregada) ou MSCE (encerramento).
	CoTipoMatriz string
	// ClasseConta é a classe de contas do PCASP:
	// 1=Ativo, 2=Passivo, 3=VPD, 4=VPA (patrimonial),
	// 5=Orçamento aprovado, 6=Execução do orçamento (orçamentária),
	// 7=Controls devedores, 8=Controls credores (controle).
	ClasseConta int
	// IdTV é o tipo de valor: beginning_balance (saldo inicial),
	// ending_balance (saldo final) ou period_change (movimento).
	IdTV string
}

// doGet executa uma requisição HTTP GET para o path informado com os query params
// e retorna o corpo da resposta em bytes. Método interno compartilhado por todos
// os endpoints para evitar duplicação de lógica HTTP.
func (c *SICONFIClient) doGet(ctx context.Context, path string, query url.Values) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("siconfi: erro ao montar url: %w", err)
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("siconfi: erro ao criar requisição: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("siconfi: erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var sErr siconfiErrorResp
		if err := json.Unmarshal(body, &sErr); err == nil && sErr.Code == "AccountIsLocked" {
			return nil, fmt.Errorf("siconfi: API temporariamente indisponivel (conta bloqueada): %w", ErrSICONFIIndisponivel)
		}
		return nil, fmt.Errorf("siconfi: status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("siconfi: erro ao ler resposta: %w", err)
	}

	return body, nil
}

// ListarEntes retorna as informações básicas de cadastro de todos os entes da
// federação (estados e municípios) cadastrados no SICONFI.
// Inclui código IBGE, nome, UF, região, esfera, população, CNPJ e capital.
// Endpoint: GET /entes
func (c *SICONFIClient) ListarEntes(ctx context.Context) ([]Ente, error) {
	log := logger.New("Clients: Client: ListarEntes")
	log.Info("SICONFI: solicitando lista de entes")

	body, err := c.doGet(ctx, "/entes", nil)
	if err != nil {
		return nil, fmt.Errorf("listar entes: %w", err)
	}

	itens, err := Decode[Ente](body)
	if err != nil {
		return nil, fmt.Errorf("listar entes: %w", err)
	}

	return itens, nil
}

// BuscarDCA retorna os dados contidos nos quadros da Declaração das Contas Anuais (DCA)
// para um ente e exercício específicos. Para o exercício de 2013, inclui também dados
// do antigo QDCC (Quadro de Detalhamento das Contas Contábeis).
//
// Parâmetros:
//   - anExercicio: exercício do relatório (ex: 2023)
//   - idEnte: código IBGE de 7 dígitos do ente
//   - noAnexo: opcional, filtra por um anexo específico (ex: "DCA-Anexo I-AB")
//
// Endpoint: GET /dca
func (c *SICONFIClient) BuscarDCA(ctx context.Context, anExercicio int64, idEnte int, noAnexo ...string) ([]DCAItem, error) {
	log := logger.New("Clients: Client: BuscarDCA")
	log.Info("SICONFI: solicitando DCA", "ente", idEnte, "exercicio", anExercicio)

	query := url.Values{}
	query.Set("an_exercicio", strconv.FormatInt(anExercicio, 10))
	query.Set("id_ente", strconv.Itoa(idEnte))
	if len(noAnexo) > 0 && noAnexo[0] != "" {
		query.Set("no_anexo", noAnexo[0])
	}

	body, err := c.doGet(ctx, "/dca", query)
	if err != nil {
		return nil, fmt.Errorf("buscar dca: %w", err)
	}

	itens, err := Decode[DCAItem](body)
	if err != nil {
		return nil, fmt.Errorf("buscar dca: %w", err)
	}

	return itens, nil
}

// BuscarRGF retorna os dados contidos nos anexos do Relatório de Gestão Fiscal (RGF)
// para um poder/órgão e período específicos.
//
// O RGF é o principal instrumento de transparência da gestão fiscal exigido pela LRF.
// Contém indicadores como: despesa com pessoal, dívida consolidada líquida,
// garantias, operações de crédito e disponibilidade de caixa.
//
// Parâmetros obrigatórios: AnExercicio, InPeriodicidade, NrPeriodo,
// CoTipoDemonstrativo, CoPoder, IdEnte.
// Parâmetros opcionais: NoAnexo, CoEsfera.
//
// Endpoint: GET /rgf
func (c *SICONFIClient) BuscarRGF(ctx context.Context, params RGFParams) ([]RGFItem, error) {
	log := logger.New("Clients: Client: BuscarRGF")
	log.Info("SICONFI: solicitando RGF", "ente", params.IdEnte, "exercicio", params.AnExercicio, "periodo", params.NrPeriodo)

	query := url.Values{}
	query.Set("an_exercicio", strconv.FormatInt(params.AnExercicio, 10))
	query.Set("in_periodicidade", params.InPeriodicidade)
	query.Set("nr_periodo", strconv.Itoa(params.NrPeriodo))
	query.Set("co_tipo_demonstration", params.CoTipoDemonstrativo)
	query.Set("co_poder", params.CoPoder)
	query.Set("id_ente", strconv.Itoa(params.IdEnte))
	if params.NoAnexo != "" {
		query.Set("no_anexo", params.NoAnexo)
	}
	if params.CoEsfera != "" {
		query.Set("co_esfera", params.CoEsfera)
	}

	body, err := c.doGet(ctx, "/rgf", query)
	if err != nil {
		return nil, fmt.Errorf("buscar rgf: %w", err)
	}

	itens, err := Decode[RGFItem](body)
	if err != nil {
		return nil, fmt.Errorf("buscar rgf: %w", err)
	}

	return itens, nil
}

// BuscarRREO retorna os dados contidos nos anexos do Relatório Resumido da Execução
// Orçamentária (RREO) para um ente e período específicos.
//
// O RREO é publicado a cada bimestre e demonstra a execução orçamentária resumida:
// balanço orçamentário, execução de restos a pagar, receita corrente líquida,
// despesas por função/subfunção e receitas/despesas previdenciárias.
// O RREO Simplificado aplica-se apenas a municípios com menos de 50 mil habitantes.
//
// Parâmetros obrigatórios: AnExercicio, NrPeriodo, CoTipoDemonstrativo, IdEnte.
// Parâmetros opcionais: NoAnexo, CoEsfera.
//
// Endpoint: GET /rreo
func (c *SICONFIClient) BuscarRREO(ctx context.Context, params RREOParams) ([]RREOItem, error) {
	log := logger.New("Clients: Client: BuscarRREO")
	log.Info("SICONFI: solicitando RREO", "ente", params.IdEnte, "exercicio", params.AnExercicio, "periodo", params.NrPeriodo)

	query := url.Values{}
	query.Set("an_exercicio", strconv.FormatInt(params.AnExercicio, 10))
	query.Set("nr_periodo", strconv.Itoa(params.NrPeriodo))
	query.Set("co_tipo_demonstration", params.CoTipoDemonstrativo)
	query.Set("id_ente", strconv.Itoa(params.IdEnte))
	if params.NoAnexo != "" {
		query.Set("no_anexo", params.NoAnexo)
	}
	if params.CoEsfera != "" {
		query.Set("co_esfera", params.CoEsfera)
	}

	body, err := c.doGet(ctx, "/rreo", query)
	if err != nil {
		return nil, fmt.Errorf("buscar rreo: %w", err)
	}

	itens, err := Decode[RREOItem](body)
	if err != nil {
		return nil, fmt.Errorf("buscar rreo: %w", err)
	}

	return itens, nil
}

// BuscarMSCPatrimonial retorna o detalhamento dos registros contábeis de natureza
// patrimonial (classes 1 a 4 do PCASP) para um ente, exercício e mês específicos.
//
// Classes cobertas:
//   - 1: Ativo (bens e direitos)
//   - 2: Passivo (obrigações)
//   - 3: Variações Patrimoniais Diminutivas (VPD)
//   - 4: Variações Patrimoniais Aumentativas (VPA)
//
// Todos os campos de MSCParams são obrigatórios.
//
// Endpoint: GET /msc_patrimonial
func (c *SICONFIClient) BuscarMSCPatrimonial(ctx context.Context, params MSCParams) ([]MSCItem, error) {
	log := logger.New("Clients: Client: BuscarMSCPatrimonial")
	log.Info("SICONFI: solicitando MSC Patrimonial", "ente", params.IdEnte, "exercicio", params.AnReferencia, "mes", params.MeReferencia, "classe", params.ClasseConta)

	body, err := c.buscarMSC(ctx, "/msc_patrimonial", params)
	if err != nil {
		return nil, err
	}

	itens, err := Decode[MSCItem](body)
	if err != nil {
		return nil, fmt.Errorf("buscar msc patrimonial: %w", err)
	}

	return itens, nil
}

// BuscarMSCOrcamentaria retorna o detalhamento dos registros contábeis de natureza
// orçamentária (classes 5 e 6 do PCASP) para um ente, exercício e mês específicos.
//
// Classes cobertas:
//   - 5: Orçamento aprovado (planejamento/dotação)
//   - 6: Execução do orçamento (empenho, liquidação, pagamento)
//
// Todos os campos de MSCParams são obrigatórios.
//
// Endpoint: GET /msc_orcamentaria
func (c *SICONFIClient) BuscarMSCOrcamentaria(ctx context.Context, params MSCParams) ([]MSCItem, error) {
	log := logger.New("Clients: Client: BuscarMSCOrcamentaria")
	log.Info("SICONFI: solicitando MSC Orçamentária", "ente", params.IdEnte, "exercicio", params.AnReferencia, "mes", params.MeReferencia, "classe", params.ClasseConta)

	body, err := c.buscarMSC(ctx, "/msc_orcamentaria", params)
	if err != nil {
		return nil, err
	}

	itens, err := Decode[MSCItem](body)
	if err != nil {
		return nil, fmt.Errorf("buscar msc orcamentaria: %w", err)
	}

	return itens, nil
}

// BuscarMSCControle retorna o detalhamento dos registros contábeis de natureza de
// controle (classes 7 e 8 do PCASP) para um ente, exercício e mês específicos.
//
// Classes cobertas:
//   - 7: Controls devedores (execução de contratos, convênios, garantias)
//   - 8: Controls credores (execução de precatórios, obrigações contratuais)
//
// Todos os campos de MSCParams são obrigatórios.
//
// Endpoint: GET /msc_controle
func (c *SICONFIClient) BuscarMSCControle(ctx context.Context, params MSCParams) ([]MSCItem, error) {
	log := logger.New("Clients: Client: BuscarMSCControle")
	log.Info("SICONFI: solicitando MSC Controle", "ente", params.IdEnte, "exercicio", params.AnReferencia, "mes", params.MeReferencia, "classe", params.ClasseConta)

	body, err := c.buscarMSC(ctx, "/msc_controle", params)
	if err != nil {
		return nil, err
	}

	itens, err := Decode[MSCItem](body)
	if err != nil {
		return nil, fmt.Errorf("buscar msc controle: %w", err)
	}

	return itens, nil
}

// buscarMSC é o método interno compartilhado pelos três endpoints de MSC
// para montar os query params e executar a requisição.
func (c *SICONFIClient) buscarMSC(ctx context.Context, path string, params MSCParams) ([]byte, error) {
	query := url.Values{}
	query.Set("id_ente", strconv.Itoa(params.IdEnte))
	query.Set("an_referencia", strconv.FormatInt(params.AnReferencia, 10))
	query.Set("me_referencia", strconv.FormatInt(params.MeReferencia, 10))
	query.Set("co_tipo_matriz", params.CoTipoMatriz)
	query.Set("classe_conta", strconv.Itoa(params.ClasseConta))
	query.Set("id_tv", params.IdTV)

	return c.doGet(ctx, path, query)
}

// BuscarExtratoEntregas retorna o extrato de entregas de relatórios e matrizes
// para um ente e exercício. Informa quais relatórios foram homologados, retificados,
// a data de envio e a forma de envio (CSV, webservice, etc.).
//
// Parâmetros:
//   - idEnte: código IBGE de 7 dígitos do ente
//   - anReferencia: exercício de referência (ex: 2023)
//
// Endpoint: GET /extrato_entregas
func (c *SICONFIClient) BuscarExtratoEntregas(ctx context.Context, idEnte int, anReferencia int64) ([]ExtratoEntregasItem, error) {
	log := logger.New("Clients: Client: BuscarExtratoEntregas")
	log.Info("SICONFI: solicitando extrato de entregas", "ente", idEnte, "exercicio", anReferencia)

	query := url.Values{}
	query.Set("id_ente", strconv.Itoa(idEnte))
	query.Set("an_referencia", strconv.FormatInt(anReferencia, 10))

	body, err := c.doGet(ctx, "/extrato_entregas", query)
	if err != nil {
		return nil, fmt.Errorf("buscar extrato entregas: %w", err)
	}

	itens, err := Decode[ExtratoEntregasItem](body)
	if err != nil {
		return nil, fmt.Errorf("buscar extrato entregas: %w", err)
	}

	return itens, nil
}

// ListarAnexosRelatorios retorna a tabela de apoio com os anexos de cada
// relatório/demonstration agrupados por esfera de governo (E=Estados,
// M=Municípios, U=União, C=Consórcio).
//
// Útil para descobrir quais anexos estão disponíveis para cada tipo de
// demonstration (DCA, RGF, RREO, QDCC) antes de fazer consultas filtradas.
//
// Endpoint: GET /anexos-relatorios
func (c *SICONFIClient) ListarAnexosRelatorios(ctx context.Context) ([]AnexoRelatorio, error) {
	log := logger.New("Clients: Client: ListarAnexosRelatorios")
	log.Info("SICONFI: solicitando tabela de anexos dos relatórios")

	body, err := c.doGet(ctx, "/anexos-relatorios", nil)
	if err != nil {
		return nil, fmt.Errorf("listar anexos relatorios: %w", err)
	}

	itens, err := Decode[AnexoRelatorio](body)
	if err != nil {
		return nil, fmt.Errorf("listar anexos relatorios: %w", err)
	}

	return itens, nil
}
