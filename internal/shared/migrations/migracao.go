package migracao

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const tabelaMigracoesProjeto = "liceu_schema_migrations"

func AplicarSQLPool(ctx context.Context, pool *pgxpool.Pool, diretorio string) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if err := garantirTabelaMigracoesPool(ctx, conn); err != nil {
		return err
	}

	entradas, err := os.ReadDir(diretorio)
	if err != nil {
		return fmt.Errorf("ler diretorio migrations: %w", err)
	}

	var arquivos []string
	for _, e := range entradas {
		if e.IsDir() {
			continue
		}
		nome := e.Name()
		if strings.HasSuffix(nome, ".up.sql") {
			arquivos = append(arquivos, nome)
		}
	}
	sort.Strings(arquivos)

	if err := registrarBaselineSeNecessarioPool(ctx, conn, arquivos); err != nil {
		return err
	}

	for _, nome := range arquivos {
		aplicada, err := migrationJaAplicadaPool(ctx, conn, nome)
		if err != nil {
			return err
		}
		if aplicada {
			continue
		}

		caminho := filepath.Join(diretorio, nome)
		sqlBytes, err := os.ReadFile(caminho)
		if err != nil {
			return fmt.Errorf("ler %s: %w", nome, err)
		}
		if _, err := conn.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("executar %s: %w", nome, err)
		}
		if err := registrarMigrationPool(ctx, conn, nome); err != nil {
			return err
		}
	}
	return nil
}

func garantirTabelaMigracoesPool(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            nome VARCHAR(255) PRIMARY KEY,
            aplicada_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )
    `, tabelaMigracoesProjeto))
	return err
}

func migrationJaAplicadaPool(ctx context.Context, conn *pgxpool.Conn, nome string) (bool, error) {
	var total int64
	row := conn.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE nome = $1`, tabelaMigracoesProjeto), nome)
	if err := row.Scan(&total); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("verificar migration %s: %w", nome, err)
	}
	return total > 0, nil
}

func registrarMigrationPool(ctx context.Context, conn *pgxpool.Conn, nome string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (nome) VALUES ($1)`, tabelaMigracoesProjeto), nome)
	return err
}

func registrarBaselineSeNecessarioPool(ctx context.Context, conn *pgxpool.Conn, arquivos []string) error {
	if len(arquivos) == 0 {
		return nil
	}

	var existeEleicao int64
	row := conn.QueryRow(ctx, `
        SELECT COUNT(*) FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'eleicao'
    `)
	if err := row.Scan(&existeEleicao); err != nil {
		return fmt.Errorf("verificar baseline do schema: %w", err)
	}
	if existeEleicao == 0 {
		return nil
	}

	var totalMigracoes int64
	row2 := conn.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, tabelaMigracoesProjeto))
	if err := row2.Scan(&totalMigracoes); err != nil {
		return fmt.Errorf("verificar historico de migrations: %w", err)
	}
	if totalMigracoes > 0 {
		return nil
	}

	for _, nome := range arquivos {
		if strings.HasPrefix(nome, "000001_") {
			return registrarMigrationPool(ctx, conn, nome)
		}
	}
	return nil
}
