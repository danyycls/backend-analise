package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	repositorio "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/domain"
	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/shared/mongodb"
	"github.com/danyele/podp/internal/shared/services"
	anomalia "github.com/danyele/podp/internal/worker/anomalia"
	"go.mongodb.org/mongo-driver/bson"
)

type AnaliseAnomaliaWorkerUseCase struct {
	normalizarSvc             services.NormalizarDocumentosService
	buscarTseSvc              services.BuscarLigacaoPoliticaTSEService
	verificarTcuSvc           services.VerificarSancoesTCUService
	verificarServSvc          services.VerificarServidorPublicoService
	verificarPessoaPublicaSvc services.VerificarPessoaPublicaService
	gerarDescricaoSvc         services.GerarDescricaoVinculoService
	db                        database.DB
	mongo                     mongodb.Client
}

func NovoAnaliseAnomaliaWorkerUseCase(
	normalizarSvc services.NormalizarDocumentosService,
	buscarTseSvc services.BuscarLigacaoPoliticaTSEService,
	verificarTcuSvc services.VerificarSancoesTCUService,
	verificarServSvc services.VerificarServidorPublicoService,
	verificarPessoaPublicaSvc services.VerificarPessoaPublicaService,
	gerarDescricaoSvc services.GerarDescricaoVinculoService,
	db database.DB,
	mongo mongodb.Client,
) *AnaliseAnomaliaWorkerUseCase {
	return &AnaliseAnomaliaWorkerUseCase{
		normalizarSvc:             normalizarSvc,
		buscarTseSvc:              buscarTseSvc,
		verificarTcuSvc:           verificarTcuSvc,
		verificarServSvc:          verificarServSvc,
		verificarPessoaPublicaSvc: verificarPessoaPublicaSvc,
		gerarDescricaoSvc:         gerarDescricaoSvc,
		db:                        db,
		mongo:                     mongo,
	}
}

func (u *AnaliseAnomaliaWorkerUseCase) Executar(
	ctx context.Context,
	req anomalia.IniciarWorkerRequest,
	eventos chan<- anomalia.WorkerEvento,
) {
	logger.Info("iniciando analise de anomalias",
		"licitacoes", len(req.Licitacoes),
	)

	eventos <- anomalia.WorkerEvento{
		Type:       "started",
		Total:      len(req.Licitacoes),
		EtapaAtual: "analisando_vinculos",
		Message:    "Iniciando análise de vínculos...",
	}

	repo := repositorio.Novo(u.db)
	totalLicitacoes := len(req.Licitacoes)
	anomaliasEncontradas := 0
	processed := 0
	errors := 0

	const maxConcorrencia = 5
	sem := make(chan struct{}, maxConcorrencia)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, item := range req.Licitacoes {
		select {
		case <-ctx.Done():
			logger.Info("analise cancelada durante iteracao", "processados", processed)
			goto cleanup
		default:
		}

		wg.Add(1)
		sem <- struct{}{}

		item := item
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			select {
			case <-ctx.Done():
				return
			default:
			}

			documentos := []documentoExtraido{}

			if item.CpfCnpj != "" {
				documentos = append(documentos, documentoExtraido{
					Documento: item.CpfCnpj,
					Origem:    "principal",
				})
			}

			for _, socio := range item.Socios {
				if socio.Documento != "" {
					documentos = append(documentos, documentoExtraido{
						Documento: socio.Documento,
						Nome:      socio.Nome,
						Origem:    "socio",
					})
				}
			}

			type docResultado struct {
				dv       *services.NormalizarDocumentosOutput
				vinculos []domain.Vinculo
			}

			var resultados []docResultado

			for _, doc := range documentos {
				select {
				case <-ctx.Done():
					return
				default:
				}
				dv, _ := u.normalizarSvc.Executar(ctx, services.NormalizarDocumentosInput{
					Documento: doc.Documento,
					Nome:      doc.Nome,
					Origem:    doc.Origem,
				})
				if dv == nil {
					continue
				}

				mu.Lock()
				eventos <- anomalia.WorkerEvento{
					Type:                 "progress",
					Message:              fmt.Sprintf("Verificando documento %s...", dv.DocumentoNormalizado),
					Processed:            processed,
					Total:                totalLicitacoes,
					AnomaliasEncontradas: anomaliasEncontradas,
					EtapaAtual:           "analisando_vinculos",
					Documento:            dv.DocumentoNormalizado,
				}
				mu.Unlock()

				vinculos := u.buscarTodosVinculos(ctx, repo, dv.DocumentoNormalizado, dv.Nome, dv.Parcial)

				if len(vinculos) == 0 {
					continue
				}

				resultados = append(resultados, docResultado{dv: dv, vinculos: vinculos})
			}

			docVinculos := montarAnormalidadesLicitacao(item)

			if len(resultados) == 0 && len(docVinculos) == 0 {
				mu.Lock()
				processed++
				eventos <- anomalia.WorkerEvento{
					Type:                 "progress",
					Processed:            processed,
					Total:                totalLicitacoes,
					Success:              processed,
					Errors:               errors,
					AnomaliasEncontradas: anomaliasEncontradas,
					EtapaAtual:           "analisando_vinculos",
				}
				mu.Unlock()
				return
			}

			// Unificar todos os resultados num único card por licitação
			documentosUnicos := make(map[string]struct{})

			for _, r := range resultados {
				key := r.dv.DocumentoNormalizado + "|" + r.dv.Origem
				if _, ok := documentosUnicos[key]; ok {
					continue
				}
				documentosUnicos[key] = struct{}{}

				docVinculos = append(docVinculos, domain.DocumentoVinculo{
					DocumentoInput:       r.dv.DocumentoInput,
					DocumentoNormalizado: r.dv.DocumentoNormalizado,
					Nome:                 r.dv.Nome,
					Parcial:              r.dv.Parcial,
					Origem:               r.dv.Origem,
					Vinculos:             r.vinculos,
				})
			}

			sociosOutput := make([]domain.SocioOutput, 0, len(item.Socios))
			for _, s := range item.Socios {
				if s.Documento != "" {
					sociosOutput = append(sociosOutput, domain.SocioOutput{
						Nome:      s.Nome,
						Documento: s.Documento,
					})
				}
			}

			vl := domain.VinculoLicitacao{
				NumeroControlePncp: item.NumeroControlePncp,
				CpfCnpj:            item.CpfCnpj,
				ValorGlobal:        item.ValorGlobal,
				NomeEmpresa:        item.NomeEmpresa,
				Socios:             sociosOutput,
				DocumentosVinculos: docVinculos,
			}

			descricao, errDesc := u.gerarDescricaoSvc.Executar(ctx, services.GerarDescricaoVinculoInput{
				VinculoLicitacao: vl,
			})
			if errDesc != nil {
				logger.Warn("erro ao gerar descricao", "erro", errDesc)
			}

			titulo := ""
			tags := []string{}
			if descricao != nil {
				titulo = descricao.Titulo
				tags = descricao.Tags
			}

			anomaliaDoc := anomalia.AnomaliaDocumento{
				NomeFornecedorPNCP:      item.NomeEmpresa,
				DocumentoFornecedorPNCP: item.CpfCnpj,
				NumeroControlePncp:      item.NumeroControlePncp,
				OrgaoCNPJ:               item.OrgaoCnpj,
				OrgaoNome:               item.OrgaoNome,
				Uf:                      item.Uf,
				Municipio:               item.Municipio,
				Socios:                  sociosOutput,
				DocumentosVinculos:      docVinculos,
				Titulo:                  titulo,
				Tags:                    tags,
				CreatedAt:               time.Now(),
			}

			existing, _ := u.mongo.Find(ctx, "anomalias", bson.M{
				"numero_controle_pncp": item.NumeroControlePncp,
			})
			if len(existing) == 0 {
				if _, err := u.mongo.InsertOne(ctx, "anomalias", anomaliaDoc); err != nil {
					logger.Error("erro ao persistir anomalia", "erro", err)
				}

				mu.Lock()
				anomaliasEncontradas++
				mu.Unlock()
			}

			mu.Lock()
			processed++
			eventos <- anomalia.WorkerEvento{
				Type:                 "progress",
				Processed:            processed,
				Total:                totalLicitacoes,
				Success:              processed,
				Errors:               errors,
				AnomaliasEncontradas: anomaliasEncontradas,
				EtapaAtual:           "analisando_vinculos",
			}
			mu.Unlock()
		}()
	}

cleanup:
	wg.Wait()

	select {
	case <-ctx.Done():
		logger.Info("analise cancelada", "licitacoes", totalLicitacoes, "errors", errors)
		eventos <- anomalia.WorkerEvento{
			Type:                 "cancelled",
			Total:                totalLicitacoes,
			Processed:            processed,
			Errors:               errors,
			AnomaliasEncontradas: anomaliasEncontradas,
			EtapaAtual:           "concluido",
			Message:              "Análise cancelada pelo usuário",
		}
		return
	default:
	}

	logger.Info("analise concluida", "licitacoes", totalLicitacoes, "anomalias", anomaliasEncontradas, "errors", errors)
	eventos <- anomalia.WorkerEvento{
		Type:                 "completed",
		Total:                totalLicitacoes,
		Processed:            processed,
		Success:              processed,
		Errors:               errors,
		AnomaliasEncontradas: anomaliasEncontradas,
		EtapaAtual:           "concluido",
	}
}

type documentoExtraido struct {
	Documento string
	Nome      string
	Origem    string
}

func (u *AnaliseAnomaliaWorkerUseCase) buscarTodosVinculos(
	ctx context.Context,
	repo *repositorio.Repositorio,
	documentoNormalizado string,
	nome string,
	parcial bool,
) []domain.Vinculo {
	var todos []domain.Vinculo

	tseResult, _ := u.buscarTseSvc.Executar(ctx, repo, services.BuscarLigacaoPoliticaTSEInput{
		DocumentoNormalizado: documentoNormalizado,
		Nome:                 nome,
		Parcial:              parcial,
	})
	if tseResult != nil {
		todos = append(todos, tseResult.Vinculos...)
	}

	tcuResult, _ := u.verificarTcuSvc.Executar(ctx, services.VerificarSancoesTCUInput{
		Documentos: []string{documentoNormalizado},
	})
	if tcuResult != nil {
		if v, ok := tcuResult.Resultados[documentoNormalizado]; ok {
			todos = append(todos, v...)
		}
	}

	servResult, _ := u.verificarServSvc.Executar(ctx, services.VerificarServidorInput{
		Documentos: []string{documentoNormalizado},
	})
	if servResult != nil {
		if v, ok := servResult.Resultados[documentoNormalizado]; ok {
			todos = append(todos, v...)
		}
	}

	pessoaPublicaResult, _ := u.verificarPessoaPublicaSvc.Executar(ctx, services.VerificarPessoaPublicaInput{
		Documentos: []string{documentoNormalizado},
	})
	if pessoaPublicaResult != nil {
		if v, ok := pessoaPublicaResult.Resultados[documentoNormalizado]; ok {
			todos = append(todos, v...)
		}
	}

	return todos
}

func montarAnormalidadesLicitacao(item anomalia.LicitacaoInput) []domain.DocumentoVinculo {
	if len(item.Anormalidades) == 0 {
		return nil
	}

	vinculos := make([]domain.Vinculo, 0, len(item.Anormalidades))
	for _, a := range item.Anormalidades {
		vinculo := domain.Vinculo{
			Tipo:      a.Tipo,
			Descricao: a.Descricao,
		}

		if a.Detalhes.DispensaValorLimite != nil {
			d := a.Detalhes.DispensaValorLimite
			vinculo.Detalhes = &domain.VinculoDetalhes{
				DispensaValorLimite: &domain.DispensaValorLimiteDetalhes{
					Modalidade:  d.Modalidade,
					Categoria:   d.Categoria,
					ValorGlobal: d.ValorGlobal,
					Limite:      d.Limite,
					Excedente:   d.Excedente,
					Objeto:      d.Objeto,
					Regra:       d.Regra,
				},
			}
		}

		vinculos = append(vinculos, vinculo)
	}

	return []domain.DocumentoVinculo{
		{
			DocumentoInput:       item.CpfCnpj,
			DocumentoNormalizado: item.CpfCnpj,
			Nome:                 item.NomeEmpresa,
			Origem:               "regra",
			Vinculos:             vinculos,
		},
	}
}
