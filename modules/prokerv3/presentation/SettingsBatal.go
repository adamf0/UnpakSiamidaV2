package presentation

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

func getProkerBatalSettingStatus(ctx context.Context, db *gorm.DB, now time.Time, targetUnitID *int) (domain.ProkerBatalSettingStatus, error) {
	if db == nil {
		return domain.ProkerBatalSettingStatus{}, errors.New("Database tidak tersedia")
	}

	var row domain.AppSetting
	err := db.WithContext(ctx).
		Where("setting_key = ?", domain.ProkerBatalAktifTanggalKey).
		Take(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.EvalProkerBatalSettingWithUnit("", now, targetUnitID), nil
		}
		return domain.ProkerBatalSettingStatus{}, err
	}

	val := ""
	if row.SettingValue != nil {
		val = *row.SettingValue
	}
	status := domain.EvalProkerBatalSettingWithUnit(val, now, targetUnitID)
	if !row.UpdatedAt.IsZero() {
		status.UpdatedAt = row.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	return status, nil
}

func GetProkerBatalSettingHandler(c *fiber.Ctx) error {
	if dbInstance == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak tersedia"})
	}

	var targetUnitID *int
	unitParam := strings.TrimSpace(c.Query("unit_id"))
	if unitParam == "" {
		unitParam = strings.TrimSpace(c.Query("id_fakultas_unit"))
	}
	if unitParam != "" && unitParam != "all" {
		if id, err := strconv.Atoi(unitParam); err == nil && id > 0 {
			targetUnitID = &id
		}
	}

	// If no query param given, evaluate based on logged-in user unit for non-admins
	if targetUnitID == nil {
		level, userFU := commonpresentation.GetUserInfo(c)
		if !strings.EqualFold(level, "admin") && userFU != "" {
			if id, err := strconv.Atoi(userFU); err == nil && id > 0 {
				targetUnitID = &id
			} else {
				var uid int
				if err := dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").Select("id").
					Where("uuid = ? OR nama_fak_prod_unit = ?", userFU, userFU).
					Take(&uid).Error; err == nil && uid > 0 {
					targetUnitID = &uid
				}
			}
		}
	}

	status, err := getProkerBatalSettingStatus(c.Context(), dbInstance, time.Now(), targetUnitID)
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
		ProkerBatalAktifTanggal *string                    `json:"proker_batal_aktif_tanggal" form:"proker_batal_aktif_tanggal"`
		GlobalDate              *string                    `json:"global_date" form:"global_date"`
		Tanggal                 *string                    `json:"tanggal" form:"tanggal"`
		ActiveDate              *string                    `json:"active_date" form:"active_date"`
		Action                  string                     `json:"action" form:"action"`
		TargetUnits             []int                      `json:"target_units"`
		UnitID                  *int                       `json:"unit_id" form:"unit_id"`
		DeleteUnitID            *int                       `json:"delete_unit_id" form:"delete_unit_id"`
		Units                   *[]domain.UnitBatalSchedule `json:"units"`
	}

	var body updateSettingBody
	_ = c.BodyParser(&body)

	// Fetch current configuration
	var currentVal string
	var currentRow domain.AppSetting
	if err := dbInstance.WithContext(c.Context()).
		Where("setting_key = ?", domain.ProkerBatalAktifTanggalKey).
		Take(&currentRow).Error; err == nil && currentRow.SettingValue != nil {
		currentVal = *currentRow.SettingValue
	}
	cfg := domain.ParseProkerBatalConfig(currentVal)
	msg := ""

	action := strings.ToLower(strings.TrimSpace(body.Action))
	deleteID := body.DeleteUnitID
	if deleteID == nil && (action == "delete_unit" || action == "remove_unit") {
		deleteID = body.UnitID
	}

	if deleteID != nil && *deleteID > 0 {
		// Delete specific unit override
		filtered := make([]domain.UnitBatalSchedule, 0, len(cfg.Units))
		for _, u := range cfg.Units {
			if u.IdFakultasUnit != *deleteID {
				filtered = append(filtered, u)
			}
		}
		cfg.Units = filtered
		msg = "Jadwal khusus fakultas/unit berhasil dihapus (kembali mengikuti jadwal global)"
	} else if action == "reset_units" {
		cfg.Units = []domain.UnitBatalSchedule{}
		msg = "Semua jadwal khusus fakultas/unit berhasil direset"
	} else if action == "reset_all" {
		cfg.GlobalDate = ""
		cfg.Units = []domain.UnitBatalSchedule{}
		msg = "Pengaturan tanggal aktif Proker Batal berhasil direset total"
	} else if len(body.TargetUnits) > 0 {
		// Set / update multiple target units
		rawDate := ""
		if body.Tanggal != nil {
			rawDate = *body.Tanggal
		} else if body.ActiveDate != nil {
			rawDate = *body.ActiveDate
		} else if body.ProkerBatalAktifTanggal != nil {
			rawDate = *body.ProkerBatalAktifTanggal
		} else if val := c.FormValue("tanggal"); val != "" {
			rawDate = val
		}
		rawDate = strings.TrimSpace(rawDate)
		if rawDate == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tanggal aktif harus diisi untuk target fakultas/unit yang dipilih",
			})
		}
		parsed, err := time.ParseInLocation("2006-01-02", rawDate, time.Local)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Format tanggal aktif Proker Batal harus YYYY-MM-DD",
			})
		}
		normalizedDate := parsed.Format("2006-01-02")
		nowStr := time.Now().Format("2006-01-02 15:04:05")

		// Query unit names from v_fakultas_unit
		namesMap := make(map[int]string)
		var fuRows []struct {
			ID   int    `gorm:"column:id"`
			Nama string `gorm:"column:nama_fak_prod_unit"`
		}
		_ = dbInstance.WithContext(c.Context()).Table("v_fakultas_unit").
			Where("id IN (?)", body.TargetUnits).
			Select("id, nama_fak_prod_unit").
			Find(&fuRows).Error
		for _, r := range fuRows {
			namesMap[r.ID] = r.Nama
		}

		for _, tid := range body.TargetUnits {
			if tid <= 0 {
				continue
			}
			foundIdx := -1
			for i := range cfg.Units {
				if cfg.Units[i].IdFakultasUnit == tid {
					foundIdx = i
					break
				}
			}
			unitName := namesMap[tid]
			if foundIdx >= 0 {
				cfg.Units[foundIdx].Tanggal = normalizedDate
				cfg.Units[foundIdx].UpdatedAt = nowStr
				if unitName != "" {
					cfg.Units[foundIdx].NamaFakultasUnit = unitName
				}
			} else {
				if unitName == "" {
					unitName = fmt.Sprintf("Unit #%d", tid)
				}
				cfg.Units = append(cfg.Units, domain.UnitBatalSchedule{
					IdFakultasUnit:   tid,
					NamaFakultasUnit: unitName,
					Tanggal:          normalizedDate,
					UpdatedAt:        nowStr,
				})
			}
		}
		msg = fmt.Sprintf("Jadwal Proker Batal untuk %d fakultas/unit berhasil disimpan", len(body.TargetUnits))
	} else if body.Units != nil {
		// Direct full units replace
		validUnits := make([]domain.UnitBatalSchedule, 0, len(*body.Units))
		for _, u := range *body.Units {
			if u.IdFakultasUnit <= 0 {
				continue
			}
			parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(u.Tanggal), time.Local)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Format tanggal tidak valid untuk unit #%d (harus YYYY-MM-DD)", u.IdFakultasUnit),
				})
			}
			u.Tanggal = parsed.Format("2006-01-02")
			validUnits = append(validUnits, u)
		}
		cfg.Units = validUnits
		msg = "Jadwal khusus fakultas/unit berhasil diperbarui"
	} else {
		// Update global date
		rawDate := ""
		if body.GlobalDate != nil {
			rawDate = *body.GlobalDate
		} else if body.ProkerBatalAktifTanggal != nil {
			rawDate = *body.ProkerBatalAktifTanggal
		} else if body.Tanggal != nil {
			rawDate = *body.Tanggal
		} else if body.ActiveDate != nil {
			rawDate = *body.ActiveDate
		} else if val := c.FormValue("global_date"); val != "" {
			rawDate = val
		} else if val := c.FormValue("proker_batal_aktif_tanggal"); val != "" {
			rawDate = val
		} else if val := c.FormValue("tanggal"); val != "" {
			rawDate = val
		}
		rawDate = strings.TrimSpace(rawDate)

		if rawDate == "" || action == "reset_global" {
			cfg.GlobalDate = ""
			msg = "Tanggal aktif Proker Batal (Global) berhasil direset"
		} else {
			parsed, err := time.ParseInLocation("2006-01-02", rawDate, time.Local)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Format tanggal aktif Proker Batal harus YYYY-MM-DD",
				})
			}
			cfg.GlobalDate = parsed.Format("2006-01-02")
			msg = "Tanggal aktif Proker Batal (Global) berhasil disimpan"
		}
	}

	serialized, err := domain.SerializeProkerBatalConfig(cfg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
		if err := dbInstance.WithContext(c.Context()).Table("usersv3").Select("id").Where("uuid = ?", strings.TrimSpace(actorSID)).Take(&uid).Error; err == nil && uid > 0 {
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
	`, domain.ProkerBatalAktifTanggalKey, &serialized, updatedByPtr).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	status, err := getProkerBatalSettingStatus(c.Context(), dbInstance, time.Now(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":                    msg,
		"setting_key":                status.SettingKey,
		"proker_batal_aktif_tanggal": status.ActiveDate,
		"global_date":                status.GlobalDate,
		"has_date":                   status.HasDate,
		"is_active":                  status.IsActive,
		"status":                     status.Status,
		"status_label":               status.StatusText,
		"updated_at":                 status.UpdatedAt,
		"detail_message":             status.Message,
		"units":                      status.Units,
	})
}
