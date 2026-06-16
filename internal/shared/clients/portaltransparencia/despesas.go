package portaltransparencia

import (
	"context"
	"strconv"
	"strings"
	"time"
)

func converterMesAnoISOparaBRRecursos(mesAnoISO string) string {
	if mesAnoISO == "" {
		return ""
	}
	t, err := time.Parse("2006-01", mesAnoISO)
	if err != nil {
		return strings.ReplaceAll(mesAnoISO, "-", "/")
	}
	return t.Format("01/2006")
}

func (c *PortalTransparenciaClient) ListarRecursosRecebidos(ctx context.Context, filtro DespesaRecursosRecebidosQueryParams) ([]PessoaRecursosRecebidosUGMesDesnormalizada, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.MesAnoInicio != "" {
		params["mesAnoInicio"] = converterMesAnoISOparaBRRecursos(filtro.MesAnoInicio)
	}
	if filtro.MesAnoFim != "" {
		params["mesAnoFim"] = converterMesAnoISOparaBRRecursos(filtro.MesAnoFim)
	}
	if filtro.NomeFavorecido != "" {
		params["nomeFavorecido"] = filtro.NomeFavorecido
	}
	if filtro.CodigoFavorecido != "" {
		params["codigoFavorecido"] = filtro.CodigoFavorecido
	}
	if filtro.TipoFavorecido != "" {
		params["tipoFavorecido"] = filtro.TipoFavorecido
	}
	if filtro.UF != "" {
		params["uf"] = filtro.UF
	}
	if filtro.CodigoIBGE != "" {
		params["codigoIBGE"] = filtro.CodigoIBGE
	}
	if filtro.OrgaoSuperior != "" {
		params["orgaoSuperior"] = filtro.OrgaoSuperior
	}
	if filtro.Orgao != "" {
		params["orgao"] = filtro.Orgao
	}
	if filtro.UnidadeGestora != "" {
		params["unidadeGestora"] = filtro.UnidadeGestora
	}
	var result []PessoaRecursosRecebidosUGMesDesnormalizada
	err := c.doGet(ctx, "/api-de-dados/despesas/recursos-recebidos", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarDespesasPorOrgao(ctx context.Context, filtro DespesaPorOrgaoQueryParams) ([]DespesaAnualPorOrgao, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
		"ano":    filtro.Ano,
	}
	if filtro.OrgaoSuperior != "" {
		params["orgaoSuperior"] = filtro.OrgaoSuperior
	}
	if filtro.Orgao != "" {
		params["orgao"] = filtro.Orgao
	}
	var result []DespesaAnualPorOrgao
	err := c.doGet(ctx, "/api-de-dados/despesas/por-orgao", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarDespesasPorFuncionalProgramatica(ctx context.Context, filtro DespesaFuncionalProgramaticaQueryParams) ([]DespesaAnualPorFuncaoESubfuncao, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
		"ano":    filtro.Ano,
	}
	if filtro.Funcao != "" {
		params["funcao"] = filtro.Funcao
	}
	if filtro.Subfuncao != "" {
		params["subfuncao"] = filtro.Subfuncao
	}
	if filtro.Programa != "" {
		params["programa"] = filtro.Programa
	}
	if filtro.Acao != "" {
		params["acao"] = filtro.Acao
	}
	var result []DespesaAnualPorFuncaoESubfuncao
	err := c.doGet(ctx, "/api-de-dados/despesas/por-funcional-programatica", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarDespesasMovimentacaoLiquida(ctx context.Context, filtro DespesaMovimentacaoLiquidaQueryParams) ([]DespesaLiquidaAnualPorFuncaoESubfuncao, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
		"ano":    filtro.Ano,
	}
	if filtro.Funcao != "" {
		params["funcao"] = filtro.Funcao
	}
	if filtro.Subfuncao != "" {
		params["subfuncao"] = filtro.Subfuncao
	}
	if filtro.Programa != "" {
		params["programa"] = filtro.Programa
	}
	if filtro.Acao != "" {
		params["acao"] = filtro.Acao
	}
	if filtro.GrupoDespesa != "" {
		params["grupoDespesa"] = filtro.GrupoDespesa
	}
	if filtro.ElementoDespesa != "" {
		params["elementoDespesa"] = filtro.ElementoDespesa
	}
	if filtro.ModalidadeAplicacao != "" {
		params["modalidadeAplicacao"] = filtro.ModalidadeAplicacao
	}
	if filtro.IDPlanoOrcamentario != "" {
		params["idPlanoOrcamentario"] = filtro.IDPlanoOrcamentario
	}
	var result []DespesaLiquidaAnualPorFuncaoESubfuncao
	err := c.doGet(ctx, "/api-de-dados/despesas/por-funcional-programatica/movimentacao-liquida", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarDespesasPlanoOrcamentario(ctx context.Context, filtro DespesaPlanoOrcamentarioQueryParams) ([]DespesasPorPlanoOrcamentario, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
		"ano":    filtro.Ano,
	}
	if filtro.CodPlanoOrcamentario != "" {
		params["codPlanoOrcamentario"] = filtro.CodPlanoOrcamentario
	}
	if filtro.DescPlanoOrcamentario != "" {
		params["descPlanoOrcamentario"] = filtro.DescPlanoOrcamentario
	}
	if filtro.CodPOIdentfAcompanhamento != "" {
		params["codPOIdentfAcompanhamento"] = filtro.CodPOIdentfAcompanhamento
	}
	var result []DespesasPorPlanoOrcamentario
	err := c.doGet(ctx, "/api-de-dados/despesas/plano-orcamentario", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarItensEmpenho(ctx context.Context, codigoDocumento string, pagina int) ([]DetalhamentoDoGasto, error) {
	params := map[string]string{
		"codigoDocumento": codigoDocumento,
		"pagina":          strconv.Itoa(pagina),
	}
	var result []DetalhamentoDoGasto
	err := c.doGet(ctx, "/api-de-dados/despesas/itens-de-empenho", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarHistoricoItemEmpenho(ctx context.Context, codigoDocumento string, sequencial int, pagina int) ([]HistoricoSubItemEmpenho, error) {
	params := map[string]string{
		"codigoDocumento": codigoDocumento,
		"sequencial":      strconv.Itoa(sequencial),
		"pagina":          strconv.Itoa(pagina),
	}
	var result []HistoricoSubItemEmpenho
	err := c.doGet(ctx, "/api-de-dados/despesas/itens-de-empenho/historico", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarSubfuncoes(ctx context.Context, filtro ListarFuncionalProgramaticaQueryParams) ([]Subfuncao, error) {
	params := map[string]string{
		"anoInicio": strconv.Itoa(filtro.AnoInicio),
		"pagina":    strconv.Itoa(filtro.Pagina),
	}
	if filtro.Codigo != "" {
		params["codigo"] = filtro.Codigo
	}
	var result []Subfuncao
	err := c.doGet(ctx, "/api-de-dados/despesas/funcional-programatica/subfuncoes", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarProgramas(ctx context.Context, filtro ListarFuncionalProgramaticaQueryParams) ([]CodigoDescricao, error) {
	params := map[string]string{
		"anoInicio": strconv.Itoa(filtro.AnoInicio),
		"pagina":    strconv.Itoa(filtro.Pagina),
	}
	if filtro.Codigo != "" {
		params["codigo"] = filtro.Codigo
	}
	var result []CodigoDescricao
	err := c.doGet(ctx, "/api-de-dados/despesas/funcional-programatica/programs", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarFuncionalProgramatica(ctx context.Context, ano int, pagina int) ([]FuncionalProgramatica, error) {
	params := map[string]string{
		"ano":    strconv.Itoa(ano),
		"pagina": strconv.Itoa(pagina),
	}
	var result []FuncionalProgramatica
	err := c.doGet(ctx, "/api-de-dados/despesas/funcional-programatica/listar", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarFuncoes(ctx context.Context, filtro ListarFuncionalProgramaticaQueryParams) ([]Funcao, error) {
	params := map[string]string{
		"anoInicio": strconv.Itoa(filtro.AnoInicio),
		"pagina":    strconv.Itoa(filtro.Pagina),
	}
	if filtro.Codigo != "" {
		params["codigo"] = filtro.Codigo
	}
	var result []Funcao
	err := c.doGet(ctx, "/api-de-dados/despesas/funcional-programatica/funcoes", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarAcoes(ctx context.Context, filtro ListarFuncionalProgramaticaQueryParams) ([]CodigoDescricao, error) {
	params := map[string]string{
		"anoInicio": strconv.Itoa(filtro.AnoInicio),
		"pagina":    strconv.Itoa(filtro.Pagina),
	}
	if filtro.Codigo != "" {
		params["codigo"] = filtro.Codigo
	}
	var result []CodigoDescricao
	err := c.doGet(ctx, "/api-de-dados/despesas/funcional-programatica/acoes", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarFavorecidosFinaisPorDocumento(ctx context.Context, codigoDocumento string, pagina int) ([]ConsultaFavorecidosFinaisPorDocumento, error) {
	params := map[string]string{
		"codigoDocumento": codigoDocumento,
		"pagina":          strconv.Itoa(pagina),
	}
	var result []ConsultaFavorecidosFinaisPorDocumento
	err := c.doGet(ctx, "/api-de-dados/despesas/favorecidos-finais-por-documento", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarEmpenhosImpactados(ctx context.Context, codigoDocumento string, fase string, pagina int) ([]EmpenhoImpactadoBasico, error) {
	params := map[string]string{
		"codigoDocumento": codigoDocumento,
		"fase":            fase,
		"pagina":          strconv.Itoa(pagina),
	}
	var result []EmpenhoImpactadoBasico
	err := c.doGet(ctx, "/api-de-dados/despesas/empenhos-impactados", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarDocumentos(ctx context.Context, filtro DespesaDocumentosQueryParams) ([]interface{}, error) {
	params := map[string]string{
		"dataEmissao": filtro.DataEmissao,
		"fase":        filtro.Fase,
		"pagina":      strconv.Itoa(filtro.Pagina),
	}
	if filtro.UnidadeGestora != "" {
		params["unidadeGestora"] = filtro.UnidadeGestora
	}
	if filtro.Gestao != "" {
		params["gestao"] = filtro.Gestao
	}
	var result []interface{}
	err := c.doGet(ctx, "/api-de-dados/despesas/documentos", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) BuscarDocumentoPorCodigo(ctx context.Context, codigo string) (*DespesasPorDocumento, error) {
	params := map[string]string{}
	var result DespesasPorDocumento
	err := c.doGet(ctx, "/api-de-dados/despesas/documentos/"+codigo, params, &result)
	return &result, err
}

func (c *PortalTransparenciaClient) ListarDocumentosRelacionados(ctx context.Context, codigoDocumento string, fase string) ([]DocumentoRelacionado, error) {
	params := map[string]string{
		"codigoDocumento": codigoDocumento,
		"fase":            fase,
	}
	var result []DocumentoRelacionado
	err := c.doGet(ctx, "/api-de-dados/despesas/documentos-relacionados", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarDocumentosPorFavorecido(ctx context.Context, filtro DespesaDocumentosPorFavorecidoQueryParams) ([]interface{}, error) {
	params := map[string]string{
		"codigoPessoa": filtro.CodigoPessoa,
		"fase":         filtro.Fase,
		"ano":          filtro.Ano,
		"pagina":       strconv.Itoa(filtro.Pagina),
	}
	if filtro.UG != "" {
		params["ug"] = filtro.UG
	}
	if filtro.Gestao != "" {
		params["gestao"] = filtro.Gestao
	}
	if filtro.OrdenacaoResultado != "" {
		params["ordenacaoResultado"] = filtro.OrdenacaoResultado
	}
	var result []interface{}
	err := c.doGet(ctx, "/api-de-dados/despesas/documentos-por-favorecido", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarTiposTransferencia(ctx context.Context) ([]CodigoDescricao, error) {
	params := map[string]string{}
	var result []CodigoDescricao
	err := c.doGet(ctx, "/api-de-dados/despesas/tipo-transferencia", params, &result)
	return result, err
}
