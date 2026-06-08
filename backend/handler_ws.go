package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PingInterval time.Duration
	MaxMessageSz int64
}

func (a *App) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws upgrade failed", "error", err)
		return
	}

	client := &Client{
		ID:   r.RemoteAddr,
		Send: make(chan Message, 64),
		Hub:  a.hub,
	}

	a.hub.Register(client)
	go writePump(client, conn, a.wsCfg)
	go readPump(client, conn, a.wsCfg)
}

func readPump(client *Client, conn *websocket.Conn, cfg WSConfig) {
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

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		slog.Debug("ws message", "type", msg.Type, "client_id", client.ID)
	}
}

func writePump(client *Client, conn *websocket.Conn, cfg WSConfig) {
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
