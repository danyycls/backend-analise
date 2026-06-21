package testkit

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/danyele/podp/internal/shared/database"

	tctest "github.com/testcontainers/testcontainers-go"
)

type IntegrationTestCase struct {
	Name     string
	Fixtures func(ctx context.Context, t *testing.T, pool *pgxpool.Pool)
	Assert   func(t *testing.T, w *httptest.ResponseRecorder)
}

func StartPostgresContainer(ctx context.Context) (*tcpostgres.PostgresContainer, *pgxpool.Pool, database.DB, error) {
	pgc, err := tcpostgres.Run(ctx, "postgres:15-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testusr"),
		tcpostgres.WithPassword("secret"),
		tctest.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("start postgres container: %w", err)
	}

	connStr, err := pgc.ConnectionString(ctx)
	if err != nil {
		pgc.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("connection string: %w", err)
	}
	connStr += "?sslmode=disable"

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		pgc.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("create pool: %w", err)
	}

	if err := RunMigrations(ctx, pool); err != nil {
		pool.Close()
		pgc.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("run migrations: %w", err)
	}

	return pgc, pool, database.NewPoolDB(pool), nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrationsDir := migrationsDirPath()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("ler diretorio de migracoes %s: %w", migrationsDir, err)
	}

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" || !strings.HasSuffix(e.Name(), ".up.sql") {
			continue
		}
		sqlBytes, err := os.ReadFile(filepath.Join(migrationsDir, e.Name()))
		if err != nil {
			return fmt.Errorf("ler migration %s: %w", e.Name(), err)
		}
		if _, err := conn.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("executar migration %s: %w", e.Name(), err)
		}
	}

	if err := registerBaselineIfNeeded(ctx, conn, entries); err != nil {
		return err
	}

	return nil
}

func migrationsDirPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	return filepath.Join(dir, "migrations", "schema")
}

func registerBaselineIfNeeded(ctx context.Context, conn *pgxpool.Conn, entries []os.DirEntry) error {
	var existe int64
	row := conn.QueryRow(ctx, `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'podp_schema_migrations'
	`)
	if err := row.Scan(&existe); err != nil {
		return err
	}
	if existe == 0 {
		_, err := conn.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS podp_schema_migrations (
				nome VARCHAR(255) PRIMARY KEY,
				aplicada_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)
		`)
		return err
	}

	var total int64
	row2 := conn.QueryRow(ctx, `SELECT COUNT(*) FROM podp_schema_migrations`)
	if err := row2.Scan(&total); err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" && strings.HasSuffix(e.Name(), ".up.sql") {
			_, err := conn.Exec(ctx, `INSERT INTO podp_schema_migrations (nome) VALUES ($1)`, e.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
