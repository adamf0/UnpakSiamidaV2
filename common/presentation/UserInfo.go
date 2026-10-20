package presentation

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetUserInfo extracts the active user role/level and fakultas_unit from Fiber context, headers, locals, or request params.
func GetUserInfo(c *fiber.Ctx) (level string, fakultasUnit string) {
	if hLevel := strings.ToLower(strings.TrimSpace(c.Get("X-Active-Role"))); hLevel != "" {
		level = hLevel
	} else if hLevel := strings.ToLower(strings.TrimSpace(c.Get("X-Level"))); hLevel != "" {
		level = hLevel
	}

	if level == "" {
		if lvl, ok := c.Locals("level").(string); ok && lvl != "" {
			level = strings.ToLower(strings.TrimSpace(lvl))
		}
	}
	if fu, ok := c.Locals("fakultas_unit").(string); ok && fu != "" {
		fakultasUnit = strings.TrimSpace(fu)
	}

	if level == "" {
		level = strings.ToLower(strings.TrimSpace(c.FormValue("level")))
		if level == "" {
			level = strings.ToLower(strings.TrimSpace(c.Query("level")))
		}
	}

	if fakultasUnit == "" {
		fakultasUnit = strings.TrimSpace(c.FormValue("fakultas_unit"))
		if fakultasUnit == "" {
			fakultasUnit = strings.TrimSpace(c.FormValue("target"))
		}
		if fakultasUnit == "" {
			fakultasUnit = strings.TrimSpace(c.Query("filter_fakultas_unit"))
		}
		if fakultasUnit == "" {
			fakultasUnit = strings.TrimSpace(c.Query("fakultas_unit"))
		}
		if fakultasUnit == "" {
			fakultasUnit = strings.TrimSpace(c.Query("target"))
		}
	}

	if u, ok := c.Locals("user").(*Account); ok && u != nil {
		if level == "" || level == "user" {
			for _, r := range u.ExtraRole {
				rRole := strings.ToLower(strings.TrimSpace(r.Role))
				if rRole == "upml" || rRole == "upmf" || rRole == "upms" {
					level = rRole
					break
				}
			}
		}

		if fakultasUnit == "" {
			for _, r := range u.ExtraRole {
				if level != "" && !strings.EqualFold(r.Role, level) && level != "user" {
					continue
				}
				if r.IdFakultasUnit > 0 {
					fakultasUnit = strconv.Itoa(int(r.IdFakultasUnit))
					break
				} else if r.NamaFakultasUnit != "" {
					fakultasUnit = r.NamaFakultasUnit
					break
				}
			}
		}
	}

	return level, fakultasUnit
}
