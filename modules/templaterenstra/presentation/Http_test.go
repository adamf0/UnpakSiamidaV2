package presentation

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestModuleTemplateRenstraRegistersPutRoute(t *testing.T) {
	app := fiber.New()
	ModuleTemplateRenstra(app)

	foundPut := false
	for _, route := range app.GetRoutes() {
		if route.Method == fiber.MethodPut && route.Path == "/templaterenstra/:uuid" {
			foundPut = true
			break
		}
	}
	if !foundPut {
		t.Fatalf("expected PUT /templaterenstra/:uuid route to be registered")
	}
}
