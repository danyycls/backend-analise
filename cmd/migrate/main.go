package main

import (
	"context"
	"time"

	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/logger"
	migracao "github.com/danyele/podp/internal/shared/migrations"
)

func main() {
	log := logger.New("Migrate")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pool, err := database.NovaPool(ctx, database.ConfigFromEnv())
	if err != nil {
		log.Fatal("erro ao criar pool", "erro", err)
	}
	defer pool.Close()

	if err := migracao.AplicarSQLPool(ctx, pool, "internal/shared/migrations/schema"); err != nil {
		log.Fatal("erro ao aplicar migrations", "erro", err)
	}
	log.Info("migrations aplicadas com sucesso")
}
