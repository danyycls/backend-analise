package usecase

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/indexes"
	parse "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/parse"
	repositorios "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/repositorios"
	"github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/service"
	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/podp/internal/shared/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NivelMetrica struct {
	Copia     float64 `json:"copia"`
	Mesclar   float64 `json:"mesclar"`
	Parse     float64 `json:"parse"`
	Total     float64 `json:"total"`
	Registros int64   `json:"registros,omitempty"`
	Operacoes int     `json:"operacoes,omitempty"`
}

type EventoProgressoImportacao struct {
	TotalArquivos          int                     `json:"total_arquivos"`
	TotalDiretorios        int                     `json:"total_diretorios"`
	DiretorioIndice        int                     `json:"diretorio_indice"`
	ArquivosLendo          int                     `json:"arquivos_lendo"`
	ArquivosLidos          int                     `json:"arquivos_lidos"`
	ArquivosPersistindo    int                     `json:"arquivos_persistindo"`
	ArquivosPersistidos    int                     `json:"arquivos_persistidos"`
	ArquivosIgnorados      int                     `json:"arquivos_ignorados"`
	ArquivosRestantes      int                     `json:"arquivos_restantes"`
	Timestamp              string                  `json:"timestamp"`
	DuracaoSegundos        float64                 `json:"duracao_segundos"`
	EtapaAtual             string                  `json:"etapa_atual"`
	EntidadeAtual          string                  `json:"entidade_atual"`
	DuracaoEtapaSegundos   float64                 `json:"duracao_etapa_segundos"`
	TaxaInsercaoPorSegundo float64                 `json:"taxa_insercao_por_segundo"`
	EtapasConcluidas       map[string]float64      `json:"etapas_concluidas,omitempty"`
	ETLTempoPorEntidade    map[string]NivelMetrica `json:"etl_tempo_por_entidade,omitempty"`
}

type ProgressoImportacao struct {
	Total             int
	TotalDiretorios   int
	DiretorioIndice   atomic.Int32
	Lendo             atomic.Int32
	Lidos             atomic.Int32
	AguardandoPersist atomic.Int32
	Persistindo       atomic.Int32
	Persistidos       atomic.Int32
	Ignorados         atomic.Int32
	InicioEtapa       atomic.Value
	EtapaAtual        atomic.Value
	EntidadeAtual     atomic.Value
	inicioGeral       time.Time
}

func (p *ProgressoImportacao) Evento() EventoProgressoImportacao {
	lendo := int(p.Lendo.Load())
	lidos := int(p.Lidos.Load())
	aguardando := int(p.AguardandoPersist.Load())
	persistindo := int(p.Persistindo.Load())
	persistidos := int(p.Persistidos.Load())
	ignorados := int(p.Ignorados.Load())

	concluidos := persistidos + ignorados
	restantes := p.Total - concluidos - lendo - aguardando - persistindo
	if restantes < 0 {
		restantes = 0
	}

	etapa := ""
	if v := p.EtapaAtual.Load(); v != nil {
		etapa = v.(string)
	}
	entidade := ""
	if v := p.EntidadeAtual.Load(); v != nil {
		entidade = v.(string)
	}
	var duracaoEtapa float64
	if v := p.InicioEtapa.Load(); v != nil {
		if inicio, ok := v.(time.Time); ok && !inicio.IsZero() {
			duracaoEtapa = time.Since(inicio).Seconds()
		}
	}
	duracaoTotal := time.Since(p.inicioGeral).Seconds()

	var taxa float64
	if duracaoTotal > 0 && p.Persistidos.Load() > 0 {
		taxa = float64(p.Persistidos.Load()) / duracaoTotal
	}

	return EventoProgressoImportacao{
		TotalArquivos:          p.Total,
		TotalDiretorios:        p.TotalDiretorios,
		DiretorioIndice:        int(p.DiretorioIndice.Load()),
		ArquivosLendo:          lendo,
		ArquivosLidos:          lidos,
		ArquivosPersistindo:    persistindo,
		ArquivosPersistidos:    persistidos,
		ArquivosIgnorados:      ignorados,
		ArquivosRestantes:      restantes,
		EtapaAtual:             etapa,
		EntidadeAtual:          entidade,
		DuracaoEtapaSegundos:   duracaoEtapa,
		DuracaoSegundos:        duracaoTotal,
		TaxaInsercaoPorSegundo: taxa,
	}
}

type ImportarCSVRequest struct{}

type MetricaImportacao struct {
	TempoCopy  string `json:"tempo_copy"`
	TempoMerge string `json:"tempo_merge"`
	Registros  int64  `json:"registros"`
}

type ImportarCSVResponse struct {
	Sucesso             bool                      `json:"sucesso"`
	Status              string                    `json:"status"`
	ArquivosProcessados []tipos.ArquivoProcessado `json:"arquivos_processados"`
	ArquivosComSucesso  []string                  `json:"arquivos_com_sucesso"`
	TotalRegistros      int                       `json:"total_registros"`
	MensagemErro        string                    `json:"mensagem_erro,omitempty"`
	Erro                *ErroImportacao           `json:"erro,omitempty"`
	Metrics             *MetricaImportacao        `json:"metrics,omitempty"`
}

type ErroImportacao struct {
	Mensagem    string `json:"mensagem"`
	Causa       string `json:"causa"`
	ArquivoTipo string `json:"arquivo_tipo,omitempty"`
	UF          string `json:"uf,omitempty"`
	Tabela      string `json:"tabela,omitempty"`
}

func novoErro(err error) *ErroImportacao {
	causa := err
	for e := err; e != nil; e = errors.Unwrap(e) {
		causa = e
	}
	return &ErroImportacao{
		Mensagem: err.Error(),
		Causa:    causa.Error(),
	}
}

type ImportarCSVUseCase interface {
	Executar(ctx context.Context, input ImportarCSVRequest) (*ImportarCSVResponse, error)
	ProgressoEvento() EventoProgressoImportacao
}

type importarCSVUseCase struct {
	pgPool           *pgxpool.Pool
	pgPoolLeitura    *pgxpool.Pool
	pgRepo           *repositorios.Repositorio
	LeitorCSVService service.LeitorCSVServiceInterface
	progression      *ProgressoImportacao
	batchSize        int
	maxWorkers       int
	arquivosPorLote  int
	resultadoMetrica *repositorios.ImportacaoResultado
}

func NovoImportarCSVUseCase(pool, poolLeitura *pgxpool.Pool, leitorCSVService service.LeitorCSVServiceInterface) ImportarCSVUseCase {
	batchSize := tipos.GetEnvInt(tipos.EnvBatchSize, 10000)
	maxWorkers := tipos.GetEnvInt(tipos.EnvMaxWorkers, runtime.NumCPU()*2)
	arquivosPorLote := tipos.GetEnvInt(tipos.EnvFilesPerBatch, 50)

	return &importarCSVUseCase{
		pgPool:           pool,
		pgPoolLeitura:    poolLeitura,
		pgRepo:           repositorios.Novo(pool, poolLeitura),
		LeitorCSVService: leitorCSVService,
		batchSize:        batchSize,
		maxWorkers:       maxWorkers,
		arquivosPorLote:  arquivosPorLote,
	}
}

func (u *importarCSVUseCase) ProgressoEvento() EventoProgressoImportacao {
	if u.progression == nil {
		return EventoProgressoImportacao{}
	}
	evt := u.progression.Evento()
	if u.resultadoMetrica != nil {
		evt.EtapasConcluidas = make(map[string]float64, len(u.resultadoMetrica.Etapas)+1)
		if u.resultadoMetrica.TempoParse > 0 {
			evt.EtapasConcluidas["parse"] = u.resultadoMetrica.TempoParse.Seconds()
		}
		if u.resultadoMetrica.TempoCOPY > 0 {
			evt.EtapasConcluidas["copia"] = u.resultadoMetrica.TempoCOPY.Seconds()
		}
		if u.resultadoMetrica.TempoMerge > 0 {
			evt.EtapasConcluidas["mesclar"] = u.resultadoMetrica.TempoMerge.Seconds()
		}

		if len(u.resultadoMetrica.Niveis) > 0 {
			etl := make(map[string]NivelMetrica, len(u.resultadoMetrica.Niveis))
			for k, v := range u.resultadoMetrica.Niveis {
				if v.Total > 0 {
					etl[k] = NivelMetrica{
						Copia:     v.Copia.Seconds(),
						Mesclar:   v.Mesclar.Seconds(),
						Parse:     v.Parse.Seconds(),
						Total:     v.Total.Seconds(),
						Registros: v.Registros,
						Operacoes: v.Operacoes,
					}
				}
			}
			evt.ETLTempoPorEntidade = etl
		}
	}
	return evt
}

func (u *importarCSVUseCase) setEtapa(etapa string) {
	if u.progression == nil {
		return
	}
	u.progression.EtapaAtual.Store(etapa)
	u.progression.InicioEtapa.Store(time.Now())
}

func (u *importarCSVUseCase) Executar(ctx context.Context, input ImportarCSVRequest) (*ImportarCSVResponse, error) {
	log := logger.New("LeitorCSV: UseCase: Execute")
	arquivos, err := u.LeitorCSVService.ListarArquivos()
	if err != nil {
		r := &ImportarCSVResponse{
			Sucesso:            false,
			Status:             "erro_leitura",
			ArquivosComSucesso: []string{},
			MensagemErro:       err.Error(),
		}
		r.Erro = novoErro(err)
		return r, err
	}

	arquivosImportados, err := u.pgRepo.ListarTodosArquivosImportados(ctx)
	if err != nil {
		log.Warn("nao foi possivel carregar arquivos ja importados", "erro", err)
		arquivosImportados = make(map[string]bool)
	}

	porDir := agruparPorDiretorio(arquivos)
	u.progression = &ProgressoImportacao{
		Total:           len(arquivos),
		TotalDiretorios: len(porDir),
		inicioGeral:     time.Now(),
	}
	u.setEtapa("leitura_parse")
	u.resultadoMetrica = &repositorios.ImportacaoResultado{
		Etapas: make(map[string]time.Duration),
		Niveis: make(map[string]*repositorios.NivelTiming),
	}

	resultado := &ImportarCSVResponse{
		Sucesso:             true,
		Status:              "concluida",
		ArquivosProcessados: make([]tipos.ArquivoProcessado, 0),
		ArquivosComSucesso:  make([]string, 0),
	}
	errGlobal := u.processarPorDiretorio(ctx, porDir, arquivosImportados, resultado)

	if errGlobal != nil {
		resultado.Sucesso = false
		resultado.Status = "erro_importacao"
		resultado.MensagemErro = errGlobal.Error()
		resultado.Erro = novoErro(errGlobal)
		return resultado, errGlobal
	}

	resultado.Metrics = &MetricaImportacao{
		TempoCopy:  u.resultadoMetrica.TempoCOPY.String(),
		TempoMerge: u.resultadoMetrica.TempoMerge.String(),
		Registros:  u.resultadoMetrica.RegistrosInseridos,
	}

	log.Info("resumo da importacao",
		"registros", u.resultadoMetrica.RegistrosInseridos,
		"parse", u.resultadoMetrica.TempoParse.String(),
		"copy", u.resultadoMetrica.TempoCOPY.String(),
		"merge", u.resultadoMetrica.TempoMerge.String())

	return resultado, nil
}

type arquivoNoLote struct {
	arquivo    tipos.ArquivoImportacao
	registros  int
	hashSHA256 string
}

type leituraOrdenada struct {
	idx      int
	arquivo  tipos.ArquivoImportacao
	leitura  *tipos.ArquivoProcessado
	dados    *tipos.DadosImportacao
	ignorado bool
	err      error
}

func prioridadeDiretorio(nome string) int {
	d := nome
	switch {
	case strings.Contains(d, "portaltransparencia"):
		return 0
	case strings.Contains(d, "consulta_cand"):
		return 1
	case strings.Contains(d, "bem_candidato"):
		return 2
	case strings.Contains(d, "candidatos"):
		return 3
	case strings.Contains(d, "orgaos_partidarios"), strings.Contains(d, "orgao_partidario"):
		return 4
	default:
		return 99
	}
}

func agruparPorDiretorio(arquivos []tipos.ArquivoImportacao) map[string][]tipos.ArquivoImportacao {
	porDir := make(map[string][]tipos.ArquivoImportacao)
	for _, a := range arquivos {
		dir := a.Diretorio
		if dir == "" {
			dir = "."
		}
		porDir[dir] = append(porDir[dir], a)
	}
	for dir := range porDir {
		sort.Slice(porDir[dir], func(i, j int) bool {
			pi := parse.PrioridadeTipoArquivo(porDir[dir][i].Tipo)
			pj := parse.PrioridadeTipoArquivo(porDir[dir][j].Tipo)
			if pi == pj {
				return porDir[dir][i].Nome < porDir[dir][j].Nome
			}
			return pi < pj
		})
	}
	return porDir
}

func ordenarDiretorios(dirs []string) {
	sort.Slice(dirs, func(i, j int) bool {
		pi := prioridadeDiretorio(dirs[i])
		pj := prioridadeDiretorio(dirs[j])
		if pi == pj {
			return dirs[i] < dirs[j]
		}
		return pi < pj
	})
}

func diretorioPrecisaCandidatos(dir string) bool {
	return strings.Contains(dir, "bem_candidato") ||
		strings.Contains(dir, "candidatos") ||
		strings.Contains(dir, "orgaos_partidarios") ||
		strings.Contains(dir, "orgao_partidario")
}

func diretorioSomenteTransacional(dir string) bool {
	return strings.Contains(dir, "bem_candidato")
}

func (u *importarCSVUseCase) dropIndexes(ctx context.Context) error {
	conn, err := u.pgPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("adquirir conexao para drop indexes: %w", err)
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := indexes.DropSecondaryIndexes(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (u *importarCSVUseCase) recreateIndexes(ctx context.Context) error {
	conn, err := u.pgPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("adquirir conexao para recreate indexes: %w", err)
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := indexes.RecreateSecondaryIndexes(ctx, tx); err != nil {
		return err
	}
	if err := indexes.AnalyzeTables(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (u *importarCSVUseCase) ensureConstraintIndexes(ctx context.Context) error {
	conn, err := u.pgPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("adquirir conexao para constraint indexes: %w", err)
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := indexes.RecreateConstraintIndexes(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (u *importarCSVUseCase) processarPorDiretorio(
	ctx context.Context,
	porDir map[string][]tipos.ArquivoImportacao,
	arquivosImportados map[string]bool,
	resultado *ImportarCSVResponse,
) error {
	log := logger.New("LeitorCSV: UseCase: processarPorDiretorio")
	dirs := make([]string, 0, len(porDir))
	dirLower := make(map[string]string, len(porDir))
	for d, arquivos := range porDir {
		dirs = append(dirs, d)
		if len(arquivos) > 0 {
			dirLower[d] = arquivos[0].DiretorioLower
		}
	}
	ordenarDiretorios(dirs)

	log.Info("removendo indices secundarios para acelerar importacao")
	if err := u.dropIndexes(ctx); err != nil {
		return fmt.Errorf("drop indexes: %w", err)
	}

	log.Info("garantindo indices de constraint para importacao")
	if err := u.ensureConstraintIndexes(ctx); err != nil {
		return fmt.Errorf("ensure constraint indexes: %w", err)
	}

	var cacheCandidatos *tipos.DadosImportacao
	candidatosCarregados := false

	for i, dir := range dirs {
		u.progression.DiretorioIndice.Store(int32(i + 1))
		log.Info("iniciando diretorio", "indice", i+1, "total", len(dirs), "diretorio", dir, "arquivos", len(porDir[dir]))
		inicioDir := time.Now()

		if diretorioPrecisaCandidatos(dirLower[dir]) && !candidatosCarregados {
			cacheCandidatos = tipos.NovoDadosImportacao()
			total, err := u.pgRepo.CarregarCandidatosNoMapa(ctx, cacheCandidatos.Candidatos)
			if err != nil {
				return fmt.Errorf("carregar cache de candidatos: %w", err)
			}
			log.Info("cache de candidatos carregado do banco", "total", total)
			candidatosCarregados = true
		}

		if err := u.processarDiretorio(ctx, dir, dirLower[dir], porDir[dir], arquivosImportados, resultado, cacheCandidatos); err != nil {
			return err
		}

		log.Info("diretorio concluido", "diretorio", dir, "duracao", time.Since(inicioDir))
	}

	if cacheCandidatos != nil {
		cacheCandidatos.Candidatos = nil
		cacheCandidatos = nil
	}

	log.Info("recriando indices secundarios e atualizando estatisticas")
	if err := u.recreateIndexes(ctx); err != nil {
		return fmt.Errorf("recreate indexes: %w", err)
	}

	u.progression.DiretorioIndice.Store(0)
	return nil
}

func (u *importarCSVUseCase) processarDiretorio(
	ctx context.Context,
	dir string,
	dirLower string,
	arquivos []tipos.ArquivoImportacao,
	arquivosImportados map[string]bool,
	resultado *ImportarCSVResponse,
	cacheCandidatos *tipos.DadosImportacao,
) error {
	log := logger.New("LeitorCSV: UseCase: processarDiretorio")
	if cacheCandidatos != nil && len(cacheCandidatos.Candidatos) == 0 {
		cacheCandidatos = nil
	}

	n := len(arquivos)
	var muResultado sync.Mutex

	tamanhoLoteFiles := u.arquivosPorLote
	if tamanhoLoteFiles < 1 {
		tamanhoLoteFiles = 1
	}

	dadosAcumulados := tipos.NovoDadosImportacao()

	for batchStart := 0; batchStart < n; batchStart += tamanhoLoteFiles {
		batchEnd := batchStart + tamanhoLoteFiles
		if batchEnd > n {
			batchEnd = n
		}
		batchArquivos := arquivos[batchStart:batchEnd]

		lotePendente, err := u.processarBatch(ctx, dir, dirLower, batchArquivos, arquivosImportados, cacheCandidatos, dadosAcumulados, resultado, &muResultado)
		if err != nil {
			return err
		}

		if len(lotePendente) > 0 {
			if err := u.validarEPersistirLote(ctx, dir, dadosAcumulados, lotePendente); err != nil {
				return err
			}
		}

		parse.LimparDadosAposPersistencia(dadosAcumulados)

		muResultado.Lock()
		if len(resultado.ArquivosProcessados) > 100 {
			resultado.ArquivosProcessados = resultado.ArquivosProcessados[len(resultado.ArquivosProcessados)-100:]
		}
		if len(resultado.ArquivosComSucesso) > 100 {
			resultado.ArquivosComSucesso = resultado.ArquivosComSucesso[len(resultado.ArquivosComSucesso)-100:]
		}
		muResultado.Unlock()

		log.Info("lote processado")
	}

	parse.LimparTodosDados(dadosAcumulados)
	dadosAcumulados = nil //nolint:ineffassign

	debug.FreeOSMemory()
	return nil
}

func (u *importarCSVUseCase) processarBatch(
	ctx context.Context,
	dir string,
	dirLower string,
	batchArquivos []tipos.ArquivoImportacao,
	arquivosImportados map[string]bool,
	cacheCandidatos *tipos.DadosImportacao,
	dadosAcumulados *tipos.DadosImportacao,
	resultado *ImportarCSVResponse,
	muResultado *sync.Mutex,
) ([]arquivoNoLote, error) {
	log := logger.New("LeitorCSV: UseCase: processarBatch")
	batchSize := len(batchArquivos)
	inicioBatch := time.Now()

	slots := make([]*leituraOrdenada, batchSize)
	jobs := make(chan int, batchSize)
	var wg sync.WaitGroup
	var errGlobal atomic.Pointer[error]

	worker := func() {
		defer wg.Done()
		for idxInBatch := range jobs {
			if errGlobal.Load() != nil {
				return
			}
			select {
			case <-ctx.Done():
				setErrAtomic(&errGlobal, ctx.Err())
				return
			default:
			}

			arquivo := batchArquivos[idxInBatch]
			u.progression.Lendo.Add(1)

			slot := &leituraOrdenada{idx: idxInBatch, arquivo: arquivo}

			if arquivosImportados[arquivo.Nome] {
				u.progression.Lendo.Add(-1)
				u.progression.Ignorados.Add(1)
				slot.ignorado = true
				slots[idxInBatch] = slot
				continue
			}

			var proc *parse.ProcessadorLeitorCSV
			if cacheCandidatos != nil && len(cacheCandidatos.Candidatos) > 0 {
				proc = parse.NovoProcessadorComCacheCandidatos(u.batchSize, cacheCandidatos.Candidatos)
			} else {
				proc = parse.NovoProcessadorLeitorCSV(u.batchSize)
			}
			proc.ComResolverCandidato(func(ctx context.Context, sq int64) (uuid.UUID, error) {
				return u.pgRepo.BuscarIDCandidatoPorSQ(ctx, sq)
			})
			leitura, err := u.lerArquivo(ctx, arquivo, proc)
			u.progression.Lendo.Add(-1)
			if err != nil {
				slot.err = err
				slots[idxInBatch] = slot
				setErrAtomic(&errGlobal, err)
				continue
			}
			slot.leitura = leitura
			slot.dados = proc.Dados()
			slots[idxInBatch] = slot
		}
	}

	workers := u.maxWorkers
	if workers > batchSize {
		workers = batchSize
	}
	if workers < 1 {
		workers = 1
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}
	for idxInBatch := 0; idxInBatch < batchSize; idxInBatch++ {
		jobs <- idxInBatch
	}
	close(jobs)
	wg.Wait()

	lotePendente := make([]arquivoNoLote, 0, batchSize)

	for idxInBatch := 0; idxInBatch < batchSize; idxInBatch++ {
		slot := slots[idxInBatch]
		if slot == nil {
			continue
		}
		if slot.err != nil {
			log.Warn("arquivo com erro ignorado no lote",
				"diretorio", dir, "arquivo", slot.arquivo.Nome, "erro", slot.err)
			continue
		}
		if slot.ignorado {
			log.Info("arquivo ja importado, ignorando", "diretorio", dir, "arquivo", slot.arquivo.Nome)
			continue
		}

		u.progression.Lidos.Add(1)
		u.progression.AguardandoPersist.Add(1)

		if diretorioSomenteTransacional(dirLower) {
			parse.MergeDadosTransacionais(dadosAcumulados, slot.dados)
		} else {
			MergeDadosParaColetor(dadosAcumulados, slot.dados)
		}
		lotePendente = append(lotePendente, arquivoNoLote{
			arquivo:    slot.arquivo,
			registros:  slot.leitura.Registros,
			hashSHA256: slot.leitura.HashSHA256,
		})

		muResultado.Lock()
		resultado.ArquivosProcessados = append(resultado.ArquivosProcessados, *slot.leitura)
		resultado.ArquivosComSucesso = append(resultado.ArquivosComSucesso, slot.leitura.NomeArquivo)
		resultado.TotalRegistros += slot.leitura.Registros
		muResultado.Unlock()
	}

	for idxInBatch := range slots {
		if slots[idxInBatch] != nil {
			slots[idxInBatch].dados = nil
			slots[idxInBatch].leitura = nil
			slots[idxInBatch] = nil
		}
	}
	slots = nil

	durParse := time.Since(inicioBatch)
	log.Info("batch lido e parseado", "arquivos", len(lotePendente), "duracao", durParse)
	if u.resultadoMetrica != nil {
		u.resultadoMetrica.TempoParse += durParse
	}
	return lotePendente, nil
}

func setErrAtomic(p *atomic.Pointer[error], err error) {
	if err == nil {
		return
	}
	p.CompareAndSwap(nil, &err)
}

func (u *importarCSVUseCase) validarEPersistirLote(
	ctx context.Context,
	dir string,
	dados *tipos.DadosImportacao,
	lotePendente []arquivoNoLote,
) error {
	qtd := int32(len(lotePendente))
	u.progression.Persistindo.Add(qtd)

	_ = parse.ValidarFKsEmMemoria(dados)

	resultadoImportacao, err := u.persistirLote(ctx, dados, lotePendente)
	if err != nil {
		u.progression.Persistindo.Add(-qtd)
		return fmt.Errorf("diretorio %s: %w", dir, err)
	}
	if resultadoImportacao != nil && u.resultadoMetrica != nil {
		u.resultadoMetrica.TempoCOPY += resultadoImportacao.TempoCOPY
		u.resultadoMetrica.TempoMerge += resultadoImportacao.TempoMerge
		u.resultadoMetrica.TempoParse += resultadoImportacao.TempoParse
		u.resultadoMetrica.RegistrosInseridos += resultadoImportacao.RegistrosInseridos
		u.resultadoMetrica.Etapas["copia"] += resultadoImportacao.TempoCOPY
		u.resultadoMetrica.Etapas["mesclar"] += resultadoImportacao.TempoMerge

		for k, v := range resultadoImportacao.Niveis {
			if u.resultadoMetrica.Niveis[k] == nil {
				u.resultadoMetrica.Niveis[k] = &repositorios.NivelTiming{}
			}
			u.resultadoMetrica.Niveis[k].Copia += v.Copia
			u.resultadoMetrica.Niveis[k].Mesclar += v.Mesclar
			u.resultadoMetrica.Niveis[k].Parse += v.Parse
			u.resultadoMetrica.Niveis[k].Total += v.Total
			u.resultadoMetrica.Niveis[k].Registros += v.Registros
			u.resultadoMetrica.Niveis[k].Operacoes += v.Operacoes
		}
	}

	u.progression.Persistindo.Add(-qtd)
	u.progression.Persistidos.Add(qtd)
	for range lotePendente {
		u.progression.AguardandoPersist.Add(-1)
	}
	return nil
}

func (u *importarCSVUseCase) persistirLote(
	ctx context.Context,
	dados *tipos.DadosImportacao,
	arquivos []arquivoNoLote,
) (*repositorios.ImportacaoResultado, error) {
	log := logger.New("LeitorCSV: UseCase: persistirLote")
	conn, err := u.pgPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("adquirir conexao para persistencia: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("iniciar transacao persistencia lote: %w", err)
	}
	defer tx.Rollback(ctx)

	log.Info("persistindo lote", "arquivos", len(arquivos), "registros_estimados", contarRegistrosEstimados(dados))
	inicio := time.Now()

	resultado := &repositorios.ImportacaoResultado{
		Etapas: make(map[string]time.Duration),
		Niveis: make(map[string]*repositorios.NivelTiming),
	}
	resultado.SetEntidade = func(entidade string) {
		u.progression.EntidadeAtual.Store(entidade)
	}

	u.setEtapa("copy")
	if err := parse.PersistirDadosImportacaoPgCopy(ctx, tx, u.pgRepo, dados, u.batchSize, resultado); err != nil {
		return nil, err
	}

	u.setEtapa("merge")
	for _, item := range arquivos {
		uf := parse.ObterUFDoNomeArquivo(item.arquivo.Nome)
		if err := u.pgRepo.RegistrarArquivoImportado(
			ctx, tx,
			item.arquivo.CaminhoRelativo,
			item.arquivo.Nome,
			item.arquivo.Tipo,
			uf,
			item.registros,
			item.hashSHA256,
		); err != nil {
			return nil, fmt.Errorf("registrar %s: %w", item.arquivo.Nome, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit persistencia lote: %w", err)
	}

	log.Info("lote persistido",
		"duracao", time.Since(inicio),
		"arquivos", len(arquivos),
		"registros", resultado.RegistrosInseridos,
		"parse", resultado.TempoParse.String(),
		"copy", resultado.TempoCOPY.String(),
		"merge", resultado.TempoMerge.String())
	return resultado, nil
}

func contarRegistrosEstimados(d *tipos.DadosImportacao) int {
	if d == nil {
		return 0
	}
	total := 0
	total += len(d.Convenios)
	total += len(d.DespesasCandidato)
	total += len(d.DespesasOrgaoPartidario)
	total += len(d.ReceitasCandidato)
	total += len(d.ReceitasOrgaoPartidario)
	total += len(d.ReceitasDoadorOriginarioCandidato)
	total += len(d.ReceitasDoadorOriginarioOrgaoPartidario)
	total += len(d.BensCandidato)
	total += len(d.Prestacoes)
	return total
}

func MergeDadosParaColetor(dst *tipos.DadosImportacao, src *tipos.DadosImportacao) {
	parse.MergeDados(dst, src)
}

func (u *importarCSVUseCase) lerArquivo(ctx context.Context, arquivo tipos.ArquivoImportacao, processador *parse.ProcessadorLeitorCSV) (*tipos.ArquivoProcessado, error) {
	log := logger.New("LeitorCSV: UseCase: lerArquivo")
	uf := parse.ObterUFDoNomeArquivo(arquivo.Nome)
	inicioArquivo := time.Now()
	log.Info("lendo arquivo", "diretorio", arquivo.Diretorio, "arquivo", arquivo.Nome, "tipo", arquivo.Tipo, "uf", uf)
	leitura, err := u.LeitorCSVService.LerArquivo(ctx, processador, arquivo)
	if err != nil {
		return nil, err
	}
	log.Info("arquivo lido", "diretorio", arquivo.Diretorio, "arquivo", arquivo.Nome, "tipo", arquivo.Tipo, "uf", uf, "duracao", time.Since(inicioArquivo))
	return leitura, nil
}
