package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Echo(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, []byte("Pong"))
	}
}
