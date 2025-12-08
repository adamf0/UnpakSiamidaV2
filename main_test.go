package main

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
	commonpresentation "UnpakSiamida/common/presentation"
)

func TestHeaderSecurityMiddleware(t *testing.T) {
	tests := []struct {
		name   string
		header string
		value  string
		want   int
	}{
		// ================= SUCCESS =================

		// {"[1] ok simple header", "User-Agent", "GoTest/1.0", 200},
		// {"[2] ok allow domain", "referer", "siamida.unpak.ac.id", 200},
		// {"[3] ok no domain inside", "X-Info", "just text here", 200},
		// {"[4] ok long but within limit", "X-Info", generate(200), 200},
		// {"[5] ok encoded domain allowlist", "Referer", "http%3A%2F%2Fsiamida.unpak.ac.id", 200},
		// {"[6] ok weird but safe", "X-Info", "__SAFE__", 200},
		// {"[7] ok random numeric", "X-Info", "123456789", 200},
		// {"[8] ok domain allowed suffix", "X-A", "test.internal.company", 200},
		// {"[9] ok unrelated words", "X-A", "company internal not domain", 200},

		// ================= FAIL =================

		{"[0] fail blacklisted header", "X-Forwarded-For", "1.2.3.4", 400},
		{"[10] fail dangerous proto", "X-A", "javascript:alert(1)", 400},
		{"[11] fail data proto", "X-A", "data:text/html", 400},
		{"[12] fail null byte", "X-A", "abc\x00def", 400},
		{"[13] fail CRLF", "X-A", "abc\r\nx:1", 400},
		{"[14] fail punycode", "X-A", "xn--evil-domain", 400},
		{"[15] fail zero-width", "X-A", "a\u200Bb", 400},
		{"[16] fail domain not allowed", "X-A", "http://evil.com", 400},
		{"[17] fail embedded domain not allowed", "X-A", "aaa http://malware.site/bb", 400},
		{"[18] fail decoded leads to dangerous proto", "X-A", "javascript%3Aalert(1)", 400},
		{"[19] fail decoded domain not allowed", "Referer", "http%3A%2F%2Fevil.com", 400},
		{"[20] fail multi-escape proto", "X-A", "javascript%253Aalert(1)", 400},
		{"[21] fail extremely long", "X-A", generate(9000), 400},

		// ================= EDGE CASES =================

		{"[22] fail host header spoof", "Host", "evil.com", 400},
		{"[23] fail domain inside text", "X-A", "aaa evil.com bbb", 400},
		{"[24] fail tricky subdomain", "X-A", "http://evil.com.example.org.evil.com", 400},
		{"[25] fail userinfo URL", "X-A", "http://admin:pw@evil.com", 400},
		{"[26] fail protocol missing but has evil domain", "X-A", "evil.com", 400},
	}

	cfg := commonpresentation.DefaultHeaderSecurityConfig()
	cfg.ResolveAndCheck = false

	app := fiber.New(fiber.Config{
		// DisableStartupMessage: true,
		ReadBufferSize: 16 * 1024,
		Prefork:      true, // gunakan semua CPU cores
		ServerHeader: "Fiber",
		// ReadTimeout: 10 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout: 10 * time.Second
	})
	app.Use(commonpresentation.LoggerMiddleware)
	app.Use(commonpresentation.HeaderSecurityMiddleware(cfg))
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(tc.header, tc.value)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != tc.want {
				t.Fatalf("got %d want %d", resp.StatusCode, tc.want)
			}
		})
	}
}

// helper
func generate(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = 'A'
	}
	return string(b)
}