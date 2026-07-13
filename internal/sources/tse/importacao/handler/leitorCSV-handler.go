// Pacote handler — endpoint para importacao de dados CSV (carga inicial do banco)
package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
	"github.com/danyele/podp/internal/sources/tse/importacao/usecase"

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
	inicio := time.Now()
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	done := make(chan struct{})
	var resultado *usecase.ImportarCSVResponse
	var err error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic na goroutine de importacao", "panic", r)
				err = fmt.Errorf("panic: %v", r)
			}
			close(done)
		}()
		resultado, err = h.useCase.Executar(ctx, usecase.ImportarCSVRequest{})
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			log.Warn("cliente desconectado, cancelando importacao", "duracao_minutos", time.Since(inicio).Minutes())
			cancel()
			return
		case <-done:
			duracaoMinutos := time.Since(inicio).Minutes()
			if err != nil {
				log.Error("erro na importacao", "erro", err, "duracao_minutos", duracaoMinutos)
				c.SSEvent("erro", gin.H{
					"sucesso":         0,
					"timestamp":       time.Now().Format(time.RFC3339Nano),
					"duracao_minutos": duracaoMinutos,
				})
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
					"timestamp":            time.Now().Format(time.RFC3339Nano),
					"duracao_minutos":      duracaoMinutos,
				})
			}
			c.Writer.Flush()
			return
		case <-ticker.C:
			progression := h.useCase.ProgressoEvento()
			progression.Timestamp = time.Now().Format(time.RFC3339Nano)
			c.SSEvent("progression", progression)
			c.Writer.Flush()
		}
	}
}
