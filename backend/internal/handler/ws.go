package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/ws"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSHandler struct {
	hub  *ws.Hub
	cfg  config.WSConfig
}

func NewWSHandler(hub *ws.Hub, cfg config.WSConfig) *WSHandler {
	return &WSHandler{hub: hub, cfg: cfg}
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
	go h.writePump(client, conn)
	go h.readPump(client, conn)
}

func (h *WSHandler) readPump(client *ws.Client, conn *websocket.Conn) {
	defer func() {
		h.hub.Unregister(client)
		conn.Close()
	}()

	conn.SetReadLimit(h.cfg.MaxMessageSz)
	conn.SetReadDeadline(time.Now().Add(h.cfg.ReadTimeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(h.cfg.ReadTimeout))
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

		slog.Debug("ws message received", "type", msg.Type, "client_id", client.ID)
	}
}

func (h *WSHandler) writePump(client *ws.Client, conn *websocket.Conn) {
	ticker := time.NewTicker(h.cfg.PingInterval)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(h.cfg.WriteTimeout))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteJSON(msg); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(h.cfg.WriteTimeout))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
