package presentation

import (
	"context"
	"html"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/unicode/norm"
	commoninfra "UnpakSiamida/common/infrastructure"
)

// =======================
// CONFIG
// =======================

type HeaderSecurityConfig struct {
	BlacklistedHeaderNames map[string]bool
	AllowDomains           []string
	MaxHeaderLen           int
	ResolveAndCheck        bool
	LookupTimeout          time.Duration
	BlockedCIDRs           []string
}

func DefaultBlacklistedHeaderNames() map[string]bool {
	names := []string{
		"x-forwarded-for", "x-forwarded-host", "forwarded", "forwarded-host",
		"x-forwarded-proto", "x-forwarded-port", "x-forwarded-scheme",
		"x-real-ip", "client-ip", "true-client-ip", "cf-connecting-ip",
		"x-remote-ip", "x-originating-ip",
		"x-original-host", "via", "x-via",
		"host", "x-host", "x-rewrite-url", "x-original-url",
		"x-request-url", "x-request-uri", "redirect", "location",
		"authorization", "proxy-authorization", "x-api-key",
		"metadata", "x-aws-ec2-metadata","referer",
	}
	out := map[string]bool{}
	for _, v := range names {
		out[strings.ToLower(v)] = true
	}
	return out
}

func DefaultHeaderSecurityConfig() *HeaderSecurityConfig {
	return &HeaderSecurityConfig{
		BlacklistedHeaderNames: DefaultBlacklistedHeaderNames(),
		AllowDomains:           []string{"siamida.unpak.ac.id","localhost:3000","thunderclient.com"},
		MaxHeaderLen:           8192,
		ResolveAndCheck:        false,
		LookupTimeout:          1 * time.Second,
		BlockedCIDRs:           []string{},
	}
}

// =======================
// REGEX
// =======================

var (
	crlfRe      = regexp.MustCompile(`[\r\n]`)
	nullRe      = regexp.MustCompile(`\x00`)
	protoRe     = regexp.MustCompile(`(?i)^(javascript|data|vbscript|file|view-source):`)
	punyRe      = regexp.MustCompile(`(?i)xn--[a-z0-9-]+`)
	zeroWidthRe = regexp.MustCompile(`[\x{200B}\x{200C}\x{200D}\x{2060}\x{FEFF}]`)
	hostExtract = regexp.MustCompile(`(?i)(?:https?://)?([a-z0-9\.\-]+\.[a-z]{2,})(:\d+)?`)
)

// =======================
// MIDDLEWARE
// =======================

func HeaderSecurityMiddleware(cfg *HeaderSecurityConfig) fiber.Handler {
	if cfg == nil {
		cfg = DefaultHeaderSecurityConfig()
	}

	var blocked []*net.IPNet
	for _, c := range cfg.BlockedCIDRs {
		_, ipnet, err := net.ParseCIDR(c)
		if err == nil {
			blocked = append(blocked, ipnet)
		}
	}

	return func(c *fiber.Ctx) error {
		for name, vals := range c.GetReqHeaders() {
			for _, val := range vals {
				// 1) Max header length
				if len(val) > cfg.MaxHeaderLen {
					return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+1]", "header too long: "+name))
				}

				// 2) Malicious control chars
				if crlfRe.MatchString(val) || nullRe.MatchString(val) {
					return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+2]", "header ctrl char: "+name))
				}

				// 3) Decode value safely
				decoded := multiUnescape(html.UnescapeString(val), 3)
				decoded = zeroWidthRe.ReplaceAllString(decoded, "")
				decoded = norm.NFKC.String(decoded)

				// 4) Dangerous protocols
				if protoRe.MatchString(decoded) {
					return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+3]", "protocol attack: "+name))
				}

				// 5) Punycode / IDN injection
				if punyRe.MatchString(decoded) {
					return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+4]", "punycode forbidden: "+name))
				}

				// 6) Zero-width check
				if zeroWidthRe.MatchString(val) {
					return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+5]", "zero width attack: "+name))
				}

				// 7) Domain allowlist cek untuk URL valid
				u, err := url.Parse(decoded)
				if err == nil && u.Host != "" {
					host := u.Hostname()
					if !domainAllowed(host, cfg.AllowDomains) {
						return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+6]", "domain not allowed: "+host))
					}

					// Optional DNS resolve
					if cfg.ResolveAndCheck {
						ctx, cancel := context.WithTimeout(context.Background(), cfg.LookupTimeout)
						defer cancel()
						ips, _ := net.DefaultResolver.LookupIP(ctx, "ip", host)
						for _, ip := range ips {
							if ipInNets(ip, blocked) {
								return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+7]", "domain resolves to forbidden IP: "+host))
							}
						}
					}
				}

				// 8) Cek host header spoof (Host harus sama atau allow domain)
				if strings.ToLower(name) == "host" {
					host := decoded
					if !domainAllowed(host, cfg.AllowDomains) {
						return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+8]", "host header spoof: "+host))
					}
				}

				// 9) Cek embedded domain dalam text
				hosts := extractHosts(decoded)
				for _, h := range hosts {
					if !domainAllowed(h, cfg.AllowDomains) {
						return c.Status(400).JSON(commoninfra.NewResponseError("common.check[A+9]", "embedded domain not allowed: "+h))
					}
				}
			}
		}

		return c.Next()
	}
}

// =======================
// HELPERS
// =======================

// func domainAllowed(host string, allow []string) bool {
// 	host = strings.ToLower(host)
// 	for _, a := range allow {
// 		if strings.HasSuffix(host, strings.ToLower(a)) {
// 			return true
// 		}
// 	}
// 	return false
// }
func domainAllowed(host string, allow []string) bool {
    host = strings.ToLower(host)
    for _, a := range allow {
        u, err := url.Parse(a)
        var domain string
        if err == nil && u.Host != "" {
            domain = u.Hostname()
        } else {
            domain = strings.ToLower(a)
        }
        if strings.HasSuffix(host, domain) {
            return true
        }
    }
    return false
}


func extractHosts(s string) []string {
	out := []string{}
	words := strings.Fields(s)
	for _, w := range words {
		u, err := url.Parse(w)
		if err == nil && u.Host != "" {
			out = append(out, u.Hostname())
			continue
		}
		m := hostExtract.FindStringSubmatch(w)
		if len(m) > 1 {
			out = append(out, m[1])
		}
	}
	return out
}

func multiUnescape(s string, n int) string {
	cur := s
	for i := 0; i < n; i++ {
		u, err := url.QueryUnescape(cur)
		if err != nil || u == cur {
			break
		}
		cur = u
	}
	return cur
}

func ipInNets(ip net.IP, nets []*net.IPNet) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
