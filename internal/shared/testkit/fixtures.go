package testkit

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var TabelasParaLimpar = []string{
	"fornecedor", "doador", "candidato", "eleicao", "partido",
	"unidade_eleitoral", "prestacao_contas", "receita_candidato",
	"receita_orgao_partidario", "despesa_candidato",
	"despesa_orgao_partidario", "bem_candidato",
	"receita_doador_originario_candidato", "receita_doador_originario_orgao_partidario",
}

func InsertFornecedor(t *testing.T, ctx context.Context, pool *pgxpool.Pool, cpfCnpj, nome string) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		INSERT INTO fornecedor (id, cpf_cnpj, nome, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
		ON CONFLICT (cpf_cnpj) DO NOTHING
	`, cpfCnpj, nome)
	require.NoError(t, err)
}

func InsertDoador(t *testing.T, ctx context.Context, pool *pgxpool.Pool, cpfCnpj, nome string) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		INSERT INTO doador (id, cpf_cnpj, nome, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
		ON CONFLICT (cpf_cnpj) DO NOTHING
	`, cpfCnpj, nome)
	require.NoError(t, err)
}

func CleanTables(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tables ...string) {
	t.Helper()
	for _, tname := range tables {
		_, err := pool.Exec(ctx, "DELETE FROM "+tname)
		require.NoError(t, err, "falha ao limpar tabela %s", tname)
	}
}

func CleanAllTables(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	CleanTables(t, ctx, pool, TabelasParaLimpar...)
}
