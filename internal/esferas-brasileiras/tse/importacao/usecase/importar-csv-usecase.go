package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	parse "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/parse"
	repositorios "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/repositorios"
	"github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/service"
	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/danyele/podp/internal/shared/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EventoProgressoImportacao payload numerico do SSE event:progression
type EventoProgressoImportacao struct {
	TotalArquivos       int `json:"total_arquivos"`
	TotalDiretorios     int `json:"total_diretorios"`
	DiretorioIndice     int `json:"diretorio_indice"`
	ArquivosLendo       int `json:"arquivos_lendo"`
	ArquivosLidos       int `json:"arquivos_lidos"`
	ArquivosPersistindo int `json:"arquivos_persistindo"`
	ArquivosPersistidos int `json:"arquivos_persistidos"`
	ArquivosIgnorados   int `json:"arquivos_ignorados"`
	ArquivosRestantes   int `json:"arquivos_restantes"`
}

// ProgressoImportacao rastreia o estado detalhado da importacao em tempo real
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

	return EventoProgressoImportacao{
		TotalArquivos:       p.Total,
		TotalDiretorios:     p.TotalDiretorios,
		DiretorioIndice:     int(p.DiretorioIndice.Load()),
		ArquivosLendo:       lendo,
		ArquivosLidos:       lidos,
		ArquivosPersistindo: persistindo,
		ArquivosPersistidos: persistidos,
		ArquivosIgnorados:   ignorados,
		ArquivosRestantes:   restantes,
	}
}

// ImportarCSVRequest define os parâmetros de entrada
type ImportarCSVRequest struct{}

// ImportarCSVResponse define a estrutura de saída
type ImportarCSVResponse struct {
	Sucesso             bool                      `json:"sucesso"`
	Status              string                    `json:"status"`
	ArquivosProcessados []tipos.ArquivoProcessado `json:"arquivos_processados"`
	ArquivosComSucesso  []string                  `json:"arquivos_com_sucesso"`
	TotalRegistros      int                       `json:"total_registros"`
	MensagemErro        string                    `json:"mensagem_erro,omitempty"`
	Erro                *ErroImportacao           `json:"erro,omitempty"`
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

// ImportarCSVUseCase define a interface para o caso de uso
type ImportarCSVUseCase interface {
	Executar(ctx context.Context, input ImportarCSVRequest) (*ImportarCSVResponse, error)
	ProgressoEvento() EventoProgressoImportacao
}

// importarCSVUseCase é a implementação
type importarCSVUseCase struct {
	pgPool           *pgxpool.Pool
	pgRepo           *repositorios.Repositorio
	LeitorCSVService service.LeitorCSVServiceInterface
	progression      *ProgressoImportacao
	batchSize        int
	maxWorkers       int
	arquivosPorLote  int
}

// NovoImportarCSVUseCase cria uma nova instância
func NovoImportarCSVUseCase(pool *pgxpool.Pool, leitorCSVService service.LeitorCSVServiceInterface) ImportarCSVUseCase {
	batchSize := 2000
	if v := os.Getenv("IMPORT_BATCH_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			batchSize = n
		}
	}
	maxWorkers := 4
	if v := os.Getenv("IMPORT_MAX_WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxWorkers = n
		}
	}
	arquivosPorLote := 2
	if v := os.Getenv("IMPORT_FILES_PER_BATCH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			arquivosPorLote = n
		}
	}

	return &importarCSVUseCase{
		pgPool:           pool,
		pgRepo:           repositorios.Novo(pool),
		LeitorCSVService: leitorCSVService,
		batchSize:        batchSize,
		maxWorkers:       maxWorkers,
		arquivosPorLote:  arquivosPorLote,
	}
}

// ProgressoEvento retorna contadores numericos para o SSE de progression
func (u *importarCSVUseCase) ProgressoEvento() EventoProgressoImportacao {
	if u.progression == nil {
		return EventoProgressoImportacao{}
	}
	return u.progression.Evento()
}

// Execute executa o caso de uso
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

	return resultado, nil
}

type arquivoNoLote struct {
	arquivo   tipos.ArquivoImportacao
	registros int
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
	d := strings.ToLower(nome)
	switch {
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
	d := strings.ToLower(dir)
	return strings.Contains(d, "bem_candidato") ||
		strings.Contains(d, "candidatos") ||
		strings.Contains(d, "orgaos_partidarios") ||
		strings.Contains(d, "orgao_partidario")
}

func diretorioSomenteTransacional(dir string) bool {
	return strings.Contains(strings.ToLower(dir), "bem_candidato")
}

func (u *importarCSVUseCase) processarPorDiretorio(
	ctx context.Context,
	porDir map[string][]tipos.ArquivoImportacao,
	arquivosImportados map[string]bool,
	resultado *ImportarCSVResponse,
) error {
	log := logger.New("LeitorCSV: UseCase: processarPorDiretorio")
	dirs := make([]string, 0, len(porDir))
	for d := range porDir {
		dirs = append(dirs, d)
	}
	ordenarDiretorios(dirs)

	for i, dir := range dirs {
		u.progression.DiretorioIndice.Store(int32(i + 1))
		log.Info("iniciando diretorio", "indice", i+1, "total", len(dirs), "diretorio", dir, "arquivos", len(porDir[dir]))
		if err := u.processarDiretorio(ctx, dir, porDir[dir], arquivosImportados, resultado); err != nil {
			return err
		}
		log.Info("diretorio concluido", "diretorio", dir)

		runtime.GC()
		debug.FreeOSMemory()
		log.Info("memoria liberada apos diretorio", "diretorio", dir)
	}

	runtime.GC()
	debug.FreeOSMemory()
	log.Info("memoria liberada ao final de todos os diretorios")

	u.progression.DiretorioIndice.Store(0)
	return nil
}

func (u *importarCSVUseCase) processarDiretorio(
	ctx context.Context,
	dir string,
	arquivos []tipos.ArquivoImportacao,
	arquivosImportados map[string]bool,
	resultado *ImportarCSVResponse,
) error {
	log := logger.New("LeitorCSV: UseCase: processarDiretorio")
	var cacheCandidatos *tipos.DadosImportacao
	if diretorioPrecisaCandidatos(dir) {
		cacheCandidatos = tipos.NovoDadosImportacao()
		total, err := u.pgRepo.CarregarCandidatosNoMapa(ctx, cacheCandidatos.Candidatos)
		if err != nil {
			return fmt.Errorf("carregar cache de candidatos: %w", err)
		}
		log.Info("cache de candidatos carregado do banco", "total", total)
	}
	defer func() {
		if cacheCandidatos != nil {
			cacheCandidatos.Candidatos = nil
			cacheCandidatos = nil
		}
	}()

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
		batchSize := len(batchArquivos)

		slots := make([]*leituraOrdenada, batchSize)
		var wg sync.WaitGroup
		jobs := make(chan int, batchSize)
		errOnce := sync.Once{}
		var errGlobal error
		setErr := func(err error) {
			if err == nil {
				return
			}
			errOnce.Do(func() { errGlobal = err })
		}

		worker := func() {
			defer wg.Done()
			for idxInBatch := range jobs {
				if errGlobal != nil {
					return
				}
				select {
				case <-ctx.Done():
					setErr(ctx.Err())
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
					setErr(err)
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

		if errGlobal != nil {
			return fmt.Errorf("diretorio %s: %w", dir, errGlobal)
		}

		dadosColetados := dadosAcumulados
		lotePendente := make([]arquivoNoLote, 0, batchSize)

		for idxInBatch := 0; idxInBatch < batchSize; idxInBatch++ {
			slot := slots[idxInBatch]
			if slot == nil {
				return fmt.Errorf("diretorio %s: leitura ausente no indice %d do lote", dir, idxInBatch)
			}
			if slot.err != nil {
				return fmt.Errorf("diretorio %s: %w", dir, slot.err)
			}

			if slot.ignorado {
				log.Info("arquivo ja importado, ignorando", "diretorio", dir, "arquivo", slot.arquivo.Nome)
				continue
			}

			u.progression.Lidos.Add(1)
			u.progression.AguardandoPersist.Add(1)

			if diretorioSomenteTransacional(dir) {
				parse.MergeDadosTransacionais(dadosColetados, slot.dados)
			} else {
				MergeDadosParaColetor(dadosColetados, slot.dados)
			}
			lotePendente = append(lotePendente, arquivoNoLote{
				arquivo:   slot.arquivo,
				registros: slot.leitura.Registros,
			})

			muResultado.Lock()
			resultado.ArquivosProcessados = append(resultado.ArquivosProcessados, *slot.leitura)
			resultado.ArquivosComSucesso = append(resultado.ArquivosComSucesso, slot.leitura.NomeArquivo)
			resultado.TotalRegistros += slot.leitura.Registros
			muResultado.Unlock()
		}

		if len(lotePendente) > 0 {
			qtd := int32(len(lotePendente))
			u.progression.Persistindo.Add(qtd)

			if err := u.persistirLote(ctx, dadosColetados, lotePendente); err != nil {
				u.progression.Persistindo.Add(-qtd)
				return fmt.Errorf("diretorio %s: %w", dir, err)
			}

			u.progression.Persistindo.Add(-qtd)
			u.progression.Persistidos.Add(qtd)
			for range lotePendente {
				u.progression.AguardandoPersist.Add(-1)
			}
		}

		for idxInBatch := range slots {
			if slots[idxInBatch] != nil {
				slots[idxInBatch].dados = nil
				slots[idxInBatch].leitura = nil
				slots[idxInBatch] = nil
			}
		}
		slots = nil
		lotePendente = nil //nolint:ineffassign

		parse.LimparDadosAposPersistencia(dadosAcumulados)

		muResultado.Lock()
		if len(resultado.ArquivosProcessados) > 100 {
			resultado.ArquivosProcessados = resultado.ArquivosProcessados[len(resultado.ArquivosProcessados)-100:]
		}
		if len(resultado.ArquivosComSucesso) > 100 {
			resultado.ArquivosComSucesso = resultado.ArquivosComSucesso[len(resultado.ArquivosComSucesso)-100:]
		}
		muResultado.Unlock()

		log.Info("forcando garbage collection pos-lote")
		runtime.GC()
		debug.FreeOSMemory()
	}

	parse.LimparTodosDados(dadosAcumulados)
	dadosAcumulados = nil //nolint:ineffassign

	return nil
}

func (u *importarCSVUseCase) persistirLote(
	ctx context.Context,
	dados *tipos.DadosImportacao,
	arquivos []arquivoNoLote,
) error {
	log := logger.New("LeitorCSV: UseCase: persistirLote")
	conn, err := u.pgPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("adquirir conexao para persistencia: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("iniciar transacao persistencia lote: %w", err)
	}
	defer tx.Rollback(ctx)

	log.Info("persistindo lote", "arquivos", len(arquivos), "registros_estimados", contarRegistrosEstimados(dados))
	inicio := time.Now()

	if err := parse.PersistirDadosImportacaoPgCopy(ctx, tx, u.pgRepo, dados, u.batchSize); err != nil {
		return err
	}

	for _, item := range arquivos {
		uf := parse.ObterUFDoNomeArquivo(item.arquivo.Nome)
		if err := u.pgRepo.RegistrarArquivoImportado(
			ctx, tx,
			item.arquivo.CaminhoRelativo,
			item.arquivo.Nome,
			item.arquivo.Tipo,
			uf,
			item.registros,
		); err != nil {
			return fmt.Errorf("registrar %s: %w", item.arquivo.Nome, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit persistencia lote: %w", err)
	}

	log.Info("lote persistido", "duracao", time.Since(inicio), "arquivos", len(arquivos))
	return nil
}

func contarRegistrosEstimados(d *tipos.DadosImportacao) int {
	if d == nil {
		return 0
	}
	total := 0
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
	log.Info("lendo arquivo", "diretorio", arquivo.Diretorio, "arquivo", arquivo.Nome, "tipo", arquivo.Tipo, "uf", uf)
	return u.LeitorCSVService.LerArquivo(ctx, processador, arquivo)
}
