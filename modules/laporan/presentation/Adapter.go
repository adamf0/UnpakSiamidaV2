package presentation

import (
	"encoding/json"
	"fmt"

	"UnpakSiamida/modules/laporan/domain"

	"github.com/gofiber/fiber/v2"
)

// =================
// Adapter Interface
// =================
type OutputAdapter[T any] interface {
	Send(c *fiber.Ctx, data domain.Paged[T]) error
}

// =================
// Paging JSON Adapter
// =================
type PagingAdapter[T any] struct{}

func (a *PagingAdapter[T]) Send(
	c *fiber.Ctx,
	data domain.Paged[T],
) error {
	return c.JSON(data)
}

// =================
// All Data JSON Adapter
// =================
type AllAdapter[T any] struct{}

func (a *AllAdapter[T]) Send(
	c *fiber.Ctx,
	data domain.Paged[T],
) error {
	return c.JSON(data.Data)
}

// =================
// NDJSON Adapter
// =================
type NDJSONAdapter[T any] struct{}

func (a *NDJSONAdapter[T]) Send(
	c *fiber.Ctx,
	data domain.Paged[T],
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
	data domain.Paged[T],
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
