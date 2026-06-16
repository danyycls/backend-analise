package websocket

import (
	"encoding/json"
	"net/http"

	gorilla "github.com/gorilla/websocket"
)

var Upgrader = gorilla.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WriteJSON(conn *gorilla.Conn, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return conn.WriteMessage(gorilla.TextMessage, data)
}
