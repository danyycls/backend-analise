package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	ws "github.com/danyele/laceu/internal/shared/websocket"

	pncp "github.com/danyele/laceu/internal/shared/clients/pncp"
)

type WSPublicacaoStreamHandler struct {
	mu   sync.Mutex
	jobs map[string]*pubJobState
}

func NovoWSPublicacaoStreamHandler(jobs map[string]*pubJobState) *WSPublicacaoStreamHandler {
	return &WSPublicacaoStreamHandler{
		jobs: jobs,
	}
}

func (h *WSPublicacaoStreamHandler) WSStream(c *gin.Context) {
	jobID := c.Param("jobId")

	h.mu.Lock()
	job, exists := h.jobs[jobID]
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"erro": "job nao encontrado"})
		return
	}

	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx := c.Request.Context()
	for {
		select {
		case event, ok := <-job.eventChan:
			if !ok {
				data, _ := json.Marshal(pncp.EventoAnalise{Type: "done"})
				conn.WriteMessage(1, data)
				return
			}
			ws.WriteJSON(conn, event)

		case <-ctx.Done():
			return
		}
	}
}
