// Pacote handler — endpoint para importacao de dados CSV (carga inicial do banco)
package handler

import (
	"context"
	"time"

	"github.com/danyele/laceu/internal/esferas-brasileiras/tse/importacao/usecase"
	"github.com/danyele/laceu/internal/shared/logger"

	"github.com/gin-gonic/gin"
)

type LeitorCSVHandler struct {
	useCase usecase.ImportarCSVUseCase
}

func NovoLeitorCSVHandler(useCase usecase.ImportarCSVUseCase) *LeitorCSVHandler {
	return &LeitorCSVHandler{useCase: useCase}
}

func (h *LeitorCSVHandler) Executar(c *gin.Context) {
	log := logger.New("LeitorCSV: Handler: Executar")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	done := make(chan struct{})
	var resultado *usecase.ImportarCSVResponse
	var err error

	go func() {
		resultado, err = h.useCase.Executar(context.Background(), usecase.ImportarCSVRequest{})
		close(done)
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			if err != nil {
				log.Error("erro na importacao", "erro", err)
				c.SSEvent("erro", gin.H{"sucesso": 0})
			} else {
				persistidos := 0
				totalRegistros := 0
				if resultado != nil {
					totalRegistros = resultado.TotalRegistros
					persistidos = len(resultado.ArquivosComSucesso)
				}
				c.SSEvent("concluido", gin.H{
					"sucesso":              1,
					"total_registros":      totalRegistros,
					"arquivos_persistidos": persistidos,
				})
			}
			c.Writer.Flush()
			return
		case <-ticker.C:
			c.SSEvent("progression", h.useCase.ProgressoEvento())
			c.Writer.Flush()
		}
	}
}
