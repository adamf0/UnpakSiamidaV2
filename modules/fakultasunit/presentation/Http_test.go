package presentation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestParseOptionalBool(t *testing.T) {
	cases := []struct {
		in    string
		want  bool
		valid bool
	}{
		{"true", true, true},
		{"1", true, true},
		{"false", false, true},
		{"0", false, true},
		{"invalid", false, false},
		{"", false, false},
	}

	for _, tc := range cases {
		got, ok := parseOptionalBool(tc.in)
		if ok != tc.valid {
			t.Fatalf("parseOptionalBool(%q) valid=%v want %v", tc.in, ok, tc.valid)
		}
		if ok && *got != tc.want {
			t.Fatalf("parseOptionalBool(%q) val=%v want %v", tc.in, *got, tc.want)
		}
	}
}

func TestModuleFakultasUnitHasNoDeleteRoute(t *testing.T) {
	app := fiber.New()
	ModuleFakultasUnit(app)

	req := httptest.NewRequest(http.MethodDelete, "/fakultasunit/428756bb-d2e4-449d-baa7-40bad1044438", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected DELETE /fakultasunit/:uuid to be disallowed, got %d", resp.StatusCode)
	}
}
