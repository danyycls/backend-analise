package parse

import (
	"context"
	"fmt"
	"time"

	"github.com/danyele/podp/internal/shared/logger"

	repositorios "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/repositorios"
	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
	"github.com/jackc/pgx/v5"
)

func lockDimensao(ctx context.Context, tx pgx.Tx, nome string) error {
	sql := `SELECT pg_advisory_xact_lock(hashtext($1))`
	_, err := tx.Exec(ctx, sql, nome)
	if err != nil {
		return fmt.Errorf("lock %s: %w", nome, err)
	}
	return nil
}

func PersistirDadosImportacaoPgCopy(
	ctx context.Context,
	tx pgx.Tx,
	repo *repositorios.Repositorio,
	dados *tipos.DadosImportacao,
	lote int,
) error {
	log := logger.New("LeitorCSV: Utils: PersistirDadosImportacaoPgCopy")
	inicio := time.Now()
	log.Info("iniciando persistencia",
		"passo", "1/15", "entidade", "Convenios", "registros", len(dados.Convenios))
	if _, err := repo.InserirEmLote(ctx, tx, dados.Convenios, lote); err != nil {
		return fmt.Errorf("convenio: %w", err)
	}
	log.Info("persistencia concluida",
		"passo", "1/15", "entidade", "Convenios", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "2/15", "entidade", "Eleicoes", "registros", len(dados.Eleicoes))
	if err := lockDimensao(ctx, tx, "eleicao"); err != nil {
		return err
	}
	mapeamento, err := repo.InserirEleicoesComRetorno(ctx, tx, valores(dados.Eleicoes), lote)
	if err != nil {
		return fmt.Errorf("eleicao: %w", err)
	}
	remapearEleicaoIDs(dados, mapeamento)
	sincronizarDependenciasDeEleicao(dados)
	log.Info("persistencia concluida",
		"passo", "2/15", "entidade", "Eleicoes", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "3/15", "entidade", "UnidadesEleitorais", "registros", len(dados.UnidadesEleitorais))
	if err := lockDimensao(ctx, tx, "unidade_eleitoral"); err != nil {
		return err
	}
	mapeamento, err = repo.InserirUnidadesEleitoraisComRetorno(ctx, tx, valores(dados.UnidadesEleitorais), lote)
	if err != nil {
		return fmt.Errorf("unidade_eleitoral: %w", err)
	}
	remapearUnidadeEleitoralIDs(dados, mapeamento)
	sincronizarDependenciasDeUnidadeEleitoral(dados)
	log.Info("persistencia concluida",
		"passo", "3/15", "entidade", "UnidadesEleitorais", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "4/15", "entidade", "Partidos", "registros", len(dados.Partidos))
	if err := lockDimensao(ctx, tx, "partido"); err != nil {
		return err
	}
	mapeamento, err = repo.InserirPartidosComRetorno(ctx, tx, valores(dados.Partidos), lote)
	if err != nil {
		return fmt.Errorf("partido: %w", err)
	}
	remapearPartidoIDs(dados, mapeamento)
	sincronizarDependenciasDePartido(dados)
	log.Info("persistencia concluida",
		"passo", "4/15", "entidade", "Partidos", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "5/15", "entidade", "Candidatos", "registros", len(dados.Candidatos))
	mapeamento, err = repo.InserirCandidatosComRetorno(ctx, tx, valores(dados.Candidatos), lote)
	if err != nil {
		return fmt.Errorf("candidato: %w", err)
	}

	log.Info("diagnostico de mapeamento",
		"passo", "5/15", "entidade", "Candidatos", "entradas_retornadas", len(mapeamento))
	remapearCandidatoIDs(dados, mapeamento)
	sincronizarDependenciasDeCandidato(dados)
	log.Info("persistencia concluida",
		"passo", "5/15", "entidade", "Candidatos", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "6/15", "entidade", "Fornecedores", "registros", len(dados.Fornecedores))
	mapeamento, err = repo.InserirFornecedoresComRetorno(ctx, tx, valores(dados.Fornecedores), lote)
	if err != nil {
		return fmt.Errorf("fornecedor: %w", err)
	}
	remapearFornecedorIDs(dados, mapeamento)
	sincronizarDependenciasDeFornecedor(dados)
	log.Info("persistencia concluida",
		"passo", "6/15", "entidade", "Fornecedores", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "7/15", "entidade", "Doadores", "registros", len(dados.Doadores))
	mapeamento, err = repo.InserirDoadoresComRetorno(ctx, tx, valores(dados.Doadores), lote)
	if err != nil {
		return fmt.Errorf("doador: %w", err)
	}
	remapearDoadorIDs(dados, mapeamento)
	sincronizarDependenciasDeDoador(dados)
	log.Info("persistencia concluida",
		"passo", "7/15", "entidade", "Doadores", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "8/15", "entidade", "PrestacoesContas", "registros", len(dados.Prestacoes))
	mapeamento, err = repo.InserirPrestacoesComRetorno(ctx, tx, valores(dados.Prestacoes), lote)
	if err != nil {
		return fmt.Errorf("prestacao_contas: %w", err)
	}
	if err := remapearPrestacaoIDsComPlaceholderPgCopy(ctx, tx, repo, dados, mapeamento); err != nil {
		return err
	}
	sincronizarDependenciasDePrestacao(dados)
	log.Info("persistencia concluida",
		"passo", "8/15", "entidade", "PrestacoesContas", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "9/15", "entidade", "DespesasCandidato", "registros", len(dados.DespesasCandidato))
	if _, err := repo.InserirDespesasCandidato(ctx, tx, dados.DespesasCandidato, lote); err != nil {
		return fmt.Errorf("despesa_candidato: %w", err)
	}
	log.Info("persistencia concluida",
		"passo", "9/15", "entidade", "DespesasCandidato", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "10/15", "entidade", "DespesasOrgaoPartidario", "registros", len(dados.DespesasOrgaoPartidario))
	if _, err := repo.InserirDespesasOrgaoPartidario(ctx, tx, dados.DespesasOrgaoPartidario, lote); err != nil {
		return fmt.Errorf("despesa_orgao_partidario: %w", err)
	}
	log.Info("persistencia concluida",
		"passo", "10/15", "entidade", "DespesasOrgaoPartidario", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "11/15", "entidade", "ReceitasCandidato", "registros", len(dados.ReceitasCandidato))
	mapeamento, err = repo.InserirReceitasCandidatoComRetorno(ctx, tx, dados.ReceitasCandidato, lote)
	if err != nil {
		return fmt.Errorf("receita_candidato: %w", err)
	}
	remapearReceitaCandidatoIDs(dados, mapeamento)
	sincronizarDependenciasDeReceitaCandidato(dados)
	log.Info("persistencia concluida",
		"passo", "11/15", "entidade", "ReceitasCandidato", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "12/15", "entidade", "ReceitasOrgaoPartidario", "registros", len(dados.ReceitasOrgaoPartidario))
	mapeamento, err = repo.InserirReceitasOrgaoComRetorno(ctx, tx, dados.ReceitasOrgaoPartidario, lote)
	if err != nil {
		return fmt.Errorf("receita_orgao_partidario: %w", err)
	}
	remapearReceitaOrgaoPartidarioIDs(dados, mapeamento)
	sincronizarDependenciasDeReceitaOrgaoPartidario(dados)
	log.Info("persistencia concluida",
		"passo", "12/15", "entidade", "ReceitasOrgaoPartidario", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "13/15", "entidade", "ReceitasDoadorOriginarioCandidato", "registros", len(dados.ReceitasDoadorOriginarioCandidato))
	if _, err := repo.InserirEmLote(ctx, tx, dados.ReceitasDoadorOriginarioCandidato, lote); err != nil {
		return fmt.Errorf("receita_doador_originario_candidato: %w", err)
	}
	log.Info("persistencia concluida",
		"passo", "13/15", "entidade", "ReceitasDoadorOriginarioCandidato", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "14/15", "entidade", "ReceitasDoadorOriginarioOrgaoPartidario", "registros", len(dados.ReceitasDoadorOriginarioOrgaoPartidario))
	if _, err := repo.InserirEmLote(ctx, tx, dados.ReceitasDoadorOriginarioOrgaoPartidario, lote); err != nil {
		return fmt.Errorf("receita_doador_originario_orgao_partidario: %w", err)
	}
	log.Info("persistencia concluida",
		"passo", "14/15", "entidade", "ReceitasDoadorOriginarioOrgaoPartidario", "duracao", time.Since(inicio).String())

	log.Info("iniciando persistencia",
		"passo", "15/15", "entidade", "BensCandidato", "registros", len(dados.BensCandidato))
	if _, err := repo.InserirEmLote(ctx, tx, dados.BensCandidato, lote); err != nil {
		return fmt.Errorf("bem_candidato: %w", err)
	}
	log.Info("persistencia concluida",
		"passo", "15/15", "entidade", "BensCandidato", "duracao", time.Since(inicio).String())

	return nil
}
