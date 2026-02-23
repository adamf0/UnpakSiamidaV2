package presentation

import (
	"encoding/json"
	"fmt"
	"strings"

	commondomain "UnpakSiamida/common/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// =================
// Adapter Interface
// =================
type OutputAdapter[T any] interface {
	Send(c *fiber.Ctx, data commondomain.Paged[T]) error
}

// =================
// Paging JSON Adapter
// =================
type PagingAdapter[T any] struct{}

func (a *PagingAdapter[T]) Send(
	c *fiber.Ctx,
	data commondomain.Paged[T],
) error {
	return c.JSON(data)
}

// =================
// All Data JSON Adapter
// =================
type AllAdapter[T any] struct{}

func (a *AllAdapter[T]) Send(
	c *fiber.Ctx,
	data commondomain.Paged[T],
) error {
	return c.JSON(data.Data)
}

// =================
// NDJSON Adapter
// =================
type NDJSONAdapter[T any] struct{}

func (a *NDJSONAdapter[T]) Send(
	c *fiber.Ctx,
	data commondomain.Paged[T],
) error {
	c.Set("Content-Type", "application/x-ndjson")

	for _, u := range data.Data {
		b, _ := json.Marshal(u)
		fmt.Fprintln(c, string(b))
	}

	return nil
}

// =================
// SSE Adapter
// =================
type SSEAdapter[T any] struct{}

func (a *SSEAdapter[T]) Send(
	c *fiber.Ctx,
	data commondomain.Paged[T],
) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	totalCount := len(data.Data)
	c.Context().Write([]byte(fmt.Sprintf("total: %d\n\n", totalCount)))

	// start event
	c.Context().Write([]byte("data: start\n\n"))

	for _, u := range data.Data {
		b, _ := json.Marshal(u)
		c.Context().Write([]byte("data: " + string(b) + "\n\n"))
	}

	// done event
	c.Context().Write([]byte("data: done\n\n"))
	return nil
}

// =================
// WebSocket Adapter
// =================
type WebSocketAdapter[T any] struct{}

func (a *WebSocketAdapter[T]) Send(
	c *fiber.Ctx,
	datas commondomain.Paged[T],
) error {

	// Harus upgrade
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	// Ambil subprotocol dan split by comma (RFC)
	protocolHeader := c.Get("Sec-WebSocket-Protocol")
	var subprotocols []string

	if protocolHeader != "" {
		for _, p := range strings.Split(protocolHeader, ",") {
			subprotocols = append(subprotocols, strings.TrimSpace(p))
		}
	}

	return websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		// Reader goroutine (agar koneksi tidak close)
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()

		// Meta
		conn.WriteJSON(map[string]interface{}{
			"type":  "meta",
			"total": len(datas.Data),
		})

		// Start
		conn.WriteJSON(map[string]interface{}{
			"type": "start",
		})

		// Stream data
		for _, item := range datas.Data {
			conn.WriteJSON(map[string]interface{}{
				"type": "data",
				"data": item,
			})
		}

		// Done
		conn.WriteJSON(map[string]interface{}{
			"type": "done",
		})

	}, websocket.Config{
		Subprotocols: subprotocols, // wajib echo balik
	})(c)
}
