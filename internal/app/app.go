package app

import (
	"context"
	"os"

	feedback "github.com/danyele/podp/internal/feedback"
	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/redis"

	handlerEstadual "github.com/danyele/podp/internal/esferas-brasileiras/estadual/handler"
	handlerDeputados "github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/handler"
	handlerPNCP "github.com/danyele/podp/internal/esferas-brasileiras/federal/pncp/handler"
	handlerPortalCartoes "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/cartoes/handler"
	handlerPortalDespesas "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/despesas/handler"
	handlerPortalEmendas "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/emendas/handler"
	handlerPortalOrgaos "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/orgaos/handler"
	handlerPortalPessoas "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/pessoas/handler"
	handlerPortalServidores "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/servidores/handler"
	handlerSenadores "github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/handler"
	handlerTCU "github.com/danyele/podp/internal/esferas-brasileiras/federal/tcu/handler"
	tseHandler "github.com/danyele/podp/internal/esferas-brasileiras/tse/handler"
	importacaoHandler "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/handler"
	importacaoRepositorios "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/repositorios"
	importacaoService "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/service"
	importacaoUseCase "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/usecase"
	repositorios "github.com/danyele/podp/internal/shared/repositorios"
	tseUseCase "github.com/danyele/podp/internal/esferas-brasileiras/tse/usecase"
	handlerLigacao "github.com/danyele/podp/internal/ligacao-politica/handler"

	anomaliaHandler "github.com/danyele/podp/internal/worker/anomalia/handler"
	anomaliaUseCase "github.com/danyele/podp/internal/worker/anomalia/usecase"

	clientDeputados "github.com/danyele/podp/internal/shared/clients/deputados"
	clientPNCP "github.com/danyele/podp/internal/shared/clients/pncp"
	clientPortal "github.com/danyele/podp/internal/shared/clients/portaltransparencia"
	clientSenadores "github.com/danyele/podp/internal/shared/clients/senado"
	clientTCU "github.com/danyele/podp/internal/shared/clients/tcu"

	"github.com/danyele/podp/internal/shared/clients/ibge"
	"github.com/danyele/podp/internal/shared/clients/opencnpj"
	"github.com/danyele/podp/internal/shared/clients/siconfi"
	"github.com/danyele/podp/internal/shared/mongodb"

	usecaseEstadual "github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase"
	dadosfinanceiros "github.com/danyele/podp/internal/esferas-brasileiras/estadual/usecase/dadosfinanceiros"
	usecaseDeputados "github.com/danyele/podp/internal/esferas-brasileiras/federal/deputados/usecase"
	usecasePNCP "github.com/danyele/podp/internal/esferas-brasileiras/federal/pncp/usecase"
	usecasePortalCartoes "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/cartoes/usecase"
	usecasePortalDespesas "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/despesas/usecase"
	usecasePortalEmendas "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/emendas/usecase"
	usecasePortalOrgaos "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/orgaos/usecase"
	usecasePortalPessoas "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/pessoas/usecase"
	usecasePortalServidores "github.com/danyele/podp/internal/esferas-brasileiras/federal/portaltransparencia/servidores/usecase"
	usecaseSenadores "github.com/danyele/podp/internal/esferas-brasileiras/federal/senadores/usecase"
	usecaseTCU "github.com/danyele/podp/internal/esferas-brasileiras/federal/tcu/usecase"
	usecaseMunicipal "github.com/danyele/podp/internal/esferas-brasileiras/municipal/usecase"

	usecaseLigacao "github.com/danyele/podp/internal/ligacao-politica/usecase"
	"github.com/danyele/podp/internal/shared/services"
	"github.com/danyele/podp/internal/stream"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB     database.DB
	PgPool *pgxpool.Pool

	PNCPClient      *clientPNCP.PNCPClient
	OpenCNPJClient  *opencnpj.OpenCNPJClient
	IBGEClient      *ibge.IBGEClient
	DeputadosClient *clientDeputados.DeputadosClient
	SenadoClient    *clientSenadores.SenadoClient
	TCUClient       *clientTCU.TCUClient
	SICONFIClient   *siconfi.SICONFIClient
	PortalClient    *clientPortal.PortalTransparenciaClient
	RedisCache      *redis.RedisCache
	FeedbackHandler *feedback.Handler

	LeitorCSVService *importacaoService.LeitorCSVService
	LeitorCSVUseCase importacaoUseCase.ImportarCSVUseCase
	LeitorCSVHandler *importacaoHandler.LeitorCSVHandler

	AnalisarLigacaoPoliticaUseCase *usecaseLigacao.AnalisarLigacaoPoliticaUseCase
	AnalisarLigacaoPoliticaHandler *handlerLigacao.AnalisarLigacaoPoliticaHandler

	AnaliseOrgaoPNCPHandler   *handlerPNCP.AnaliseOrgaoPNCPHandler
	AnalisePublicacaoHandler  *handlerPNCP.AnalisePublicacaoHandler
	BuscarLicitacoesUFHandler *handlerPNCP.BuscarLicitacoesUFHandler
	ListarMunicipiosHandler   *handlerPNCP.ListarMunicipiosHandler

	HandlerBuscaRelacoes    *tseHandler.BuscaRelacoesHandler
	HandlerConsultaEntidade *tseHandler.ConsultaEntidadeHandler

	ListarCargosTSEHandler  *tseHandler.ListarCargosHandler
	ListarPartidosHandler   *tseHandler.ListarPartidosHandler
	BuscarCandidatosHandler *tseHandler.BuscarCandidatosHandler
	BuscarDoadorHandler     *tseHandler.BuscarDoadorHandler
	BuscarFornecedorHandler *tseHandler.BuscarFornecedorHandler

	ContasIrregularesHandler *handlerTCU.ContasIrregularesHandler
	FinsEleitoraisHandler    *handlerTCU.FinsEleitoraisHandler
	InabilitadosHandler      *handlerTCU.InabilitadosHandler
	InidoneosHandler         *handlerTCU.InidoneosHandler

	BuscarDeputadosAtivosHandler        *handlerDeputados.EsferaFederalBuscarDeputadosAtivosHandler
	BuscarDetalhesDeputadoHandler       *handlerDeputados.EsferaFederalBuscarDetalhesDeputadoHandler
	BuscarDespesasDeputadoHandler       *handlerDeputados.EsferaFederalBuscarDespesasDeputadoHandler
	BuscarOrgaoAssociadoDeputadoHandler *handlerDeputados.EsferaFederalBuscarOrgaoAssociadoDeputadoHandler

	DepListarPartidosHandler        *handlerDeputados.EsferaFederalListarPartidosHandler
	DepBuscarPartidoHandler         *handlerDeputados.EsferaFederalBuscarPartidoHandler
	DepListarMembrosPartidoHandler  *handlerDeputados.EsferaFederalListarMembrosPartidoHandler
	DepListarProposicoesHandler     *handlerDeputados.ListarProposicoesHandler
	DepBuscarProposicaoHandler      *handlerDeputados.BuscarProposicaoHandler
	DepListarTramitacoesHandler     *handlerDeputados.ListarTramitacoesHandler
	DepListarAutoresHandler         *handlerDeputados.ListarAutoresHandler
	DepListarTemasHandler           *handlerDeputados.ListarTemasHandler
	DepListarRelacionadasHandler    *handlerDeputados.ListarRelacionadasHandler
	DepListarEventosHandler         *handlerDeputados.EsferaFederalListarEventosHandler
	DepBuscarEventoHandler          *handlerDeputados.EsferaFederalBuscarEventoHandler
	DepListarOrgaosHandler          *handlerDeputados.EsferaFederalListarOrgaosHandler
	DepBuscarOrgaoHandler           *handlerDeputados.EsferaFederalBuscarOrgaoHandler
	DepListarMembrosOrgaoHandler    *handlerDeputados.EsferaFederalListarMembrosOrgaoHandler
	DepListarBlocosHandler          *handlerDeputados.EsferaFederalListarBlocosHandler
	DepBuscarBlocoHandler           *handlerDeputados.EsferaFederalBuscarBlocoHandler
	DepListarPartidosDoBlocoHandler *handlerDeputados.EsferaFederalListarPartidosDoBlocoHandler
	DepListarFrentesHandler         *handlerDeputados.EsferaFederalListarFrentesHandler
	DepBuscarFrenteHandler          *handlerDeputados.EsferaFederalBuscarFrenteHandler
	DepListarMembrosFrenteHandler   *handlerDeputados.EsferaFederalListarMembrosFrenteHandler
	DepListarGruposHandler          *handlerDeputados.EsferaFederalListarGruposHandler
	DepBuscarGrupoHandler           *handlerDeputados.EsferaFederalBuscarGrupoHandler
	DepListarLegislaturasHandler    *handlerDeputados.EsferaFederalListarLegislaturasHandler
	DepBuscarLegislaturaHandler     *handlerDeputados.EsferaFederalBuscarLegislaturaHandler
	DepListarVotacoesHandler        *handlerDeputados.EsferaFederalListarVotacoesHandler
	DepBuscarVotacaoHandler         *handlerDeputados.EsferaFederalBuscarVotacaoHandler
	DepListarVotosHandler           *handlerDeputados.EsferaFederalListarVotosHandler

	BaixarDocumentoEmendaHandler *handlerSenadores.BaixarDocumentoEmendaHandler

	ListarSenadoresHandler                   *handlerSenadores.ListarSenadoresHandler
	BuscarSenadorHandler                     *handlerSenadores.BuscarSenadorHandler
	ListarCargosSenadorHandler               *handlerSenadores.ListarCargosHandler
	ListarComissoesSenadorHandler            *handlerSenadores.ListarComissoesHandler
	ListarMandatosHandler                    *handlerSenadores.ListarMandatosHandler
	ListarOrcamentoHandler                   *handlerSenadores.ListarOrcamentoHandler
	ListarProcessosHandler                   *handlerSenadores.ListarProcessosHandler
	ListarProcessoAssuntosHandler            *handlerSenadores.ListarProcessoAssuntosHandler
	ListarProcessoEmendasHandler             *handlerSenadores.ListarProcessoEmendasHandler
	BuscarProcessoHandler                    *handlerSenadores.BuscarProcessoHandler
	ListarVotacoesHandler                    *handlerSenadores.ListarVotacoesHandler
	ListarVotacoesComissaoHandler            *handlerSenadores.ListarVotacoesComissaoHandler
	ListarVotacoesComissaoParlamentarHandler *handlerSenadores.ListarVotacoesComissaoParlamentarHandler
	ListarMateriaTramitacaoHandler           *handlerSenadores.ListarMateriaTramitacaoHandler
	ListarAgendaDiaHandler                   *handlerSenadores.ListarAgendaDiaHandler
	ListarAgendaMesHandler                   *handlerSenadores.ListarAgendaMesHandler
	BuscarEncontroHandler                    *handlerSenadores.BuscarEncontroHandler
	ListarTodasComissoesHandler              *handlerSenadores.ListarTodasComissoesHandler
	BuscarComissaoHandler                    *handlerSenadores.BuscarComissaoHandler

	ListarEstadosHandler                  *handlerEstadual.EsferaEstadualListarEstadosHandler
	BuscarDadosEstadoHandler              *handlerEstadual.EsferaEstadualBuscarDadosCompletosEstadoHandler
	BuscarBasicoEstadoHandler             *handlerEstadual.EsferaEstadualBuscarDadosBasicosEstadoHandler
	BuscarCandidatosEstadoHandler         *handlerEstadual.EsferaEstadualBuscarCandidatosHandler
	BuscarDeputadosEstadoHandler          *handlerEstadual.EsferaEstadualBuscarDeputadosHandler
	BuscarSenadoresEstadoHandler          *handlerEstadual.EsferaEstadualBuscarSenadoresHandler
	BuscarMunicipiosPopulacaoHandler      *handlerEstadual.EsferaEstadualBuscarMunicipiosPopulacaoHandler
	BuscarRecursosFederaisCompletoHandler *handlerEstadual.EsferaEstadualBuscarRecursosFederaisCompletoHandler
	WSHub                                 *stream.Hub

	BuscarSIAPEHandler    *handlerPortalOrgaos.BuscarSIAPEHandler
	BuscarSIAFIHandler    *handlerPortalOrgaos.BuscarSIAFIHandler
	BuscarFisicaHandler   *handlerPortalPessoas.BuscarFisicaHandler
	BuscarJuridicaHandler *handlerPortalPessoas.BuscarJuridicaHandler
	BuscarCartoesHandler  *handlerPortalCartoes.BuscarCartoesHandler

	BuscarServidoresHandler         *handlerPortalServidores.BuscarServidoresHandler
	BuscarServidorPorIDHandler      *handlerPortalServidores.BuscarServidorPorIDHandler
	BuscarRemuneracaoHandler        *handlerPortalServidores.BuscarRemuneracaoHandler
	BuscarServidoresPorOrgaoHandler *handlerPortalServidores.BuscarServidoresPorOrgaoHandler
	BuscarFuncoesECargosHandler     *handlerPortalServidores.BuscarFuncoesECargosHandler
	BuscarPEPsHandler               *handlerPortalServidores.BuscarPEPsHandler

	BuscarTiposTransferenciaHandler       *handlerPortalDespesas.BuscarTiposTransferenciaHandler
	BuscarRecursosRecebidosHandler        *handlerPortalDespesas.BuscarRecursosRecebidosHandler
	BuscarDespesasPorOrgaoHandler         *handlerPortalDespesas.BuscarDespesasPorOrgaoHandler
	BuscarPorFuncionalProgramaticaHandler *handlerPortalDespesas.BuscarPorFuncionalProgramaticaHandler
	BuscarMovimentacaoLiquidaHandler      *handlerPortalDespesas.BuscarMovimentacaoLiquidaHandler
	BuscarPlanoOrcamentarioHandler        *handlerPortalDespesas.BuscarPlanoOrcamentarioHandler
	BuscarItensEmpenhoHandler             *handlerPortalDespesas.BuscarItensEmpenhoHandler
	BuscarHistoricoItemEmpenhoHandler     *handlerPortalDespesas.BuscarHistoricoItemEmpenhoHandler
	BuscarSubfuncoesHandler               *handlerPortalDespesas.BuscarSubfuncoesHandler
	BuscarProgramasHandler                *handlerPortalDespesas.BuscarProgramasHandler
	ListarFuncionalProgramaticaHandler    *handlerPortalDespesas.ListarFuncionalProgramaticaHandler
	BuscarFuncoesHandler                  *handlerPortalDespesas.BuscarFuncoesHandler
	BuscarAcoesHandler                    *handlerPortalDespesas.BuscarAcoesHandler
	BuscarFavorecidosFinaisHandler        *handlerPortalDespesas.BuscarFavorecidosFinaisPorDocumentoHandler
	BuscarEmpenhosImpactadosHandler       *handlerPortalDespesas.BuscarEmpenhosImpactadosHandler
	BuscarDocumentosHandler               *handlerPortalDespesas.BuscarDocumentosHandler
	BuscarDocumentoPorCodigoHandler       *handlerPortalDespesas.BuscarDocumentoPorCodigoHandler
	BuscarDocumentosRelacionadosHandler   *handlerPortalDespesas.BuscarDocumentosRelacionadosHandler
	BuscarDocumentosPorFavorecidoHandler  *handlerPortalDespesas.BuscarDocumentosPorFavorecidoHandler

	BuscarEmendasHandler          *handlerPortalEmendas.BuscarEmendasHandler
	BuscarDocumentosEmendaHandler *handlerPortalEmendas.BuscarDocumentosEmendaHandler

	MongoClient             mongodb.Client
	GerarDescricaoSvc       services.GerarDescricaoVinculoService
	AnomaliaWorkerHandler   *anomaliaHandler.AnomaliaWorkerHandler
	AnomaliaConsultaHandler *anomaliaHandler.AnomaliaConsultaHandler

	ConvenioHandler *handlerPNCP.ConsultaConvenioHandler
}

func NovoApp(db database.DB, diretorioCSV string) *App {
	log := logger.New("App")

	pncpClient := clientPNCP.NovoPNCPClient(os.Getenv("PNCP_BASE_URL"))
	opencnpjClient := opencnpj.NovoOpenCNPJClient(os.Getenv("OPENCNPJ_BASE_URL"))
	ibgeClient := ibge.NovoIBGEClient(os.Getenv("IBGE_BASE_URL"), os.Getenv("IBGE_AGREGADOS_BASE_URL"))
	tcuClient := clientTCU.NovoTCUClient(os.Getenv("TCU_BASE_URL"))
	deputadosClient := clientDeputados.NovoDeputadosClient(os.Getenv("DEPUTADOS_BASE_URL"))
	senadoClient := clientSenadores.NovoSenadoClient(os.Getenv("SENADO_BASE_URL"))
	siconfiClient := siconfi.NovoSICONFIClient(os.Getenv("SICONFI_BASE_URL"))
	redisCache := redis.NovoRedisCache()
	if err := redisCache.Ping(context.Background()); err != nil {
		log.Fatal("redis indisponivel", "erro", err)
	}

	portalClient := clientPortal.NovoPortalTransparenciaClient(
		os.Getenv("PORTAL_TRANSPARENCIA_API_KEY"),
		os.Getenv("PORTAL_TRANSPARENCIA_BASE_URL"),
	)

	leitorCSVService := importacaoService.NovoLeitorCSVService(diretorioCSV)

	pgPool := db.Pool()

	pgPoolLeitura, err := importacaoRepositorios.NovoPool(context.Background(), importacaoRepositorios.ConfigLeituraFromEnv())
	if err != nil {
		log.Fatal("criar pool de leitura para importacao", "erro", err)
	}

	pncpRepo := repositorios.NovoPNCPRepository(db)

	leitorCSVUseCase := importacaoUseCase.NovoImportarCSVUseCase(pgPool, pgPoolLeitura, leitorCSVService)
	leitorCSVHandler := importacaoHandler.NovoLeitorCSVHandler(leitorCSVUseCase)

	normalizarSvc := services.NovoNormalizarDocumentosService()
	buscarSvc := services.NovoBuscarLigacaoPoliticaTSEService()
	verificarTcuSvc := services.NovoVerificarSancoesTCUService(tcuClient)
	verificarServSvc := services.NovoVerificarServidorPublicoService(portalClient)
	verificarPessoaPublicaSvc := services.NovoVerificarPessoaPublicaService(portalClient)
	gerarDescricaoSvc := services.NovoGerarDescricaoVinculoService()

	casoUsoLigacao := usecaseLigacao.NovoAnalisarLigacaoPoliticaUseCase(db, normalizarSvc, buscarSvc, verificarTcuSvc, verificarServSvc)
	handlerLigacao := handlerLigacao.NovoAnalisarLigacaoPoliticaHandler(casoUsoLigacao, redisCache)

	mongoClient, errMongo := mongodb.NovoMongoClient(context.Background())
	if errMongo != nil {
		log.Fatal("mongodb indisponivel", "erro", errMongo)
	}

	relacoesHandlerBusca := tseHandler.NovoBuscarRelacoesHandler(tseUseCase.NovoBuscarRelacoesUseCase(db))
	relacoesHandlerEntidade := tseHandler.NovoConsultarEntidadeHandler(tseUseCase.NovoConsultarEntidadeUseCase(db))

	buscaCandidatosUC := tseUseCase.NovoBuscarCandidatosUseCase(db)
	buscarDoadorUC := tseUseCase.NovoBuscarDoadorUseCase(db)
	buscarFornecedorUC := tseUseCase.NovoBuscarFornecedorUseCase(db)

	pncpAnaliseOrgaoHandler := handlerPNCP.NovoAnaliseOrgaoPNCPHandler(
		usecasePNCP.NovoConsultaCNPJOrgaoPNCPUseCase(pncpClient, opencnpjClient, redisCache, pncpRepo),
		redisCache,
	)
	pncpAnalisePubHandler := handlerPNCP.NovoAnalisePublicacaoHandler(
		usecasePNCP.NovoConsultaPublicacaoPNCPUseCase(pncpClient, opencnpjClient, redisCache, pncpRepo),
		redisCache,
	)

	anomaliaWorkerUC := anomaliaUseCase.NovoAnaliseAnomaliaWorkerUseCase(
		normalizarSvc,
		buscarSvc,
		verificarTcuSvc,
		verificarServSvc,
		verificarPessoaPublicaSvc,
		gerarDescricaoSvc,
		db,
		mongoClient,
	)
	anomaliaWorkerHandler := anomaliaHandler.NovoAnomaliaWorkerHandler(anomaliaWorkerUC)

	anomaliaConsutaUC := anomaliaUseCase.NovoAnomaliaConsultaUseCase(mongoClient)
	anomaliaConsultaHandler := anomaliaHandler.NovoAnomaliaConsultaHandler(anomaliaConsutaUC)

	convenioUC := usecasePNCP.NovoConsultaConvenioUseCase(db)
	convenioHandler := handlerPNCP.NovoConsultaConvenioHandler(convenioUC)

	contasIrregularesUC := usecaseTCU.NovoContasIrregularesUseCase(tcuClient)
	finsEleitoraisUC := usecaseTCU.NovoFinsEleitoraisUseCase(tcuClient)
	inabilitadosUC := usecaseTCU.NovoInabilitadosUseCase(tcuClient)
	inidoneosUC := usecaseTCU.NovoInidoneosUseCase(tcuClient)

	depAtivosUC := usecaseDeputados.NovoEsferaFederalBuscarDeputadosAtivosUseCase(deputadosClient)
	depDetalhesUC := usecaseDeputados.NovoEsferaFederalBuscarDetalhesDeputadoUseCase(deputadosClient)
	depDespesasUC := usecaseDeputados.NovoEsferaFederalBuscarDespesasDeputadoUseCase(deputadosClient)
	depOrgaoUC := usecaseDeputados.NovoEsferaFederalBuscarOrgaoAssociadoDeputadoUseCase(deputadosClient)

	listarPartidosUC := usecaseDeputados.NovoEsferaFederalListarPartidosUseCase(deputadosClient)
	buscarPartidoUC := usecaseDeputados.NovoEsferaFederalBuscarPartidoUseCase(deputadosClient)
	listarMembrosPartidoUC := usecaseDeputados.NovoEsferaFederalListarMembrosPartidoUseCase(deputadosClient)
	listarProposicoesUC := usecaseDeputados.NovoListarProposicoesUseCase(deputadosClient)
	buscarProposicaoUC := usecaseDeputados.NovoBuscarProposicaoUseCase(deputadosClient)
	listarTramitacoesUC := usecaseDeputados.NovoListarTramitacoesUseCase(deputadosClient)
	listarAutoresUC := usecaseDeputados.NovoListarAutoresUseCase(deputadosClient)
	listarTemasUC := usecaseDeputados.NovoListarTemasUseCase(deputadosClient)
	listarRelacionadasUC := usecaseDeputados.NovoListarRelacionadasUseCase(deputadosClient)
	listarEventosUC := usecaseDeputados.NovoEsferaFederalListarEventosUseCase(deputadosClient)
	buscarEventoUC := usecaseDeputados.NovoEsferaFederalBuscarEventoUseCase(deputadosClient)
	listarOrgaosUC := usecaseDeputados.NovoEsferaFederalListarOrgaosUseCase(deputadosClient)
	buscarOrgaoUC := usecaseDeputados.NovoEsferaFederalBuscarOrgaoUseCase(deputadosClient)
	listarMembrosOrgaoUC := usecaseDeputados.NovoEsferaFederalListarMembrosOrgaoUseCase(deputadosClient)
	listarBlocosUC := usecaseDeputados.NovoEsferaFederalListarBlocosUseCase(deputadosClient)
	buscarBlocoUC := usecaseDeputados.NovoEsferaFederalBuscarBlocoUseCase(deputadosClient)
	listarPartidosDoBlocoUC := usecaseDeputados.NovoEsferaFederalListarPartidosDoBlocoUseCase(deputadosClient)
	listarFrentesUC := usecaseDeputados.NovoEsferaFederalListarFrentesUseCase(deputadosClient)
	buscarFrenteUC := usecaseDeputados.NovoEsferaFederalBuscarFrenteUseCase(deputadosClient)
	listarMembrosFrenteUC := usecaseDeputados.NovoEsferaFederalListarMembrosFrenteUseCase(deputadosClient)
	listarGruposUC := usecaseDeputados.NovoEsferaFederalListarGruposUseCase(deputadosClient)
	buscarGrupoUC := usecaseDeputados.NovoEsferaFederalBuscarGrupoUseCase(deputadosClient)
	listarLegislaturasUC := usecaseDeputados.NovoEsferaFederalListarLegislaturasUseCase(deputadosClient)
	buscarLegislaturaUC := usecaseDeputados.NovoEsferaFederalBuscarLegislaturaUseCase(deputadosClient)
	depListarVotacoesUC := usecaseDeputados.NovoEsferaFederalListarVotacoesUseCase(deputadosClient)
	buscarVotacaoUC := usecaseDeputados.NovoEsferaFederalBuscarVotacaoUseCase(deputadosClient)
	listarVotosUC := usecaseDeputados.NovoEsferaFederalListarVotosUseCase(deputadosClient)

	baixarDocumentoEmendaUC := usecaseSenadores.NovoBaixarDocumentoEmendaUseCase(senadoClient)

	listarSenadoresUC := usecaseSenadores.NovoListarSenadoresUseCase(senadoClient)
	buscarSenadorUC := usecaseSenadores.NovoBuscarSenadorUseCase(senadoClient)
	listarCargosSenUC := usecaseSenadores.NovoListarCargosUseCase(senadoClient)
	listarComissoesSenUC := usecaseSenadores.NovoListarComissoesUseCase(senadoClient)
	listarMandatosUC := usecaseSenadores.NovoListarMandatosUseCase(senadoClient)
	listarOrcamentoUC := usecaseSenadores.NovoListarOrcamentoUseCase(senadoClient)
	listarProcessosUC := usecaseSenadores.NovoListarProcessosUseCase(senadoClient)
	listarProcessoAssuntosUC := usecaseSenadores.NovoListarProcessoAssuntosUseCase(senadoClient)
	listarProcessoEmendasUC := usecaseSenadores.NovoListarProcessoEmendasUseCase(senadoClient)
	buscarProcessoUC := usecaseSenadores.NovoBuscarProcessoUseCase(senadoClient)
	listarVotacoesUC := usecaseSenadores.NovoListarVotacoesUseCase(senadoClient)
	listarVotacoesComissaoUC := usecaseSenadores.NovoListarVotacoesComissaoUseCase(senadoClient)
	listarVotacoesComParlUC := usecaseSenadores.NovoListarVotacoesComissaoParlamentarUseCase(senadoClient)
	listarMateriaTramitacaoUC := usecaseSenadores.NovoListarMateriaTramitacaoUseCase(senadoClient)
	listarAgendaDiaUC := usecaseSenadores.NovoListarAgendaDiaUseCase(senadoClient)
	listarAgendaMesUC := usecaseSenadores.NovoListarAgendaMesUseCase(senadoClient)
	buscarEncontroUC := usecaseSenadores.NovoBuscarEncontroUseCase(senadoClient)
	listarTodasComissoesUC := usecaseSenadores.NovoListarTodasComissoesUseCase(senadoClient)
	buscarComissaoUC := usecaseSenadores.NovoBuscarComissaoUseCase(senadoClient)

	listarEstadosUC := usecaseEstadual.NovoEsferaEstadualListarEstadosUseCase(ibgeClient)
	dadosCompletosUC := usecaseEstadual.NovoEsferaEstadualBuscarDadosCompletosEstadoUseCase(db, ibgeClient, deputadosClient, senadoClient)
	basicoEstadoUC := usecaseEstadual.NovoEsferaEstadualBuscarDadosBasicosEstadoUseCase(ibgeClient)
	candidatosEstadoUC := usecaseEstadual.NovoEsferaEstadualBuscarCandidatosUseCase(db)
	deputadosEstadoUC := usecaseEstadual.NovoEsferaEstadualBuscarDeputadosUseCase(deputadosClient)
	senadoresEstadoUC := usecaseEstadual.NovoEsferaEstadualBuscarSenadoresUseCase(senadoClient)
	municipiosPopulacaoUC := usecaseEstadual.NovoEsferaEstadualBuscarMunicipiosPopulacaoUseCase(ibgeClient)

	baseFinanceiroUC := dadosfinanceiros.NovoBaseFinanceiroUseCase(siconfiClient, ibgeClient, redisCache)
	despesaPessoalUC := dadosfinanceiros.NovoEsferaEstadualBuscarDespesaPessoalUseCase(baseFinanceiroUC)
	despesaCategoriaUC := dadosfinanceiros.NovoEsferaEstadualBuscarDespesaCategoriaUseCase(baseFinanceiroUC)
	rreoUC := dadosfinanceiros.NovoEsferaEstadualBuscarRREOUseCase(baseFinanceiroUC)
	recursosFederaisUC := dadosfinanceiros.NovoEsferaEstadualBuscarRecursosFederaisUseCase(portalClient, redisCache)
	recursosFederaisCompletoUC := dadosfinanceiros.NovoEsferaEstadualBuscarRecursosFederaisCompletoUseCase(portalClient, redisCache)

	detalhesMunicipioUC := usecaseMunicipal.NovoEsferaMunicipalBuscarDetalhesUseCase(siconfiClient, pncpClient, redisCache, pncpRepo)

	siapeUC := usecasePortalOrgaos.NovoBuscarOrgaosSIAPEUseCase(portalClient)
	siafiUC := usecasePortalOrgaos.NovoBuscarOrgaosSIAFIUseCase(portalClient)
	fisicaUC := usecasePortalPessoas.NovoBuscarPessoasFisicasUseCase(portalClient)
	juridicaUC := usecasePortalPessoas.NovoBuscarPessoasJuridicasUseCase(portalClient)
	cartoesUC := usecasePortalCartoes.NovoBuscarCartoesUseCase(portalClient)

	servidoresUC := usecasePortalServidores.NovoBuscarServidoresUseCase(portalClient)
	servidorPorIDUC := usecasePortalServidores.NovoBuscarServidorPorIDUseCase(portalClient)
	remuneracaoUC := usecasePortalServidores.NovoBuscarRemuneracaoServidoresUseCase(portalClient)
	servidoresPorOrgaoUC := usecasePortalServidores.NovoBuscarServidoresPorOrgaoUseCase(portalClient)
	funcoesECargosUC := usecasePortalServidores.NovoBuscarFuncoesECargosUseCase(portalClient)
	pepsUC := usecasePortalServidores.NovoBuscarPEPsUseCase(portalClient)

	recursosRecebidosUC := usecasePortalDespesas.NovoBuscarRecursosRecebidosUseCase(portalClient)
	despesasPorOrgaoUC := usecasePortalDespesas.NovoBuscarDespesasPorOrgaoUseCase(portalClient)
	funcionalProgramaticaUC := usecasePortalDespesas.NovoBuscarDespesasPorFuncionalProgramaticaUseCase(portalClient)
	movLiquidaUC := usecasePortalDespesas.NovoBuscarMovimentacaoLiquidaUseCase(portalClient)
	planoOrcamentarioUC := usecasePortalDespesas.NovoBuscarPlanoOrcamentarioUseCase(portalClient)
	itensEmpenhoUC := usecasePortalDespesas.NovoBuscarItensEmpenhoUseCase(portalClient)
	historicoEmpenhoUC := usecasePortalDespesas.NovoBuscarHistoricoItemEmpenhoUseCase(portalClient)
	subfuncoesUC := usecasePortalDespesas.NovoBuscarSubfuncoesUseCase(portalClient)
	programasUC := usecasePortalDespesas.NovoBuscarProgramasUseCase(portalClient)
	listarFuncProgramaticaUC := usecasePortalDespesas.NovoListarFuncionalProgramaticaUseCase(portalClient)
	funcoesUC := usecasePortalDespesas.NovoBuscarFuncoesUseCase(portalClient)
	acoesUC := usecasePortalDespesas.NovoBuscarAcoesUseCase(portalClient)
	favorecidosUC := usecasePortalDespesas.NovoBuscarFavorecidosFinaisPorDocumentoUseCase(portalClient)
	empenhosImpactadosUC := usecasePortalDespesas.NovoBuscarEmpenhosImpactadosUseCase(portalClient)
	documentosUC := usecasePortalDespesas.NovoBuscarDocumentosUseCase(portalClient)
	documentoPorCodigoUC := usecasePortalDespesas.NovoBuscarDocumentoPorCodigoUseCase(portalClient)
	documentosRelacionadosUC := usecasePortalDespesas.NovoBuscarDocumentosRelacionadosUseCase(portalClient)
	documentosPorFavorecidoUC := usecasePortalDespesas.NovoBuscarDocumentosPorFavorecidoUseCase(portalClient)
	tiposTransferenciaUC := usecasePortalDespesas.NovoBuscarTiposTransferenciaUseCase(portalClient)

	emendasUC := usecasePortalEmendas.NovoBuscarEmendasUseCase(portalClient)
	documentosEmendaUC := usecasePortalEmendas.NovoBuscarDocumentosEmendaUseCase(portalClient)

	feedbackUsecase := &feedback.SaveFeedbackUsecase{Redis: redisCache}
	feedbackHandler := &feedback.Handler{Usecase: feedbackUsecase}

	return &App{
		DB:     db,
		PgPool: pgPool,

		PNCPClient:      pncpClient,
		OpenCNPJClient:  opencnpjClient,
		IBGEClient:      ibgeClient,
		DeputadosClient: deputadosClient,
		SenadoClient:    senadoClient,
		TCUClient:       tcuClient,
		SICONFIClient:   siconfiClient,
		PortalClient:    portalClient,
		RedisCache:      redisCache,
		FeedbackHandler: feedbackHandler,

		LeitorCSVService: leitorCSVService,
		LeitorCSVUseCase: leitorCSVUseCase,
		LeitorCSVHandler: leitorCSVHandler,

		AnalisarLigacaoPoliticaUseCase: casoUsoLigacao,
		AnalisarLigacaoPoliticaHandler: handlerLigacao,

		AnaliseOrgaoPNCPHandler:   pncpAnaliseOrgaoHandler,
		AnalisePublicacaoHandler:  pncpAnalisePubHandler,
		BuscarLicitacoesUFHandler: handlerPNCP.NovoBuscarLicitacoesUFHandler(usecasePNCP.NovoConsultaPublicacaoPNCPUseCase(pncpClient, opencnpjClient, redisCache, pncpRepo)),
		ListarMunicipiosHandler:   handlerPNCP.NovoListarMunicipiosHandler(ibgeClient),

		HandlerBuscaRelacoes:    relacoesHandlerBusca,
		HandlerConsultaEntidade: relacoesHandlerEntidade,

		ListarCargosTSEHandler:  tseHandler.NovoListarCargosHandler(buscaCandidatosUC),
		ListarPartidosHandler:   tseHandler.NovoListarPartidosHandler(buscaCandidatosUC),
		BuscarCandidatosHandler: tseHandler.NovoBuscarCandidatosHandler(buscaCandidatosUC),
		BuscarDoadorHandler:     tseHandler.NovoBuscarDoadorHandler(buscarDoadorUC),
		BuscarFornecedorHandler: tseHandler.NovoBuscarFornecedorHandler(buscarFornecedorUC),

		ContasIrregularesHandler: handlerTCU.NovoContasIrregularesHandler(contasIrregularesUC),
		FinsEleitoraisHandler:    handlerTCU.NovoFinsEleitoraisHandler(finsEleitoraisUC),
		InabilitadosHandler:      handlerTCU.NovoInabilitadosHandler(inabilitadosUC),
		InidoneosHandler:         handlerTCU.NovoInidoneosHandler(inidoneosUC),

		BuscarDeputadosAtivosHandler:        handlerDeputados.NovoEsferaFederalBuscarDeputadosAtivosHandler(depAtivosUC, redisCache),
		BuscarDetalhesDeputadoHandler:       handlerDeputados.NovoEsferaFederalBuscarDetalhesDeputadoHandler(depDetalhesUC),
		BuscarDespesasDeputadoHandler:       handlerDeputados.NovoEsferaFederalBuscarDespesasDeputadoHandler(depDespesasUC),
		BuscarOrgaoAssociadoDeputadoHandler: handlerDeputados.NovoEsferaFederalBuscarOrgaoAssociadoDeputadoHandler(depOrgaoUC),

		DepListarPartidosHandler:        handlerDeputados.NovoEsferaFederalListarPartidosHandler(listarPartidosUC),
		DepBuscarPartidoHandler:         handlerDeputados.NovoEsferaFederalBuscarPartidoHandler(buscarPartidoUC),
		DepListarMembrosPartidoHandler:  handlerDeputados.NovoEsferaFederalListarMembrosPartidoHandler(listarMembrosPartidoUC),
		DepListarProposicoesHandler:     handlerDeputados.NovoListarProposicoesHandler(listarProposicoesUC),
		DepBuscarProposicaoHandler:      handlerDeputados.NovoBuscarProposicaoHandler(buscarProposicaoUC),
		DepListarTramitacoesHandler:     handlerDeputados.NovoListarTramitacoesHandler(listarTramitacoesUC),
		DepListarAutoresHandler:         handlerDeputados.NovoListarAutoresHandler(listarAutoresUC),
		DepListarTemasHandler:           handlerDeputados.NovoListarTemasHandler(listarTemasUC),
		DepListarRelacionadasHandler:    handlerDeputados.NovoListarRelacionadasHandler(listarRelacionadasUC),
		DepListarEventosHandler:         handlerDeputados.NovoEsferaFederalListarEventosHandler(listarEventosUC),
		DepBuscarEventoHandler:          handlerDeputados.NovoEsferaFederalBuscarEventoHandler(buscarEventoUC),
		DepListarOrgaosHandler:          handlerDeputados.NovoEsferaFederalListarOrgaosHandler(listarOrgaosUC),
		DepBuscarOrgaoHandler:           handlerDeputados.NovoEsferaFederalBuscarOrgaoHandler(buscarOrgaoUC),
		DepListarMembrosOrgaoHandler:    handlerDeputados.NovoEsferaFederalListarMembrosOrgaoHandler(listarMembrosOrgaoUC),
		DepListarBlocosHandler:          handlerDeputados.NovoEsferaFederalListarBlocosHandler(listarBlocosUC),
		DepBuscarBlocoHandler:           handlerDeputados.NovoEsferaFederalBuscarBlocoHandler(buscarBlocoUC),
		DepListarPartidosDoBlocoHandler: handlerDeputados.NovoEsferaFederalListarPartidosDoBlocoHandler(listarPartidosDoBlocoUC),
		DepListarFrentesHandler:         handlerDeputados.NovoEsferaFederalListarFrentesHandler(listarFrentesUC),
		DepBuscarFrenteHandler:          handlerDeputados.NovoEsferaFederalBuscarFrenteHandler(buscarFrenteUC),
		DepListarMembrosFrenteHandler:   handlerDeputados.NovoEsferaFederalListarMembrosFrenteHandler(listarMembrosFrenteUC),
		DepListarGruposHandler:          handlerDeputados.NovoEsferaFederalListarGruposHandler(listarGruposUC),
		DepBuscarGrupoHandler:           handlerDeputados.NovoEsferaFederalBuscarGrupoHandler(buscarGrupoUC),
		DepListarLegislaturasHandler:    handlerDeputados.NovoEsferaFederalListarLegislaturasHandler(listarLegislaturasUC),
		DepBuscarLegislaturaHandler:     handlerDeputados.NovoEsferaFederalBuscarLegislaturaHandler(buscarLegislaturaUC),
		DepListarVotacoesHandler:        handlerDeputados.NovoEsferaFederalListarVotacoesHandler(depListarVotacoesUC),
		DepBuscarVotacaoHandler:         handlerDeputados.NovoEsferaFederalBuscarVotacaoHandler(buscarVotacaoUC),
		DepListarVotosHandler:           handlerDeputados.NovoEsferaFederalListarVotosHandler(listarVotosUC),

		BaixarDocumentoEmendaHandler: handlerSenadores.NovoBaixarDocumentoEmendaHandler(baixarDocumentoEmendaUC),

		ListarSenadoresHandler: handlerSenadores.NovoListarSenadoresHandler(listarSenadoresUC, redisCache),

		BuscarSenadorHandler:                     handlerSenadores.NovoBuscarSenadorHandler(buscarSenadorUC),
		ListarCargosSenadorHandler:               handlerSenadores.NovoListarCargosHandler(listarCargosSenUC),
		ListarComissoesSenadorHandler:            handlerSenadores.NovoListarComissoesHandler(listarComissoesSenUC),
		ListarMandatosHandler:                    handlerSenadores.NovoListarMandatosHandler(listarMandatosUC),
		ListarOrcamentoHandler:                   handlerSenadores.NovoListarOrcamentoHandler(listarOrcamentoUC),
		ListarProcessosHandler:                   handlerSenadores.NovoListarProcessosHandler(listarProcessosUC),
		ListarProcessoAssuntosHandler:            handlerSenadores.NovoListarProcessoAssuntosHandler(listarProcessoAssuntosUC),
		ListarProcessoEmendasHandler:             handlerSenadores.NovoListarProcessoEmendasHandler(listarProcessoEmendasUC),
		BuscarProcessoHandler:                    handlerSenadores.NovoBuscarProcessoHandler(buscarProcessoUC),
		ListarVotacoesHandler:                    handlerSenadores.NovoListarVotacoesHandler(listarVotacoesUC),
		ListarVotacoesComissaoHandler:            handlerSenadores.NovoListarVotacoesComissaoHandler(listarVotacoesComissaoUC),
		ListarVotacoesComissaoParlamentarHandler: handlerSenadores.NovoListarVotacoesComissaoParlamentarHandler(listarVotacoesComParlUC),
		ListarMateriaTramitacaoHandler:           handlerSenadores.NovoListarMateriaTramitacaoHandler(listarMateriaTramitacaoUC),
		ListarAgendaDiaHandler:                   handlerSenadores.NovoListarAgendaDiaHandler(listarAgendaDiaUC),
		ListarAgendaMesHandler:                   handlerSenadores.NovoListarAgendaMesHandler(listarAgendaMesUC),
		BuscarEncontroHandler:                    handlerSenadores.NovoBuscarEncontroHandler(buscarEncontroUC),
		ListarTodasComissoesHandler:              handlerSenadores.NovoListarTodasComissoesHandler(listarTodasComissoesUC),
		BuscarComissaoHandler:                    handlerSenadores.NovoBuscarComissaoHandler(buscarComissaoUC),

		ListarEstadosHandler:                  handlerEstadual.NovoEsferaEstadualListarEstadosHandler(listarEstadosUC),
		BuscarDadosEstadoHandler:              handlerEstadual.NovoEsferaEstadualBuscarDadosCompletosEstadoHandler(dadosCompletosUC),
		BuscarBasicoEstadoHandler:             handlerEstadual.NovoEsferaEstadualBuscarDadosBasicosEstadoHandler(basicoEstadoUC),
		BuscarCandidatosEstadoHandler:         handlerEstadual.NovoEsferaEstadualBuscarCandidatosHandler(candidatosEstadoUC),
		BuscarDeputadosEstadoHandler:          handlerEstadual.NovoEsferaEstadualBuscarDeputadosHandler(deputadosEstadoUC),
		BuscarSenadoresEstadoHandler:          handlerEstadual.NovoEsferaEstadualBuscarSenadoresHandler(senadoresEstadoUC),
		BuscarMunicipiosPopulacaoHandler:      handlerEstadual.NovoEsferaEstadualBuscarMunicipiosPopulacaoHandler(municipiosPopulacaoUC),
		BuscarRecursosFederaisCompletoHandler: handlerEstadual.NovoEsferaEstadualBuscarRecursosFederaisCompletoHandler(recursosFederaisCompletoUC),

		WSHub: stream.NewHub(
			pncpAnaliseOrgaoHandler,
			pncpAnalisePubHandler,
			anomaliaWorkerHandler,
			despesaPessoalUC,
			despesaCategoriaUC,
			rreoUC,
			recursosFederaisUC,
			detalhesMunicipioUC,
		),

		BuscarSIAPEHandler:    handlerPortalOrgaos.NovoBuscarSIAPEHandler(siapeUC),
		BuscarSIAFIHandler:    handlerPortalOrgaos.NovoBuscarSIAFIHandler(siafiUC),
		BuscarFisicaHandler:   handlerPortalPessoas.NovoBuscarFisicaHandler(fisicaUC),
		BuscarJuridicaHandler: handlerPortalPessoas.NovoBuscarJuridicaHandler(juridicaUC),
		BuscarCartoesHandler:  handlerPortalCartoes.NovoBuscarCartoesHandler(cartoesUC),

		BuscarServidoresHandler:         handlerPortalServidores.NovoBuscarServidoresHandler(servidoresUC),
		BuscarServidorPorIDHandler:      handlerPortalServidores.NovoBuscarServidorPorIDHandler(servidorPorIDUC),
		BuscarRemuneracaoHandler:        handlerPortalServidores.NovoBuscarRemuneracaoHandler(remuneracaoUC),
		BuscarServidoresPorOrgaoHandler: handlerPortalServidores.NovoBuscarServidoresPorOrgaoHandler(servidoresPorOrgaoUC),
		BuscarFuncoesECargosHandler:     handlerPortalServidores.NovoBuscarFuncoesECargosHandler(funcoesECargosUC),
		BuscarPEPsHandler:               handlerPortalServidores.NovoBuscarPEPsHandler(pepsUC),

		BuscarTiposTransferenciaHandler:       handlerPortalDespesas.NovoBuscarTiposTransferenciaHandler(tiposTransferenciaUC),
		BuscarRecursosRecebidosHandler:        handlerPortalDespesas.NovoBuscarRecursosRecebidosHandler(recursosRecebidosUC),
		BuscarDespesasPorOrgaoHandler:         handlerPortalDespesas.NovoBuscarDespesasPorOrgaoHandler(despesasPorOrgaoUC),
		BuscarPorFuncionalProgramaticaHandler: handlerPortalDespesas.NovoBuscarPorFuncionalProgramaticaHandler(funcionalProgramaticaUC),
		BuscarMovimentacaoLiquidaHandler:      handlerPortalDespesas.NovoBuscarMovimentacaoLiquidaHandler(movLiquidaUC),
		BuscarPlanoOrcamentarioHandler:        handlerPortalDespesas.NovoBuscarPlanoOrcamentarioHandler(planoOrcamentarioUC),
		BuscarItensEmpenhoHandler:             handlerPortalDespesas.NovoBuscarItensEmpenhoHandler(itensEmpenhoUC),
		BuscarHistoricoItemEmpenhoHandler:     handlerPortalDespesas.NovoBuscarHistoricoItemEmpenhoHandler(historicoEmpenhoUC),
		BuscarSubfuncoesHandler:               handlerPortalDespesas.NovoBuscarSubfuncoesHandler(subfuncoesUC),
		BuscarProgramasHandler:                handlerPortalDespesas.NovoBuscarProgramasHandler(programasUC),
		ListarFuncionalProgramaticaHandler:    handlerPortalDespesas.NovoListarFuncionalProgramaticaHandler(listarFuncProgramaticaUC),
		BuscarFuncoesHandler:                  handlerPortalDespesas.NovoBuscarFuncoesHandler(funcoesUC),
		BuscarAcoesHandler:                    handlerPortalDespesas.NovoBuscarAcoesHandler(acoesUC),
		BuscarFavorecidosFinaisHandler:        handlerPortalDespesas.NovoBuscarFavorecidosFinaisPorDocumentoHandler(favorecidosUC),
		BuscarEmpenhosImpactadosHandler:       handlerPortalDespesas.NovoBuscarEmpenhosImpactadosHandler(empenhosImpactadosUC),
		BuscarDocumentosHandler:               handlerPortalDespesas.NovoBuscarDocumentosHandler(documentosUC),
		BuscarDocumentoPorCodigoHandler:       handlerPortalDespesas.NovoBuscarDocumentoPorCodigoHandler(documentoPorCodigoUC),
		BuscarDocumentosRelacionadosHandler:   handlerPortalDespesas.NovoBuscarDocumentosRelacionadosHandler(documentosRelacionadosUC),
		BuscarDocumentosPorFavorecidoHandler:  handlerPortalDespesas.NovoBuscarDocumentosPorFavorecidoHandler(documentosPorFavorecidoUC),

		BuscarEmendasHandler:          handlerPortalEmendas.NovoBuscarEmendasHandler(emendasUC),
		BuscarDocumentosEmendaHandler: handlerPortalEmendas.NovoBuscarDocumentosEmendaHandler(documentosEmendaUC),

		MongoClient:             mongoClient,
		GerarDescricaoSvc:       gerarDescricaoSvc,
		AnomaliaWorkerHandler:   anomaliaWorkerHandler,
		AnomaliaConsultaHandler: anomaliaConsultaHandler,

		ConvenioHandler: convenioHandler,
	}
}
