package presentation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	commonpresentation "UnpakSiamida/common/presentation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TestUpdateProkerBatalSettingHandler_RejectsNonAdmin(t *testing.T) {
	prevDB := dbInstance
	dbInstance = &gorm.DB{}
	defer func() { dbInstance = prevDB }()

	app := fiber.New()
	app.Put("/settings/proker-batal", func(c *fiber.Ctx) error {
		c.Locals("user", &commonpresentation.Account{
			UUID:  "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
			Level: "user",
		})
		c.Locals("level", "user")
		return UpdateProkerBatalSettingHandler(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/settings/proker-batal", strings.NewReader(`{"proker_batal_aktif_tanggal":"2026-11-15"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for non-admin, got %d", resp.StatusCode)
	}
}
