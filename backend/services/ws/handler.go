package wssvc

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/cubenet-cms/backend/internal/ws"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Config struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PingInterval time.Duration
	MaxMessageSz int64
}

type Handler struct {
	hub *ws.Hub
	cfg Config
}

func NewHandler(hub *ws.Hub, cfg Config) *Handler {
	return &Handler{hub: hub, cfg: cfg}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	go writePump(client, conn, h.cfg)
	go readPump(client, conn, h.cfg)
}

func readPump(client *ws.Client, conn *websocket.Conn, cfg Config) {
	defer func() {
		client.Hub.Unregister(client)
		conn.Close()
	}()

	conn.SetReadLimit(cfg.MaxMessageSz)
	conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
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

func writePump(client *ws.Client, conn *websocket.Conn, cfg Config) {
	ticker := time.NewTicker(cfg.PingInterval)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteJSON(msg); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
