package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/cubenet-cms/cms/pkg/ws"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSHandler struct {
	hub      *ws.Hub
	timeout  time.Duration
	interval time.Duration
	maxSize  int64
}

func NewWSHandler(hub *ws.Hub, timeout, interval time.Duration, maxSize int64) *WSHandler {
	return &WSHandler{hub: hub, timeout: timeout, interval: interval, maxSize: maxSize}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade failed", "error", err)
		return
	}

	client := &ws.Client{
		ID:   r.RemoteAddr,
		Send: make(chan ws.Message, 64),
		Hub:  h.hub,
	}

	h.hub.Register(client)
	go writePump(client, conn, h.timeout, h.interval)
	go readPump(client, conn, h.timeout, h.maxSize)
}

func readPump(client *ws.Client, conn *websocket.Conn, timeout time.Duration, maxSize int64) {
	defer func() {
		client.Hub.Unregister(client)
		conn.Close()
	}()

	conn.SetReadLimit(maxSize)
	conn.SetReadDeadline(time.Now().Add(timeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(timeout))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg ws.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		slog.Debug("ws message", "type", msg.Type, "client_id", client.ID)
	}
}

func writePump(client *ws.Client, conn *websocket.Conn, timeout, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(timeout))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(timeout))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
