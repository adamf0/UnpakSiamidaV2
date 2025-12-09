package presentation

import (
    "fmt"
    "encoding/json"
    "github.com/gofiber/fiber/v2"
    domainrenstra "UnpakSiamida/modules/renstra/domain"
)
// -----------------
// Adapter Interface
// -----------------
type OutputAdapter interface {
    Send(c *fiber.Ctx, renstras domainrenstra.PagedRenstras) error
}

// -----------------
// Paging JSON
// -----------------
type PagingAdapter struct{}

func (a *PagingAdapter) Send(c *fiber.Ctx, renstras domainrenstra.PagedRenstras) error {
    return c.JSON(renstras)
}

// -----------------
// All data JSON
// -----------------
type AllAdapter struct{}

func (a *AllAdapter) Send(c *fiber.Ctx, renstras domainrenstra.PagedRenstras) error {
    return c.JSON(renstras.Data)
}

// -----------------
// NDJSON
// -----------------
type NDJSONAdapter struct{}

func (a *NDJSONAdapter) Send(c *fiber.Ctx, renstras domainrenstra.PagedRenstras) error {
    c.Set("Content-Type", "application/x-ndjson")
    for _, u := range renstras.Data {
        line, _ := json.Marshal(u)
        fmt.Fprintln(c, string(line))
    }
    return nil
}

// -----------------
// SSE
// -----------------
type SSEAdapter struct{}

func (a *SSEAdapter) Send(c *fiber.Ctx, renstras domainrenstra.PagedRenstras) error {
    c.Set("Content-Type", "text/event-stream")
    c.Set("Cache-Control", "no-cache")
    c.Set("Connection", "keep-alive")

    totalCount := len(renstras.Data)
	c.Context().Write([]byte(fmt.Sprintf("total: %d\n\n", totalCount)))

    // Start event
    c.Context().Write([]byte("data: start\n\n"))
    
    for _, u := range renstras.Data {
        line, _ := json.Marshal(u)
        c.Context().Write([]byte("data: " + string(line) + "\n\n"))
    }

    // Done event
    c.Context().Write([]byte("data: done\n\n"))
    return nil
}