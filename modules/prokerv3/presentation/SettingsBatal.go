package presentation

import (
	"context"
	"errors"
	"strings"
	"time"

	commonpresentation "UnpakSiamida/common/presentation"
	domain "UnpakSiamida/modules/prokerv3/domain"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ensureAppSettingsTable(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS app_settings (
			id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
			setting_key varchar(100) NOT NULL,
			setting_value text DEFAULT NULL,
			updated_by int(11) DEFAULT NULL,
			created_at timestamp NULL DEFAULT current_timestamp(),
			updated_at timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
			PRIMARY KEY (id),
			UNIQUE KEY uk_app_settings_key (setting_key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`).Error
}

func getProkerBatalSettingStatus(ctx context.Context, db *gorm.DB, now time.Time) (domain.ProkerBatalSettingStatus, error) {
	if db == nil {
		return domain.ProkerBatalSettingStatus{}, errors.New("Database tidak tersedia")
	}

	var row domain.AppSetting
	err := db.WithContext(ctx).
		Where("setting_key = ?", domain.ProkerBatalAktifTanggalKey).
		Take(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.EvalProkerBatalSetting("", now), nil
		}
		return domain.ProkerBatalSettingStatus{}, err
	}

	val := ""
	if row.SettingValue != nil {
		val = *row.SettingValue
	}
	status := domain.EvalProkerBatalSetting(val, now)
	if !row.UpdatedAt.IsZero() {
		status.UpdatedAt = row.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	return status, nil
}

func GetProkerBatalSettingHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}
	status, err := getProkerBatalSettingStatus(c.Context(), dbInstance, time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(status)
}

func UpdateProkerBatalSettingHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	isAdminRole := false
	if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
		isAdminRole = strings.EqualFold(strings.TrimSpace(u.Level), "admin")
	}
	if !isAdminRole {
		level, _ := commonpresentation.GetUserInfo(c)
		isAdminRole = strings.EqualFold(strings.TrimSpace(level), "admin")
	}
	if !isAdminRole {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Hanya Admin yang dapat mengubah tanggal aktif Proker Batal",
		})
	}

	type updateSettingBody struct {
		ProkerBatalAktifTanggal *string `json:"proker_batal_aktif_tanggal" form:"proker_batal_aktif_tanggal"`
		Tanggal                 *string `json:"tanggal" form:"tanggal"`
		ActiveDate              *string `json:"active_date" form:"active_date"`
	}

	var body updateSettingBody
	_ = c.BodyParser(&body)

	rawDate := ""
	if body.ProkerBatalAktifTanggal != nil {
		rawDate = *body.ProkerBatalAktifTanggal
	} else if body.Tanggal != nil {
		rawDate = *body.Tanggal
	} else if body.ActiveDate != nil {
		rawDate = *body.ActiveDate
	} else if val := c.FormValue("proker_batal_aktif_tanggal"); val != "" {
		rawDate = val
	} else if val := c.FormValue("tanggal"); val != "" {
		rawDate = val
	} else if val := c.FormValue("active_date"); val != "" {
		rawDate = val
	}
	rawDate = strings.TrimSpace(rawDate)

	var valuePtr *string
	if rawDate != "" {
		parsed, err := time.ParseInLocation("2006-01-02", rawDate, time.Local)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Format tanggal aktif Proker Batal harus YYYY-MM-DD",
			})
		}
		normalized := parsed.Format("2006-01-02")
		valuePtr = &normalized
	}

	var updatedByPtr *int
	actorSID, _ := c.Locals("sid").(string)
	if strings.TrimSpace(actorSID) == "" {
		if u, ok := c.Locals("user").(*commonpresentation.Account); ok && u != nil {
			actorSID = u.UUID
		}
	}
	if strings.TrimSpace(actorSID) != "" {
		var uid int
		if err := dbInstance.WithContext(c.Context()).Table("users").Select("id").Where("uuid = ?", strings.TrimSpace(actorSID)).Take(&uid).Error; err == nil && uid > 0 {
			updatedByPtr = &uid
		}
	}

	if err := dbInstance.WithContext(c.Context()).Exec(`
		INSERT INTO app_settings (setting_key, setting_value, updated_by, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			setting_value = VALUES(setting_value),
			updated_by = VALUES(updated_by),
			updated_at = NOW()
	`, domain.ProkerBatalAktifTanggalKey, valuePtr, updatedByPtr).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	status, err := getProkerBatalSettingStatus(c.Context(), dbInstance, time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	msg := "Tanggal aktif Proker Batal berhasil disimpan"
	if valuePtr == nil {
		msg = "Tanggal aktif Proker Batal berhasil direset (fitur dinonaktifkan)"
	}

	return c.JSON(fiber.Map{
		"message":                   msg,
		"setting_key":               status.SettingKey,
		"proker_batal_aktif_tanggal": status.ActiveDate,
		"has_date":                  status.HasDate,
		"is_active":                 status.IsActive,
		"status":                    status.Status,
		"status_label":              status.StatusText,
		"updated_at":                status.UpdatedAt,
		"detail_message":            status.Message,
	})
}
