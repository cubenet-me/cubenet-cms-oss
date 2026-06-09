package ws

import (
	"encoding/json"
	"log/slog"
	"sync"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Client struct {
	ID   string
	Send chan Message
	Hub  *Hub
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = true
	slog.Debug("ws client connected", "client_id", c.ID)
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.Send)
		slog.Debug("ws client disconnected", "client_id", c.ID)
	}
}

func (h *Hub) Broadcast(msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		select {
		case c.Send <- msg:
		default:
		}
	}
}

func (h *Hub) SendTo(clientID string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.ID == clientID {
			select {
			case c.Send <- msg:
			default:
			}
			return
		}
	}
}
