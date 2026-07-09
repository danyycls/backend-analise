package main

import (
	"context"
	"fmt"
	"time"

	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/logger"
	migracao "github.com/danyele/podp/internal/shared/migrations"
)

func main() {
	log := logger.New("Migrate")
	if err := run(); err != nil {
		log.Fatal("erro fatal", "erro", err)
	}
}

func run() error {
	log := logger.New("Migrate")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pool, err := database.NovaPool(ctx, database.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("erro ao criar pool: %w", err)
	}
	defer pool.Close()

	if err := migracao.AplicarSQLPool(ctx, pool, "internal/shared/migrations/schema"); err != nil {
		return fmt.Errorf("erro ao aplicar migrations: %w", err)
	}
	log.Info("migrations aplicadas com sucesso")
	return nil
}
