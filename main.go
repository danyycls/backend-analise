package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danyele/laceu/internal/shared/logger"

	"github.com/danyele/laceu/internal/app"
	database "github.com/danyele/laceu/internal/shared/database"
	migracao "github.com/danyele/laceu/internal/shared/migrations"
)

func main() {
	log := logger.New("Main: main: ListenAndServe")
	ctx, cancelar := context.WithTimeout(context.Background(), 30*time.Minute)

	pool, err := database.NovaPool(ctx, database.ConfigFromEnv())
	if err != nil {
		log := logger.New("Main: main: NovaPool")
		log.Fatal("erro ao criar pgx pool", "erro", err)
	}
	defer cancelar()
	defer pool.Close()

	poolDB := database.NewPoolDB(pool)

	if err := migracao.AplicarSQLPool(ctx, pool, "internal/shared/migrations/schema"); err != nil {
		logger.New("Main: main: AplicarSQLPool").Error("erro ao aplicar migrations", "erro", err)
		os.Exit(1) //nolint:gocritic
	}

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
	if strings.TrimSpace(os.Getenv("DB_HOST")) == "" {
		_ = os.Setenv("DB_HOST", "localhost")
	}
	if strings.TrimSpace(os.Getenv("DB_PORT")) == "" {
		_ = os.Setenv("DB_PORT", "5432")
	}
	if strings.TrimSpace(os.Getenv("DB_USER")) == "" {
		_ = os.Setenv("DB_USER", "postgres")
	}
	if strings.TrimSpace(os.Getenv("DB_PASSWORD")) == "" {
		_ = os.Setenv("DB_PASSWORD", "postgres")
	}
	if strings.TrimSpace(os.Getenv("DB_NAME")) == "" {
		_ = os.Setenv("DB_NAME", "tse_data")
	}
	if strings.TrimSpace(os.Getenv("PORTAL_TRANSPARENCIA_BASE_URL")) == "" {
		_ = os.Setenv("PORTAL_TRANSPARENCIA_BASE_URL", "https://api.portaldatransparencia.gov.br")
	}
	if strings.TrimSpace(os.Getenv("PNCP_BASE_URL")) == "" {
		_ = os.Setenv("PNCP_BASE_URL", "https://pncp.gov.br/pncp-consulta/v1")
	}
	if strings.TrimSpace(os.Getenv("TCU_BASE_URL")) == "" {
		_ = os.Setenv("TCU_BASE_URL", "https://certidoes.apps.gov.br/api/publico")
	}
	if strings.TrimSpace(os.Getenv("SICONFI_BASE_URL")) == "" {
		_ = os.Setenv("SICONFI_BASE_URL", "https://apidatalake.tesouro.gov.br/ords/siconfi/tt")
	}
	if strings.TrimSpace(os.Getenv("DEPUTADOS_BASE_URL")) == "" {
		_ = os.Setenv("DEPUTADOS_BASE_URL", "https://dadosabertos.camara.leg.br/api/v2")
	}
	if strings.TrimSpace(os.Getenv("IBGE_BASE_URL")) == "" {
		_ = os.Setenv("IBGE_BASE_URL", "https://servicodados.ibge.gov.br/api/v1/localidades")
	}
	if strings.TrimSpace(os.Getenv("IBGE_AGREGADOS_BASE_URL")) == "" {
		_ = os.Setenv("IBGE_AGREGADOS_BASE_URL", "https://servicodados.ibge.gov.br/api/v3/agregados")
	}
	if strings.TrimSpace(os.Getenv("OPENCNPJ_BASE_URL")) == "" {
		_ = os.Setenv("OPENCNPJ_BASE_URL", "https://api.opencnpj.org/%s")
	}
	if strings.TrimSpace(os.Getenv("SENADO_BASE_URL")) == "" {
		_ = os.Setenv("SENADO_BASE_URL", "https://legis.senado.leg.br/dadosabertos")
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
	if porta := strings.TrimSpace(os.Getenv("PORT")); porta != "" {
		return porta
	}
	return "8080"
}

func obterDiretorioCSV() string {
	if diretorio := strings.TrimSpace(os.Getenv("IMPORTACAO_DIRETORIO_CSV")); diretorio != "" {
		return diretorio
	}
	return "dataCSV"
}
