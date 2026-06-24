package app

import (
	"github.com/gin-gonic/gin"
)

func NovoRoteador(app *App) *gin.Engine {
	roteador := gin.New()
	roteador.Use(gin.Logger(), gin.Recovery())

	if app.LeitorCSVHandler != nil {
		roteador.POST("/import", app.LeitorCSVHandler.Executar)
	}

	if app.AnaliseOrgaoPNCPHandler != nil {
		roteador.POST("/orgao/analise", app.AnaliseOrgaoPNCPHandler.AnaliseOrgaoPNCP)
		roteador.GET("/orgao/analise/batch/:jobId", app.AnaliseOrgaoPNCPHandler.BuscarResultadosBatch)
	}

	if app.AnalisarLigacaoPoliticaHandler != nil {
		roteador.POST("/busca/contexto", app.AnalisarLigacaoPoliticaHandler.Analisar)
	}

	if app.HandlerBuscaRelacoes != nil {
		roteador.POST("/busca/relacoes", app.HandlerBuscaRelacoes.Buscar)
	}

	if app.HandlerConsultaEntidade != nil {
		roteador.POST("/entidade", app.HandlerConsultaEntidade.Consultar)
	}

	if app.AnalisePublicacaoHandler != nil {
		roteador.POST("/publicacao/analise", app.AnalisePublicacaoHandler.AnalisePublicacao)
		roteador.GET("/publicacao/analise/batch/:jobId", app.AnalisePublicacaoHandler.BuscarResultadosBatch)
		roteador.GET("/ibge/municipios/:uf", app.ListarMunicipiosHandler.ListarMunicipios)
	}

	if app.ContasIrregularesHandler != nil {
		roteador.POST("/tcu/contas-irregulares", app.ContasIrregularesHandler.Buscar)
	}
	if app.FinsEleitoraisHandler != nil {
		roteador.POST("/tcu/fins-eleitorais", app.FinsEleitoraisHandler.Buscar)
	}
	if app.InabilitadosHandler != nil {
		roteador.POST("/tcu/inabilitados", app.InabilitadosHandler.Buscar)
	}
	if app.InidoneosHandler != nil {
		roteador.POST("/tcu/inidoneos", app.InidoneosHandler.Buscar)
	}

	if app.ListarCargosTSEHandler != nil {
		roteador.GET("/busca/cargos", app.ListarCargosTSEHandler.ListarCargos)
		roteador.GET("/busca/partidos", app.ListarPartidosHandler.ListarPartidos)
		roteador.POST("/busca/candidatos", app.BuscarCandidatosHandler.BuscarCandidatos)
		roteador.POST("/busca/doadores", app.BuscarDoadorHandler.BuscarDoador)
		roteador.POST("/busca/fornecedores", app.BuscarFornecedorHandler.BuscarFornecedor)
	}

	if app.BuscarDeputadosAtivosHandler != nil {
		roteador.GET("/deputados", app.BuscarDeputadosAtivosHandler.BuscarDeputadosAtivos)
	}
	if app.BuscarDetalhesDeputadoHandler != nil {
		roteador.GET("/deputados/:id/completo", app.BuscarDetalhesDeputadoHandler.BuscarDetalhesDeputado)
	}
	if app.BuscarDespesasDeputadoHandler != nil {
		roteador.GET("/deputados/:id/despesas", app.BuscarDespesasDeputadoHandler.BuscarDespesasDeputado)
	}
	if app.BuscarOrgaoAssociadoDeputadoHandler != nil {
		roteador.GET("/deputados/:id/orgaos", app.BuscarOrgaoAssociadoDeputadoHandler.BuscarOrgaoAssociadoDeputado)
	}

	if app.ListarSenadoresHandler != nil {
		roteador.GET("/senado/senadores", app.ListarSenadoresHandler.Listar)
	}
	if app.BuscarSenadorHandler != nil {
		roteador.GET("/senado/senadores/:codigo/completo", app.BuscarSenadorHandler.Buscar)
	}
	if app.ListarCargosSenadorHandler != nil {
		roteador.GET("/senado/senadores/:codigo/cargos", app.ListarCargosSenadorHandler.Listar)
	}
	if app.ListarComissoesSenadorHandler != nil {
		roteador.GET("/senado/senadores/:codigo/comissoes", app.ListarComissoesSenadorHandler.Listar)
	}
	if app.ListarMandatosHandler != nil {
		roteador.GET("/senado/senadores/:codigo/mandatos", app.ListarMandatosHandler.Listar)
	}
	if app.ListarOrcamentoHandler != nil {
		roteador.GET("/senado/orcamento", app.ListarOrcamentoHandler.Listar)
	}
	if app.ListarProcessosHandler != nil {
		roteador.GET("/senado/processors", app.ListarProcessosHandler.Listar)
	}
	if app.ListarProcessoAssuntosHandler != nil {
		roteador.GET("/senado/processo/assuntos", app.ListarProcessoAssuntosHandler.Listar)
	}
	if app.ListarProcessoEmendasHandler != nil {
		roteador.GET("/senado/processo/emendas", app.ListarProcessoEmendasHandler.Listar)
	}
	if app.BuscarProcessoHandler != nil {
		roteador.GET("/senado/processo/:id", app.BuscarProcessoHandler.Buscar)
	}
	if app.ListarVotacoesHandler != nil {
		roteador.GET("/senado/votacoes", app.ListarVotacoesHandler.Listar)
	}
	if app.ListarVotacoesComissaoHandler != nil {
		roteador.GET("/senado/votacoes/comissao/:sigla", app.ListarVotacoesComissaoHandler.Listar)
	}
	if app.ListarVotacoesComissaoParlamentarHandler != nil {
		roteador.GET("/senado/votacoes/parlamentar/:codigo", app.ListarVotacoesComissaoParlamentarHandler.Listar)
	}
	if app.ListarMateriaTramitacaoHandler != nil {
		roteador.GET("/senado/materia/tramitacao", app.ListarMateriaTramitacaoHandler.Listar)
	}
	if app.ListarAgendaDiaHandler != nil {
		roteador.GET("/senado/agenda/dia/:data", app.ListarAgendaDiaHandler.Listar)
	}
	if app.ListarAgendaMesHandler != nil {
		roteador.GET("/senado/agenda/mes/:data", app.ListarAgendaMesHandler.Listar)
	}
	if app.BuscarEncontroHandler != nil {
		roteador.GET("/senado/encontro/:codigo", app.BuscarEncontroHandler.Buscar)
	}
	if app.ListarTodasComissoesHandler != nil {
		roteador.GET("/senado/comissoes", app.ListarTodasComissoesHandler.Listar)
	}
	if app.BuscarComissaoHandler != nil {
		roteador.GET("/senado/comissoes/:codigo", app.BuscarComissaoHandler.Buscar)
	}

	if app.ListarEstadosHandler != nil {
		roteador.GET("/ibge/estados", app.ListarEstadosHandler.ListarEstados)
		roteador.GET("/estado/:uf/dados-completos", app.BuscarDadosEstadoHandler.BuscarDadosEstado)
		roteador.GET("/estado/:uf/basico", app.BuscarBasicoEstadoHandler.BuscarBasicoEstado)
		roteador.GET("/estado/:uf/candidatos", app.BuscarCandidatosEstadoHandler.BuscarCandidatosEstado)
		roteador.GET("/estado/:uf/deputados", app.BuscarDeputadosEstadoHandler.BuscarDeputadosEstado)
		roteador.GET("/estado/:uf/senadores", app.BuscarSenadoresEstadoHandler.BuscarSenadoresEstado)
		roteador.GET("/ibge/municipios-populacao/:uf", app.BuscarMunicipiosPopulacaoHandler.BuscarMunicipiosPopulacao)
		roteador.GET("/estado/:uf/recursos-federais", app.BuscarRecursosFederaisCompletoHandler.Buscar)
	}

	if app.BuscarLicitacoesUFHandler != nil {
		roteador.GET("/estado/:uf/licitacoes", app.BuscarLicitacoesUFHandler.Buscar)
		roteador.GET("/estado/:uf/licitacoes/municipio/:codigo", app.BuscarLicitacoesUFHandler.Buscar)
	}

	if app.WSHub != nil {
		roteador.GET("/ws", app.WSHub.Handle)
	}

	if app.BuscarSIAPEHandler != nil {
		roteador.GET("/portal-transparencia/orgaos/siape", app.BuscarSIAPEHandler.BuscarSIAPE)
		roteador.GET("/portal-transparencia/orgaos/siafi", app.BuscarSIAFIHandler.BuscarSIAFI)
	}

	if app.BuscarFisicaHandler != nil {
		roteador.GET("/portal-transparencia/pessoas/fisica", app.BuscarFisicaHandler.BuscarFisica)
		roteador.GET("/portal-transparencia/pessoas/juridica", app.BuscarJuridicaHandler.BuscarJuridica)
	}

	if app.BuscarCartoesHandler != nil {
		roteador.GET("/portal-transparencia/cartoes", app.BuscarCartoesHandler.Buscar)
	}

	if app.BuscarServidoresHandler != nil {
		roteador.GET("/portal-transparencia/servidores", app.BuscarServidoresHandler.Buscar)
		roteador.GET("/portal-transparencia/servidores/:id", app.BuscarServidorPorIDHandler.BuscarPorID)
		roteador.GET("/portal-transparencia/servidores/remuneracao", app.BuscarRemuneracaoHandler.BuscarRemuneracao)
		roteador.GET("/portal-transparencia/servidores/por-orgao", app.BuscarServidoresPorOrgaoHandler.BuscarPorOrgao)
		roteador.GET("/portal-transparencia/servidores/funcoes-e-cargos", app.BuscarFuncoesECargosHandler.BuscarFuncoesECargos)
		roteador.GET("/portal-transparencia/servidores/peps", app.BuscarPEPsHandler.BuscarPEPs)
	}

	if app.BuscarTiposTransferenciaHandler != nil {
		roteador.GET("/portal-transparencia/despesas/tipo-transferencia", app.BuscarTiposTransferenciaHandler.BuscarTiposTransferencia)
		roteador.GET("/portal-transparencia/despesas/recursos-recebidos", app.BuscarRecursosRecebidosHandler.BuscarRecursosRecebidos)
		roteador.GET("/portal-transparencia/despesas/por-orgao", app.BuscarDespesasPorOrgaoHandler.BuscarPorOrgao)
		roteador.GET("/portal-transparencia/despesas/por-funcional-programatica", app.BuscarPorFuncionalProgramaticaHandler.BuscarPorFuncionalProgramatica)
		roteador.GET("/portal-transparencia/despesas/por-funcional-programatica/movimentacao-liquida", app.BuscarMovimentacaoLiquidaHandler.BuscarMovimentacaoLiquida)
		roteador.GET("/portal-transparencia/despesas/plano-orcamentario", app.BuscarPlanoOrcamentarioHandler.BuscarPlanoOrcamentario)
		roteador.GET("/portal-transparencia/despesas/itens-de-empenho", app.BuscarItensEmpenhoHandler.BuscarItensEmpenho)
		roteador.GET("/portal-transparencia/despesas/itens-de-empenho/historico", app.BuscarHistoricoItemEmpenhoHandler.BuscarHistoricoEmpenho)
		roteador.GET("/portal-transparencia/despesas/funcional-programatica/subfuncoes", app.BuscarSubfuncoesHandler.BuscarSubfuncoes)
		roteador.GET("/portal-transparencia/despesas/funcional-programatica/programs", app.BuscarProgramasHandler.BuscarProgramas)
		roteador.GET("/portal-transparencia/despesas/funcional-programatica/listar", app.ListarFuncionalProgramaticaHandler.ListarFuncionalProgramatica)
		roteador.GET("/portal-transparencia/despesas/funcional-programatica/funcoes", app.BuscarFuncoesHandler.BuscarFuncoes)
		roteador.GET("/portal-transparencia/despesas/funcional-programatica/acoes", app.BuscarAcoesHandler.BuscarAcoes)
		roteador.GET("/portal-transparencia/despesas/favorecidos-finais-por-documento", app.BuscarFavorecidosFinaisHandler.BuscarFavorecidosFinaisPorDocumento)
		roteador.GET("/portal-transparencia/despesas/empenhos-impactados", app.BuscarEmpenhosImpactadosHandler.BuscarEmpenhosImpactados)
		roteador.GET("/portal-transparencia/despesas/documentos", app.BuscarDocumentosHandler.BuscarDocumentos)
		roteador.GET("/portal-transparencia/despesas/documentos/:codigo", app.BuscarDocumentoPorCodigoHandler.BuscarDocumentoPorCodigo)
		roteador.GET("/portal-transparencia/despesas/documentos-relacionados", app.BuscarDocumentosRelacionadosHandler.BuscarDocumentosRelacionados)
		roteador.GET("/portal-transparencia/despesas/documentos-por-favorecido", app.BuscarDocumentosPorFavorecidoHandler.BuscarDocumentosPorFavorecido)
	}

	if app.BuscarEmendasHandler != nil {
		roteador.GET("/portal-transparencia/emendas", app.BuscarEmendasHandler.Buscar)
		roteador.GET("/portal-transparencia/emendas/documentos/:codigo", app.BuscarDocumentosEmendaHandler.BuscarDocumentos)
	}

	return roteador
}
