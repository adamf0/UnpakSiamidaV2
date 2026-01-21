package presentation

import (
	domainBeritaAcara "UnpakSiamida/modules/beritaacara/domain"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// -----------------
// Adapter Interface
// -----------------
type OutputAdapter interface {
	Send(c *fiber.Ctx, BeritaAcaras domainBeritaAcara.PagedBeritaAcaras) error
}

// -----------------
// Paging JSON
// -----------------
type PagingAdapter struct{}

func (a *PagingAdapter) Send(c *fiber.Ctx, BeritaAcaras domainBeritaAcara.PagedBeritaAcaras) error {
	return c.JSON(BeritaAcaras)
}

// -----------------
// All data JSON
// -----------------
type AllAdapter struct{}

func (a *AllAdapter) Send(c *fiber.Ctx, BeritaAcaras domainBeritaAcara.PagedBeritaAcaras) error {
	return c.JSON(BeritaAcaras.Data)
}

// -----------------
// NDJSON
// -----------------
type NDJSONAdapter struct{}

func (a *NDJSONAdapter) Send(c *fiber.Ctx, BeritaAcaras domainBeritaAcara.PagedBeritaAcaras) error {
	c.Set("Content-Type", "application/x-ndjson")
	for _, u := range BeritaAcaras.Data {
		line, _ := json.Marshal(u)
		fmt.Fprintln(c, string(line))
	}
	return nil
}

// -----------------
// SSE
// -----------------
type SSEAdapter struct{}

func (a *SSEAdapter) Send(c *fiber.Ctx, BeritaAcaras domainBeritaAcara.PagedBeritaAcaras) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	totalCount := len(BeritaAcaras.Data)
	c.Context().Write([]byte(fmt.Sprintf("total: %d\n\n", totalCount)))

	// Start event
	c.Context().Write([]byte("data: start\n\n"))

	for _, u := range BeritaAcaras.Data {
		line, _ := json.Marshal(u)
		c.Context().Write([]byte("data: " + string(line) + "\n\n"))
	}

	// Done event
	c.Context().Write([]byte("data: done\n\n"))
	return nil
}
