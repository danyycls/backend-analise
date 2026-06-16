package portaltransparencia

import (
	"context"
	"strconv"
)

func (c *PortalTransparenciaClient) ListarServidores(ctx context.Context, filtro ServidorQueryParams) ([]CadastroServidor, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.CPF != "" {
		params["cpf"] = filtro.CPF
	}
	if filtro.Nome != "" {
		params["nome"] = filtro.Nome
	}
	if filtro.OrgaoServidorLotacao != "" {
		params["orgaoServidorLotacao"] = filtro.OrgaoServidorLotacao
	}
	if filtro.OrgaoServidorExercicio != "" {
		params["orgaoServidorExercicio"] = filtro.OrgaoServidorExercicio
	}
	if filtro.SituacaoServidor != "" {
		params["situacaoServidor"] = filtro.SituacaoServidor
	}
	if filtro.TipoServidor != "" {
		params["tipoServidor"] = filtro.TipoServidor
	}
	if filtro.CodigoFuncaoCargo != "" {
		params["codigoFuncaoCargo"] = filtro.CodigoFuncaoCargo
	}
	var result []CadastroServidor
	err := c.doGet(ctx, "/api-de-dados/servidores", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) BuscarServidorPorID(ctx context.Context, id int) (*CadastroServidor, error) {
	params := map[string]string{}
	var result CadastroServidor
	err := c.doGet(ctx, "/api-de-dados/servidores/"+strconv.Itoa(id), params, &result)
	return &result, err
}

func (c *PortalTransparenciaClient) ListarRemuneracaoServidores(ctx context.Context, filtro ServidorRemuneracaoQueryParams) ([]ServidorRemuneracao, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
		"mesAno": filtro.MesAno,
	}
	if filtro.CPF != "" {
		params["cpf"] = filtro.CPF
	}
	if filtro.IDservidorAposentadoPensionista != "" {
		params["id"] = filtro.IDservidorAposentadoPensionista
	}
	var result []ServidorRemuneracao
	err := c.doGet(ctx, "/api-de-dados/servidores/remuneracao", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarServidoresPorOrgao(ctx context.Context, filtro ServidorPorOrgaoQueryParams) ([]ServidorPorOrgao, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.OrgaoLotacao != "" {
		params["orgaoLotacao"] = filtro.OrgaoLotacao
	}
	if filtro.OrgaoExercicio != "" {
		params["orgaoExercicio"] = filtro.OrgaoExercicio
	}
	if filtro.TipoServidor != "" {
		params["tipoServidor"] = filtro.TipoServidor
	}
	if filtro.TipoVinculo != "" {
		params["tipoVinculo"] = filtro.TipoVinculo
	}
	if filtro.Licenca != "" {
		params["licenca"] = filtro.Licenca
	}
	var result []ServidorPorOrgao
	err := c.doGet(ctx, "/api-de-dados/servidores/por-orgao", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarFuncoesECargos(ctx context.Context, filtro FuncaoCargoQueryParams) ([]FuncaoServidor, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.CodigoFuncaoCargo != "" {
		params["codigoFuncaoCargo"] = filtro.CodigoFuncaoCargo
	}
	if filtro.DescricaoFuncaoCargo != "" {
		params["descricaoFuncaoCargo"] = filtro.DescricaoFuncaoCargo
	}
	var result []FuncaoServidor
	err := c.doGet(ctx, "/api-de-dados/servidores/funcoes-e-cargos", params, &result)
	return result, err
}

func (c *PortalTransparenciaClient) ListarPEPs(ctx context.Context, filtro PEPQueryParams) ([]PEP, error) {
	params := map[string]string{
		"pagina": strconv.Itoa(filtro.Pagina),
	}
	if filtro.CPF != "" {
		params["cpf"] = filtro.CPF
	}
	if filtro.Nome != "" {
		params["nome"] = filtro.Nome
	}
	if filtro.DescricaoFuncao != "" {
		params["descricaoFuncao"] = filtro.DescricaoFuncao
	}
	if filtro.OrgaoServidorLotacao != "" {
		params["orgaoServidorLotacao"] = filtro.OrgaoServidorLotacao
	}
	if filtro.DataInicioExercicioDe != "" {
		params["dataInicioExercicioDe"] = filtro.DataInicioExercicioDe
	}
	if filtro.DataInicioExercicioAte != "" {
		params["dataInicioExercicioAte"] = filtro.DataInicioExercicioAte
	}
	if filtro.DataFimExercicioDe != "" {
		params["dataFimExercicioDe"] = filtro.DataFimExercicioDe
	}
	if filtro.DataFimExercicioAte != "" {
		params["dataFimExercicioAte"] = filtro.DataFimExercicioAte
	}
	var result []PEP
	err := c.doGet(ctx, "/api-de-dados/peps", params, &result)
	return result, err
}
