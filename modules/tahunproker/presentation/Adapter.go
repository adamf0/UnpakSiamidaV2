package presentation

import (
	domainTahunProker "UnpakSiamida/modules/tahunproker/domain"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// -----------------
// Adapter Interface
// -----------------
type OutputAdapter interface {
	Send(c *fiber.Ctx, TahunProkers domainTahunProker.PagedTahunProkers) error
}

// -----------------
// Paging JSON
// -----------------
type PagingAdapter struct{}

func (a *PagingAdapter) Send(c *fiber.Ctx, TahunProkers domainTahunProker.PagedTahunProkers) error {
	return c.JSON(TahunProkers)
}

// -----------------
// All data JSON
// -----------------
type AllAdapter struct{}

func (a *AllAdapter) Send(c *fiber.Ctx, TahunProkers domainTahunProker.PagedTahunProkers) error {
	return c.JSON(TahunProkers.Data)
}

// -----------------
// NDJSON
// -----------------
type NDJSONAdapter struct{}

func (a *NDJSONAdapter) Send(c *fiber.Ctx, TahunProkers domainTahunProker.PagedTahunProkers) error {
	c.Set("Content-Type", "application/x-ndjson")
	for _, u := range TahunProkers.Data {
		line, _ := json.Marshal(u)
		fmt.Fprintln(c, string(line))
	}
	return nil
}

// -----------------
// SSE
// -----------------
type SSEAdapter struct{}

func (a *SSEAdapter) Send(c *fiber.Ctx, TahunProkers domainTahunProker.PagedTahunProkers) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	totalCount := len(TahunProkers.Data)
	c.Context().Write([]byte(fmt.Sprintf("total: %d\n\n", totalCount)))

	// Start event
	c.Context().Write([]byte("data: start\n\n"))

	for _, u := range TahunProkers.Data {
		line, _ := json.Marshal(u)
		c.Context().Write([]byte("data: " + string(line) + "\n\n"))
	}

	// Done event
	c.Context().Write([]byte("data: done\n\n"))
	return nil
}
