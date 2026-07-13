package stream

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"

	ws "github.com/danyele/podp/internal/shared/websocket"
	pncp "github.com/danyele/podp/internal/sources/pncp/client"
	handlerPNCP "github.com/danyele/podp/internal/sources/pncp/handler"
)

type StreamMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ClientMessage struct {
	Channel    string `json:"channel"`
	JobID      string `json:"job_id,omitempty"`
	UF         string `json:"uf,omitempty"`
	CodigoIBGE int    `json:"codigo_ibge,omitempty"`
	Exercicio  int64  `json:"exercicio,omitempty"`
}

type Hub struct {
	orgaoHandler              *handlerPNCP.AnaliseOrgaoPNCPHandler
	analiseUFMunicipioHandler *handlerPNCP.AnaliseUFMunicipioHandler
}

func NewHub(
	orgaoHandler *handlerPNCP.AnaliseOrgaoPNCPHandler,
	analiseUFMunicipioHandler *handlerPNCP.AnaliseUFMunicipioHandler,
) *Hub {
	return &Hub{
		orgaoHandler:              orgaoHandler,
		analiseUFMunicipioHandler: analiseUFMunicipioHandler,
	}
}

func (h *Hub) Handle(c *gin.Context) {
	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		return
	}

	var msg ClientMessage
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return
	}

	ctx := c.Request.Context()

	switch msg.Channel {
	case "orgao_analise":
		h.streamOrgao(ctx, conn, msg.JobID)
	case "uf_municipio_analise":
		h.streamUFMunicipio(ctx, conn, msg.JobID)
	}
}

func (h *Hub) streamOrgao(ctx context.Context, conn *gorilla.Conn, jobID string) {
	eventChan, exists := h.orgaoHandler.GetJobChan(jobID)
	if !exists {
		ws.WriteJSON(conn, pncp.EventoAnalise{Type: "error", Message: "job nao encontrado"})
		return
	}
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				ws.WriteJSON(conn, pncp.EventoAnalise{Type: "done"})
				return
			}
			if err := ws.WriteJSON(conn, event); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) streamUFMunicipio(ctx context.Context, conn *gorilla.Conn, jobID string) {
	eventChan, exists := h.analiseUFMunicipioHandler.GetJobChan(jobID)
	if !exists {
		ws.WriteJSON(conn, pncp.EventoAnalise{Type: "error", Message: "job nao encontrado"})
		return
	}
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				ws.WriteJSON(conn, pncp.EventoAnalise{Type: "done"})
				return
			}
			if err := ws.WriteJSON(conn, event); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
