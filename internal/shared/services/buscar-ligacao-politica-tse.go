package services

import (
	"context"
	"fmt"
	"strings"

	repositorio "github.com/danyele/podp/internal/esferas-brasileiras/tse/repositorio"
	tsetypes "github.com/danyele/podp/internal/esferas-brasileiras/tse/types"
	"github.com/danyele/podp/internal/shared/domain"
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
		recsCand, _ := repo.ReceitasCandidatoBuscarPorDoadorID(ctx, d.ID)
		recsPart, _ := repo.ReceitasOrgaoBuscarPorDoadorID(ctx, d.ID)

		detalhes := &domain.VinculoDetalhes{Doador: d}

		for _, r := range recsCand {
			det := utils.MontarReceitaCandidatoDetalhada(ctx, repo, r)
			detalhes.ReceitasCandidato = append(detalhes.ReceitasCandidato, &det)
		}
		for _, r := range recsPart {
			det := utils.MontarReceitaOrgaoPartidarioDetalhada(ctx, repo, r)
			detalhes.ReceitasOrgaoPartidario = append(detalhes.ReceitasOrgaoPartidario, &det)
		}

		vinculos = append(vinculos, domain.Vinculo{
			Tipo:      "doador",
			Descricao: descricaoDoador(input.Nome, doc, d.Nome, detalhes),
			Detalhes:  detalhes,
		})
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
		recsCand, _ := repo.ReceitasCandidatoBuscarPorDoadorID(ctx, d.ID)
		recsPart, _ := repo.ReceitasOrgaoBuscarPorDoadorID(ctx, d.ID)

		detalhes := &domain.VinculoDetalhes{Doador: d}

		for _, r := range recsCand {
			det := utils.MontarReceitaCandidatoDetalhada(ctx, repo, r)
			detalhes.ReceitasCandidato = append(detalhes.ReceitasCandidato, &det)
		}
		for _, r := range recsPart {
			det := utils.MontarReceitaOrgaoPartidarioDetalhada(ctx, repo, r)
			detalhes.ReceitasOrgaoPartidario = append(detalhes.ReceitasOrgaoPartidario, &det)
		}

		vinculos = append(vinculos, domain.Vinculo{
			Tipo:      "doador",
			Descricao: descricaoDoador(nome, doc, d.Nome, detalhes),
			Detalhes:  detalhes,
		})
	}

	return &BuscarLigacaoPoliticaTSEOutput{Vinculos: vinculos}, nil
}

func descricaoFornecedor(doc, nome string, dto *tsetypes.FornecedorDetalhado) string {
	label := doc
	if nome != "" {
		label = nome + " (" + doc + ")"
	}

	base := fmt.Sprintf("%s é fornecedor em campanha política", label)

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

func descricaoDoador(nome, doc, nomeDoador string, detalhes *domain.VinculoDetalhes) string {
	label := nomeDoador
	if nome != "" && nome != nomeDoador {
		label = nome + " (" + doc + ")"
	} else if doc != "" {
		label = nomeDoador + " (" + doc + ")"
	}

	qtdCand := len(detalhes.ReceitasCandidato)
	qtdPart := len(detalhes.ReceitasOrgaoPartidario)

	if qtdCand == 0 && qtdPart == 0 {
		return fmt.Sprintf("Doador registrado: %s", label)
	}

	partes := make([]string, 0)
	if qtdCand > 0 {
		partes = append(partes, fmt.Sprintf("%d doação(ões) a candidato(s)", qtdCand))
	}
	if qtdPart > 0 {
		partes = append(partes, fmt.Sprintf("%d doação(ões) a partido(s)", qtdPart))
	}
	return fmt.Sprintf("%s é doador eleitoral com %s", label, strings.Join(partes, " e "))
}
