package emailhub

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// EmailStatusEvent represents a real-time email status broadcast event
type EmailStatusEvent struct {
	Type           string     `json:"type"` // "EMAIL_STATUS_UPDATE" or "EMAIL_INSERTED"
	ID             uint       `json:"id,omitempty"`
	OutboxID       string     `json:"outbox_id"`
	UserUUID       string     `json:"user_uuid"`
	RecipientEmail string     `json:"recipient_email"`
	Username       string     `json:"username,omitempty"`
	Subject        string     `json:"subject,omitempty"`
	Status         string     `json:"status"` // PENDING, SENDING, SENT, FAILED
	Provider       string     `json:"provider,omitempty"`
	MessageID      string     `json:"message_id,omitempty"`
	ErrorMessage   *string    `json:"error_message,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// WSClient wraps websocket connection with a thread-safe send channel
type WSClient struct {
	conn *websocket.Conn
	send chan []byte
}

// EmailHub coordinates real-time event broadcasting to WS and SSE clients
type EmailHub struct {
	mu         sync.RWMutex
	wsClients  map[*WSClient]bool
	sseClients map[chan []byte]bool
}

var (
	defaultHub *EmailHub
	once       sync.Once
)

// GetHub returns the singleton EmailHub instance
func GetHub() *EmailHub {
	once.Do(func() {
		defaultHub = &EmailHub{
			wsClients:  make(map[*WSClient]bool),
			sseClients: make(map[chan []byte]bool),
		}
	})
	return defaultHub
}

// RegisterWS registers and starts writer pump for a websocket client
func (h *EmailHub) RegisterWS(conn *websocket.Conn) *WSClient {
	client := &WSClient{
		conn: conn,
		send: make(chan []byte, 32),
	}

	h.mu.Lock()
	h.wsClients[client] = true
	h.mu.Unlock()

	// Writer pump goroutine
	go func() {
		defer func() {
			_ = conn.Close()
		}()
		for msg := range client.send {
			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()

	return client
}

// UnregisterWS cleans up and removes a websocket client
func (h *EmailHub) UnregisterWS(client *WSClient) {
	if client == nil {
		return
	}
	h.mu.Lock()
	if _, ok := h.wsClients[client]; ok {
		delete(h.wsClients, client)
		close(client.send)
	}
	h.mu.Unlock()
}

// RegisterSSE registers an SSE channel
func (h *EmailHub) RegisterSSE() chan []byte {
	ch := make(chan []byte, 32)
	h.mu.Lock()
	h.sseClients[ch] = true
	h.mu.Unlock()
	return ch
}

// UnregisterSSE removes an SSE channel
func (h *EmailHub) UnregisterSSE(ch chan []byte) {
	h.mu.Lock()
	if _, ok := h.sseClients[ch]; ok {
		delete(h.sseClients, ch)
		close(ch)
	}
	h.mu.Unlock()
}

// Broadcast sends the event to all registered WS and SSE listeners
func (h *EmailHub) Broadcast(event EmailStatusEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("[EmailHub] Failed to marshal event: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	// Broadcast to WebSocket clients
	for client := range h.wsClients {
		select {
		case client.send <- payload:
		default:
			// Buffer full, skip to avoid blocking worker
		}
	}

	// Broadcast to SSE clients
	for sseCh := range h.sseClients {
		select {
		case sseCh <- payload:
		default:
			// Buffer full, skip
		}
	}
}

// HandleSSE creates an SSE stream writer for Fiber
func HandleSSE(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	hub := GetHub()
	ch := hub.RegisterSSE()

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer hub.UnregisterSSE(ch)

		// Initial greeting/ping
		_, _ = fmt.Fprintf(w, "event: connected\ndata: {\"connected\":true}\n\n")
		_ = w.Flush()

		ticker := time.NewTicker(25 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				_, err := fmt.Fprintf(w, "data: %s\n\n", string(msg))
				if err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}
			case <-ticker.C:
				_, err := fmt.Fprintf(w, ": ping\n\n")
				if err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}
			}
		}
	})

	return nil
}
