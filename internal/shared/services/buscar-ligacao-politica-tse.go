package services

import (
	"context"
	"fmt"
	"strings"

	repositorio "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/shared/domain"
	"github.com/danyele/podp/internal/shared/types"
	"github.com/danyele/podp/internal/shared/utils"
)

type BuscarLigacaoPoliticaTSEInput struct {
	DocumentoNormalizado string
	Nome                 string
	Parcial              bool
}

type BuscarLigacaoPoliticaTSEOutput struct {
	Vinculos []domain.Vinculo
}

type BuscarLigacaoPoliticaTSEService interface {
	Executar(ctx context.Context, repo *repositorio.Repositorio, input BuscarLigacaoPoliticaTSEInput) (*BuscarLigacaoPoliticaTSEOutput, error)
}

type buscarLigacaoPoliticaTSEServiceImpl struct{}

func NovoBuscarLigacaoPoliticaTSEService() BuscarLigacaoPoliticaTSEService {
	return &buscarLigacaoPoliticaTSEServiceImpl{}
}

func (s *buscarLigacaoPoliticaTSEServiceImpl) Executar(
	ctx context.Context,
	repo *repositorio.Repositorio,
	input BuscarLigacaoPoliticaTSEInput,
) (*BuscarLigacaoPoliticaTSEOutput, error) {
	if input.Parcial && input.Nome != "" {
		return s.buscarParcial(ctx, repo, input)
	}
	return s.buscarExato(ctx, repo, input)
}

func (s *buscarLigacaoPoliticaTSEServiceImpl) buscarExato(
	ctx context.Context,
	repo *repositorio.Repositorio,
	input BuscarLigacaoPoliticaTSEInput,
) (*BuscarLigacaoPoliticaTSEOutput, error) {
	doc := input.DocumentoNormalizado
	isCPF := len(doc) == 11
	var vinculos []domain.Vinculo

	docsFornecedor := []string{doc}
	if isCPF {
		docsFornecedor = append(docsFornecedor, "000"+doc)
	}
	forns, _ := repo.FornecedoresBuscarPorDocumento(ctx, docsFornecedor)
	for _, f := range forns {
		full := utils.MontarFornecedorDetalhado(ctx, repo, f)
		copyFull := *full
		vinculos = append(vinculos, domain.Vinculo{
			Tipo:      "fornecedor",
			Descricao: descricaoFornecedor(doc, input.Nome, &copyFull),
			Detalhes:  &domain.VinculoDetalhes{Fornecedor: &copyFull},
		})
	}

	docsDoador := []string{doc}
	if isCPF {
		docsDoador = append(docsDoador, "000"+doc)
	}
	doadores, _ := repo.DoadoresBuscarPorDocumento(ctx, docsDoador)
	for _, d := range doadores {
		vinculos = append(vinculos, domain.Vinculo{
			Tipo:      "doador",
			Descricao: fmt.Sprintf("Doador registrado: %s", d.Nome),
			Detalhes: &domain.VinculoDetalhes{
				Doador: d,
			},
		})

		recsCand, _ := repo.ReceitasCandidatoBuscarPorDoadorID(ctx, d.ID)
		for _, r := range recsCand {
			vinculos = append(vinculos, domain.Vinculo{
				Tipo:      "receita_candidato",
				Descricao: fmt.Sprintf("Doação (sq_receita=%d) a candidato", r.SQReceita),
				Detalhes: &domain.VinculoDetalhes{
					ReceitasCandidato: []*types.ReceitaCandidato{r},
				},
			})
		}

		recsPart, _ := repo.ReceitasOrgaoBuscarPorDoadorID(ctx, d.ID)
		for _, r := range recsPart {
			vinculos = append(vinculos, domain.Vinculo{
				Tipo:      "receita_orgao_partidario",
				Descricao: fmt.Sprintf("Doação (sq_receita=%d) a partido", r.SQReceita),
				Detalhes: &domain.VinculoDetalhes{
					ReceitasOrgaoPartidario: []*types.ReceitaOrgaoPartidario{r},
				},
			})
		}
	}

	return &BuscarLigacaoPoliticaTSEOutput{Vinculos: vinculos}, nil
}

func (s *buscarLigacaoPoliticaTSEServiceImpl) buscarParcial(
	ctx context.Context,
	repo *repositorio.Repositorio,
	input BuscarLigacaoPoliticaTSEInput,
) (*BuscarLigacaoPoliticaTSEOutput, error) {
	doc := input.DocumentoNormalizado
	nome := input.Nome
	var vinculos []domain.Vinculo

	forns, _ := repo.FornecedoresBuscarPorDocumentoParcialENome(ctx, "%"+doc+"%", nome)
	for _, f := range forns {
		full := utils.MontarFornecedorDetalhado(ctx, repo, f)
		copyFull := *full
		vinculos = append(vinculos, domain.Vinculo{
			Tipo:      "fornecedor",
			Descricao: descricaoFornecedor(doc, nome, &copyFull),
			Detalhes:  &domain.VinculoDetalhes{Fornecedor: &copyFull},
		})
	}

	doadores, _ := repo.DoadoresBuscarPorDocumentoParcial(ctx, "%"+doc+"%", nome)
	for _, d := range doadores {
		vinculos = append(vinculos, domain.Vinculo{
			Tipo:      "doador",
			Descricao: fmt.Sprintf("Doador registrado (parcial): %s", d.Nome),
			Detalhes:  &domain.VinculoDetalhes{Doador: d},
		})

		recsCand, _ := repo.ReceitasCandidatoBuscarPorDoadorID(ctx, d.ID)
		for _, r := range recsCand {
			vinculos = append(vinculos, domain.Vinculo{
				Tipo:      "receita_candidato",
				Descricao: fmt.Sprintf("Doação (sq_receita=%d) a candidato", r.SQReceita),
				Detalhes: &domain.VinculoDetalhes{
					ReceitasCandidato: []*types.ReceitaCandidato{r},
				},
			})
		}

		recsPart, _ := repo.ReceitasOrgaoBuscarPorDoadorID(ctx, d.ID)
		for _, r := range recsPart {
			vinculos = append(vinculos, domain.Vinculo{
				Tipo:      "receita_orgao_partidario",
				Descricao: fmt.Sprintf("Doação (sq_receita=%d) a partido", r.SQReceita),
				Detalhes: &domain.VinculoDetalhes{
					ReceitasOrgaoPartidario: []*types.ReceitaOrgaoPartidario{r},
				},
			})
		}
	}

	return &BuscarLigacaoPoliticaTSEOutput{Vinculos: vinculos}, nil
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
