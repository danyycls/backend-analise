package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danyele/podp/internal/app"
	database "github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/logger"
)

func main() {
	log := logger.New("Main")
	log.Info("build info", "version", "dev", "commit", "none", "data", "2026-06-28")

	ctx, cancelar := context.WithTimeout(context.Background(), 30*time.Minute)

	pool, err := database.NovaPool(ctx, database.ConfigFromEnv())
	if err != nil {
		log.Fatal("erro ao criar pgx pool", "erro", err)
	}
	defer cancelar()
	defer pool.Close()

	poolDB := database.NewPoolDB(pool)

	a := app.NovoApp(poolDB, obterDiretorioCSV())

	roteador := app.NovoRoteador(a)

	handlerComCORS := corsHandler(roteador)
	endereco := ":" + obterPorta()
	log.Info("API iniciada", "endereco", "http://localhost"+endereco)
	if err := http.ListenAndServe(endereco, handlerComCORS); err != nil {
		log.Fatal("erro ao iniciar API", "erro", err)
	}
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func init() {
	if _, existe := os.LookupEnv("DB_HOST"); !existe {
		_ = carregarArquivoEnv(".env")
	}
}

func carregarArquivoEnv(caminho string) error {
	conteudo, err := os.ReadFile(caminho)
	if err != nil {
		return err
	}

	for _, linha := range strings.Split(string(conteudo), "\n") {
		linha = strings.TrimSpace(linha)
		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}

		partes := strings.SplitN(linha, "=", 2)
		if len(partes) != 2 {
			continue
		}

		chave := strings.TrimSpace(partes[0])
		valor := strings.TrimSpace(partes[1])
		if _, existe := os.LookupEnv(chave); !existe {
			_ = os.Setenv(chave, valor)
		}
	}

	return nil
}

func obterPorta() string {
	return strings.TrimSpace(os.Getenv("PORT"))
}

func obterDiretorioCSV() string {
	return strings.TrimSpace(os.Getenv("IMPORTACAO_DIRETORIO_CSV"))
}
