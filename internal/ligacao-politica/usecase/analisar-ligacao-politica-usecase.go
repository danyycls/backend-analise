package usecase

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/logger"

	repositorio "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/shared/clients/opencnpj"
	"github.com/danyele/podp/internal/shared/clients/portaltransparencia"
	tcu "github.com/danyele/podp/internal/shared/clients/tcu"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/shared/utils"
)

type AnalisarLigacaoPoliticaUseCase struct {
	db                  database.DB
	opencnpj            opencnpj.Client
	tcu                 tcu.Client
	portaltransparencia *portaltransparencia.PortalTransparenciaClient
}

func NovoAnalisarLigacaoPoliticaUseCase(db database.DB, opencnpj opencnpj.Client, tcu tcu.Client, portaltransparencia *portaltransparencia.PortalTransparenciaClient) *AnalisarLigacaoPoliticaUseCase {
	return &AnalisarLigacaoPoliticaUseCase{db: db, opencnpj: opencnpj, tcu: tcu, portaltransparencia: portaltransparencia}
}

func (u *AnalisarLigacaoPoliticaUseCase) Executar(ctx context.Context, licitacoes []AnalisarLigacaoPoliticaRequest) (*AnalisarLigacaoPoliticaResponse, error) {
	log := logger.New("LigacaoPolitica: UseCase: Executar")
	log.Info("licitacoes recebidas", "quantidade", len(licitacoes))

	repo := repositorio.Novo(u.db)

	resultados := make([]VinculoLicitacao, 0, len(licitacoes))

	for _, lic := range licitacoes {
		vl := u.processarLicitacao(ctx, repo, lic)
		resultados = append(resultados, vl)
	}

	u.enriquecerFornecedores(ctx, resultados)
	u.enriquecerComTCU(ctx, resultados)
	u.enriquecerComServidorPublico(ctx, resultados)

	return &AnalisarLigacaoPoliticaResponse{
		DocumentosProcessados: len(resultados),
		Resultados:            resultados,
	}, nil
}

func (u *AnalisarLigacaoPoliticaUseCase) processarLicitacao(
	ctx context.Context,
	repo *repositorio.Repositorio,
	lic AnalisarLigacaoPoliticaRequest,
) VinculoLicitacao {
	vl := VinculoLicitacao{
		NumeroControlePncp: lic.NumeroControlePncp,
		CpfCnpj:            lic.CpfCnpj,
	}

	sociosOut := make([]SocioOutput, 0, len(lic.Socios))
	for _, s := range lic.Socios {
		sociosOut = append(sociosOut, SocioOutput(s))

		dv := u.processarDocumento(ctx, repo, s.Documento, s.Nome, "socio")
		if dv != nil {
			vl.Documentos = append(vl.Documentos, *dv)
		}
	}
	vl.Socios = sociosOut

	dv := u.processarDocumento(ctx, repo, lic.CpfCnpj, "", "principal")
	if dv != nil {
		vl.Documentos = append(vl.Documentos, *dv)
	}

	return vl
}

func (u *AnalisarLigacaoPoliticaUseCase) processarDocumento(
	ctx context.Context,
	repo *repositorio.Repositorio,
	docInput string,
	nome string,
	origem string,
) *DocumentoVinculo {
	doc, parcial := utils.NormalizarDocumento(docInput)
	if len(doc) < 3 {
		return nil
	}

	dv := &DocumentoVinculo{
		DocumentoInput:       docInput,
		DocumentoNormalizado: doc,
		Nome:                 nome,
		Parcial:              parcial,
		Origem:               origem,
	}

	if parcial && nome != "" {
		u.buscarLigacoesPorDocumentoParcial(ctx, repo, dv, doc, nome)
	} else {
		u.buscarLigacoesPorDocumentoExato(ctx, repo, dv, doc)
	}

	return dv
}

func (u *AnalisarLigacaoPoliticaUseCase) buscarLigacoesPorDocumentoExato(
	ctx context.Context,
	repo *repositorio.Repositorio,
	dv *DocumentoVinculo,
	doc string,
) {
	isCPF := len(doc) == 11

	docsFornecedor := []string{doc}
	if isCPF {
		docsFornecedor = append(docsFornecedor, "000"+doc)
	}
	forns, _ := repo.FornecedoresBuscarPorDocumento(ctx, docsFornecedor)
	for _, f := range forns {
		full := utils.MontarFornecedorDetalhado(ctx, repo, f)
		copyFull := *full
		dv.Vinculos = append(dv.Vinculos, Vinculo{
			Tipo:      "fornecedor",
			Descricao: descricaoFornecedor(doc, dv.Nome, &copyFull),
			Detalhes:  &VinculoDetalhes{Fornecedor: &copyFull},
		})
	}

	docsDoador := []string{doc}
	if isCPF {
		docsDoador = append(docsDoador, "000"+doc)
	}
	doadores, _ := repo.DoadoresBuscarPorDocumento(ctx, docsDoador)
	for _, d := range doadores {
		dv.Vinculos = append(dv.Vinculos, Vinculo{
			Tipo:      "doador",
			Descricao: fmt.Sprintf("Doador registrado: %s", d.Nome),
			Detalhes: &VinculoDetalhes{
				Doador: d,
			},
		})

		recsCand, _ := repo.ReceitasCandidatoBuscarPorDoadorID(ctx, d.ID)
		for _, r := range recsCand {
			dv.Vinculos = append(dv.Vinculos, Vinculo{
				Tipo:      "receita_candidato",
				Descricao: fmt.Sprintf("Doação (sq_receita=%d) a candidato", r.SQReceita),
				Detalhes: &VinculoDetalhes{
					ReceitasCandidato: []*types.ReceitaCandidato{r},
				},
			})
		}

		recsPart, _ := repo.ReceitasOrgaoBuscarPorDoadorID(ctx, d.ID)
		for _, r := range recsPart {
			dv.Vinculos = append(dv.Vinculos, Vinculo{
				Tipo:      "receita_orgao_partidario",
				Descricao: fmt.Sprintf("Doação (sq_receita=%d) a partido", r.SQReceita),
				Detalhes: &VinculoDetalhes{
					ReceitasOrgaoPartidario: []*types.ReceitaOrgaoPartidario{r},
				},
			})
		}
	}
}

func (u *AnalisarLigacaoPoliticaUseCase) buscarLigacoesPorDocumentoParcial(
	ctx context.Context,
	repo *repositorio.Repositorio,
	dv *DocumentoVinculo,
	doc string,
	nome string,
) {
	forns, _ := repo.FornecedoresBuscarPorDocumentoParcialENome(ctx, "%"+doc+"%", nome)
	for _, f := range forns {
		full := utils.MontarFornecedorDetalhado(ctx, repo, f)
		copyFull := *full
		dv.Vinculos = append(dv.Vinculos, Vinculo{
			Tipo:      "fornecedor",
			Descricao: descricaoFornecedor(dv.DocumentoNormalizado, nome, &copyFull),
			Detalhes:  &VinculoDetalhes{Fornecedor: &copyFull},
		})
	}

	doadores, _ := repo.DoadoresBuscarPorDocumentoParcial(ctx, "%"+doc+"%", nome)
	for _, d := range doadores {
		dv.Vinculos = append(dv.Vinculos, Vinculo{
			Tipo:      "doador",
			Descricao: fmt.Sprintf("Doador registrado (parcial): %s", d.Nome),
			Detalhes:  &VinculoDetalhes{Doador: d},
		})

		recsCand, _ := repo.ReceitasCandidatoBuscarPorDoadorID(ctx, d.ID)
		for _, r := range recsCand {
			dv.Vinculos = append(dv.Vinculos, Vinculo{
				Tipo:      "receita_candidato",
				Descricao: fmt.Sprintf("Doação (sq_receita=%d) a candidato", r.SQReceita),
				Detalhes:  &VinculoDetalhes{ReceitasCandidato: []*types.ReceitaCandidato{r}},
			})
		}

		recsPart, _ := repo.ReceitasOrgaoBuscarPorDoadorID(ctx, d.ID)
		for _, r := range recsPart {
			dv.Vinculos = append(dv.Vinculos, Vinculo{
				Tipo:      "receita_orgao_partidario",
				Descricao: fmt.Sprintf("Doação (sq_receita=%d) a partido", r.SQReceita),
				Detalhes:  &VinculoDetalhes{ReceitasOrgaoPartidario: []*types.ReceitaOrgaoPartidario{r}},
			})
		}
	}
}

func (u *AnalisarLigacaoPoliticaUseCase) enriquecerFornecedores(
	ctx context.Context,
	resultados []VinculoLicitacao,
) {
	log := logger.New("LigacaoPolitica: UseCase: enriquecerFornecedores")
	type fornecerRef struct {
		licIdx int
		docIdx int
		vinIdx int
		cnpj   string
	}

	var refs []fornecerRef

	for li := range resultados {
		for di := range resultados[li].Documentos {
			for vi := range resultados[li].Documentos[di].Vinculos {
				v := &resultados[li].Documentos[di].Vinculos[vi]
				if v.Tipo != "fornecedor" || v.Detalhes == nil || v.Detalhes.Fornecedor == nil {
					continue
				}
				cnpj := v.Detalhes.Fornecedor.Fornecedor.CPFCNPJ
				if len(cnpj) != 14 {
					continue
				}
				refs = append(refs, fornecerRef{
					licIdx: li, docIdx: di, vinIdx: vi, cnpj: cnpj,
				})
			}
		}
	}

	const maxConcorrencia = 5
	sem := make(chan struct{}, maxConcorrencia)
	var wg sync.WaitGroup

	for _, ref := range refs {
		wg.Add(1)
		sem <- struct{}{}
		go func(r fornecerRef) {
			defer wg.Done()
			defer func() { <-sem }()

			fn := &resultados[r.licIdx].Documentos[r.docIdx].Vinculos[r.vinIdx].Detalhes.Fornecedor.Fornecedor
			resp, err := u.opencnpj.Buscar(ctx, fn.CPFCNPJ)
			if err != nil {
				log.Error("erro ao enriquecer CNPJ", "cnpj", fn.CPFCNPJ, "erro", err)
				return
			}
			fn.Enriquecimento = utils.BuildFornecedorDTO(resp)
		}(ref)
	}
	wg.Wait()
}

func (u *AnalisarLigacaoPoliticaUseCase) enriquecerComTCU(
	ctx context.Context,
	resultados []VinculoLicitacao,
) {
	type docRef struct {
		licIdx int
		docIdx int
		doc    string
	}

	var refs []docRef

	for li := range resultados {
		for di := range resultados[li].Documentos {
			doc := resultados[li].Documentos[di].DocumentoNormalizado
			if len(doc) < 3 {
				continue
			}
			refs = append(refs, docRef{licIdx: li, docIdx: di, doc: doc})

			for vi := range resultados[li].Documentos[di].Vinculos {
				v := &resultados[li].Documentos[di].Vinculos[vi]
				if v.Detalhes == nil {
					continue
				}
				if v.Detalhes.Fornecedor != nil {
					doc2 := v.Detalhes.Fornecedor.Fornecedor.CPFCNPJ
					if doc2 != "" {
						refs = append(refs, docRef{licIdx: li, docIdx: di, doc: doc2})
					}
				}
				if v.Detalhes.Doador != nil {
					doc2 := v.Detalhes.Doador.CPFCNPJ
					if doc2 != "" {
						refs = append(refs, docRef{licIdx: li, docIdx: di, doc: doc2})
					}
				}
			}
		}
	}

	if len(refs) == 0 {
		return
	}

	seen := make(map[string]bool)
	uniqueRefs := make([]docRef, 0, len(refs))
	for _, ref := range refs {
		key := ref.doc
		if !seen[key] {
			seen[key] = true
			uniqueRefs = append(uniqueRefs, ref)
		}
	}

	const maxConcorrencia = 5
	sem := make(chan struct{}, maxConcorrencia)
	var wg sync.WaitGroup

	for _, ref := range uniqueRefs {
		wg.Add(1)
		sem <- struct{}{}
		go func(r docRef) {
			defer wg.Done()
			defer func() { <-sem }()

			dv := &resultados[r.licIdx].Documentos[r.docIdx]
			isCPF := len(r.doc) == 11
			filterCPF := tcu.TCUQueryParams{CPF: r.doc}
			filterCNPJ := tcu.TCUQueryParams{CNPJ: r.doc}

			if resp, err := u.tcu.BuscarContasIrregulares(ctx, filterCPF); err == nil && len(resp) > 0 {
				dv.Vinculos = append(dv.Vinculos, Vinculo{
					Tipo:      "tcu_contas_irregulares",
					Descricao: fmt.Sprintf("%d registro(s) de contas julgadas irregulares no TCU", len(resp)),
					Detalhes: &VinculoDetalhes{
						ContasIrregulares: resp,
					},
				})
			}

			if respCNPJ, err := u.tcu.BuscarContasIrregulares(ctx, filterCNPJ); err == nil && len(respCNPJ) > 0 {
				dv.Vinculos = append(dv.Vinculos, Vinculo{
					Tipo:      "tcu_contas_irregulares",
					Descricao: fmt.Sprintf("%d registro(s) de contas julgadas irregulares no TCU", len(respCNPJ)),
					Detalhes: &VinculoDetalhes{
						ContasIrregulares: respCNPJ,
					},
				})
			}

			if isCPF {
				if resp, err := u.tcu.BuscarInabilitados(ctx, filterCPF); err == nil && len(resp) > 0 {
					dv.Vinculos = append(dv.Vinculos, Vinculo{
						Tipo:      "tcu_inabilitado",
						Descricao: fmt.Sprintf("%d registro(s) de inabilitado no TCU", len(resp)),
						Detalhes: &VinculoDetalhes{
							Inabilitados: resp,
						},
					})
				}
			}

			if resp, err := u.tcu.BuscarInidoneos(ctx, filterCPF); err == nil && len(resp) > 0 {
				dv.Vinculos = append(dv.Vinculos, Vinculo{
					Tipo:      "tcu_inidoneo",
					Descricao: fmt.Sprintf("%d registro(s) de inidôneo no TCU", len(resp)),
					Detalhes: &VinculoDetalhes{
						Inidoneos: resp,
					},
				})
			}

			if !isCPF {
				if respCNPJ, err := u.tcu.BuscarInidoneos(ctx, filterCNPJ); err == nil && len(respCNPJ) > 0 {
					dv.Vinculos = append(dv.Vinculos, Vinculo{
						Tipo:      "tcu_inidoneo",
						Descricao: fmt.Sprintf("%d registro(s) de inidôneo no TCU", len(respCNPJ)),
						Detalhes: &VinculoDetalhes{
							Inidoneos: respCNPJ,
						},
					})
				}
			}
		}(ref)
	}
	wg.Wait()
}

func (u *AnalisarLigacaoPoliticaUseCase) enriquecerComServidorPublico(
	ctx context.Context,
	resultados []VinculoLicitacao,
) {
	if u.portaltransparencia == nil {
		return
	}

	type docRef struct {
		licIdx int
		docIdx int
		doc    string
	}

	var refs []docRef

	for li := range resultados {
		for di := range resultados[li].Documentos {
			doc := resultados[li].Documentos[di].DocumentoNormalizado
			if len(doc) != 11 {
				continue
			}
			refs = append(refs, docRef{licIdx: li, docIdx: di, doc: doc})

			for vi := range resultados[li].Documentos[di].Vinculos {
				v := &resultados[li].Documentos[di].Vinculos[vi]
				if v.Detalhes == nil {
					continue
				}
				if v.Detalhes.Fornecedor != nil {
					doc2 := v.Detalhes.Fornecedor.Fornecedor.CPFCNPJ
					if len(doc2) == 11 {
						refs = append(refs, docRef{licIdx: li, docIdx: di, doc: doc2})
					}
				}
				if v.Detalhes.Doador != nil {
					doc2 := v.Detalhes.Doador.CPFCNPJ
					if len(doc2) == 11 {
						refs = append(refs, docRef{licIdx: li, docIdx: di, doc: doc2})
					}
				}
			}
		}
	}

	if len(refs) == 0 {
		return
	}

	seen := make(map[string]bool)
	uniqueRefs := make([]docRef, 0, len(refs))
	for _, ref := range refs {
		key := ref.doc
		if !seen[key] {
			seen[key] = true
			uniqueRefs = append(uniqueRefs, ref)
		}
	}

	const maxConcorrencia = 5
	sem := make(chan struct{}, maxConcorrencia)
	var wg sync.WaitGroup

	for _, ref := range uniqueRefs {
		wg.Add(1)
		sem <- struct{}{}
		go func(r docRef) {
			defer wg.Done()
			defer func() { <-sem }()

			servidores, err := u.portaltransparencia.ListarServidores(ctx, portaltransparencia.ServidorQueryParams{
				Pagina: 1,
				CPF:    r.doc,
			})
			if err != nil || len(servidores) == 0 {
				return
			}

			dv := &resultados[r.licIdx].Documentos[r.docIdx]
			dv.Vinculos = append(dv.Vinculos, Vinculo{
				Tipo:      "servidor_publico",
				Descricao: fmt.Sprintf("%d registro(s) de servidor público no Portal da Transparência", len(servidores)),
				Detalhes: &VinculoDetalhes{
					ServidoresPublicos: servidores,
				},
			})
		}(ref)
	}
	wg.Wait()
}

func descricaoFornecedor(doc, nome string, dto *tsetypes.FornecedorDetalhado) string {
	label := doc
	if nome != "" {
		label = nome + " (" + doc + ")"
	}

	base := fmt.Sprintf("%s é fornecedor de campanha", label)

	qtdDespCand := len(dto.DespesasCandidato)
	qtdDespPart := len(dto.DespesasOrgaoPartidario)
	if qtdDespCand > 0 || qtdDespPart > 0 {
		partes := make([]string, 0)
		if qtdDespCand > 0 {
			partes = append(partes, fmt.Sprintf("%d despesa(s) de candidato", qtdDespCand))
		}
		if qtdDespPart > 0 {
			partes = append(partes, fmt.Sprintf("%d despesa(s) partidária(s)", qtdDespPart))
		}
		base += " com " + strings.Join(partes, " e ")
	}

	return base
}
