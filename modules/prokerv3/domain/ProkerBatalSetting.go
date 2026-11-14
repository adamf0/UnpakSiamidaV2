package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const ProkerBatalAktifTanggalKey = "proker_batal_aktif_tanggal"

type AppSetting struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SettingKey   string    `gorm:"column:setting_key;type:varchar(100);uniqueIndex;not null" json:"setting_key"`
	SettingValue *string   `gorm:"column:setting_value;type:text" json:"setting_value"`
	UpdatedBy    *int      `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (AppSetting) TableName() string {
	return "app_settings"
}

// UnitBatalSchedule holds the activation date for a specific fakultas / unit.
type UnitBatalSchedule struct {
	IdFakultasUnit   int    `json:"id_fakultas_unit"`
	UuidFakultasUnit string `json:"uuid_fakultas_unit,omitempty"`
	NamaFakultasUnit string `json:"nama_fakultas_unit,omitempty"`
	Tanggal          string `json:"tanggal"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

// ProkerBatalConfig stores both the global fallback date and unit-specific overrides.
type ProkerBatalConfig struct {
	GlobalDate string              `json:"global_date"`
	Units      []UnitBatalSchedule `json:"units"`
}

// UnitBatalStatus is the evaluated status for a specific unit.
type UnitBatalStatus struct {
	IdFakultasUnit   int    `json:"id_fakultas_unit"`
	UuidFakultasUnit string `json:"uuid_fakultas_unit,omitempty"`
	NamaFakultasUnit string `json:"nama_fakultas_unit,omitempty"`
	Tanggal          string `json:"tanggal"`
	HasDate          bool   `json:"has_date"`
	IsActive         bool   `json:"is_active"`
	Status           string `json:"status"`
	StatusText       string `json:"status_label"`
	Message          string `json:"message"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

// ProkerBatalSettingStatus holds the evaluated activation state for Proker Batal.
type ProkerBatalSettingStatus struct {
	SettingKey     string            `json:"setting_key"`
	ActiveDate     string            `json:"proker_batal_aktif_tanggal"`
	GlobalDate     string            `json:"global_date"`
	HasDate        bool              `json:"has_date"`
	IsActive       bool              `json:"is_active"`
	Status         string            `json:"status"`
	StatusText     string            `json:"status_label"`
	Message        string            `json:"message"`
	UpdatedAt      string            `json:"updated_at,omitempty"`
	Scope          string            `json:"scope,omitempty"` // "unit" | "global" | "none"
	TargetUnitID   *int              `json:"target_unit_id,omitempty"`
	TargetUnitUUID string            `json:"target_unit_uuid,omitempty"`
	Units          []UnitBatalStatus `json:"units,omitempty"`
}

func formatTanggalSettingIndo(t time.Time) string {
	bulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	m := int(t.Month())
	if m >= 1 && m <= 12 {
		return fmt.Sprintf("%02d %s %d", t.Day(), bulan[m], t.Year())
	}
	return t.Format("02-01-2006")
}

// ParseProkerBatalConfig decodes the stored setting string.
// Supports both modern JSON config (with global_date & units) and legacy plain date string.
func ParseProkerBatalConfig(raw string) ProkerBatalConfig {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ProkerBatalConfig{
			GlobalDate: "",
			Units:      []UnitBatalSchedule{},
		}
	}
	if strings.HasPrefix(trimmed, "{") {
		var cfg ProkerBatalConfig
		if err := json.Unmarshal([]byte(trimmed), &cfg); err == nil {
			if cfg.Units == nil {
				cfg.Units = []UnitBatalSchedule{}
			}
			return cfg
		}
	}
	// Legacy fallback: plain string date
	return ProkerBatalConfig{
		GlobalDate: trimmed,
		Units:      []UnitBatalSchedule{},
	}
}

// SerializeProkerBatalConfig converts config struct to JSON string for storage in app_settings.
func SerializeProkerBatalConfig(cfg ProkerBatalConfig) (string, error) {
	if cfg.Units == nil {
		cfg.Units = []UnitBatalSchedule{}
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// EvalDate evaluates a single YYYY-MM-DD date string against now.
func EvalDate(dateStr string, now time.Time, unitName string) (hasDate bool, isActive bool, status string, statusText string, message string) {
	trimmed := strings.TrimSpace(dateStr)
	if trimmed == "" {
		msg := "Tanggal aktif fitur Proker Batal belum diatur oleh Admin."
		if unitName != "" {
			msg = fmt.Sprintf("Tanggal aktif fitur Proker Batal untuk %s belum diatur oleh Admin.", unitName)
		}
		return false, false, "belum_diatur", "Belum Diatur", msg
	}

	loc := now.Location()
	parsed, err := time.ParseInLocation("2006-01-02", trimmed, loc)
	if err != nil {
		msg := "Konfigurasi tanggal aktif Proker Batal tidak valid."
		if unitName != "" {
			msg = fmt.Sprintf("Konfigurasi tanggal aktif Proker Batal untuk %s tidak valid.", unitName)
		}
		return false, false, "tidak_valid", "Tanggal Tidak Valid", msg
	}

	startOfActive := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	formatted := formatTanggalSettingIndo(parsed)

	if !today.Before(startOfActive) {
		msg := fmt.Sprintf("Fitur Proker Batal aktif mulai %s", formatted)
		if unitName != "" {
			msg = fmt.Sprintf("Fitur Proker Batal untuk %s aktif mulai %s", unitName, formatted)
		}
		return true, true, "aktif", "Aktif", msg
	}

	msg := fmt.Sprintf("Fitur Proker Batal akan aktif pada tanggal %s", formatted)
	if unitName != "" {
		msg = fmt.Sprintf("Fitur Proker Batal untuk %s akan aktif pada tanggal %s", unitName, formatted)
	}
	return true, false, "belum_aktif", "Belum Aktif", msg
}

// EvalProkerBatalSettingWithUnit evaluates whether Proker Batal is active for targetUnitID on `now`.
// If targetUnitID has a unit-specific override in cfg.Units, it evaluates that override.
// Otherwise, it falls back to cfg.GlobalDate.
// If targetUnitID is nil or 0, it evaluates the GlobalDate and returns the full list of Units.
func EvalProkerBatalSettingWithUnit(rawSettingStr string, now time.Time, targetUnitID *int) ProkerBatalSettingStatus {
	cfg := ParseProkerBatalConfig(rawSettingStr)

	// If evaluating for a specific unit
	if targetUnitID != nil && *targetUnitID > 0 {
		var matchedUnit *UnitBatalSchedule
		for i := range cfg.Units {
			if cfg.Units[i].IdFakultasUnit == *targetUnitID {
				matchedUnit = &cfg.Units[i]
				break
			}
		}

		if matchedUnit != nil {
			hasDate, isActive, status, statusText, msg := EvalDate(matchedUnit.Tanggal, now, matchedUnit.NamaFakultasUnit)
			actDate := matchedUnit.Tanggal
			if !hasDate || status == "tidak_valid" {
				actDate = ""
			}
			return ProkerBatalSettingStatus{
				SettingKey:   ProkerBatalAktifTanggalKey,
				ActiveDate:   actDate,
				GlobalDate:   cfg.GlobalDate,
				HasDate:      hasDate,
				IsActive:     isActive,
				Status:       status,
				StatusText:   statusText,
				Message:      msg,
				Scope:        "unit",
				TargetUnitID: targetUnitID,
			}
		}

		// Unit has no override -> fallback to global
		hasDate, isActive, status, statusText, msg := EvalDate(cfg.GlobalDate, now, "")
		actDate := cfg.GlobalDate
		if !hasDate || status == "tidak_valid" {
			actDate = ""
		}
		scope := "global"
		if !hasDate {
			scope = "none"
		}
		return ProkerBatalSettingStatus{
			SettingKey:   ProkerBatalAktifTanggalKey,
			ActiveDate:   actDate,
			GlobalDate:   cfg.GlobalDate,
			HasDate:      hasDate,
			IsActive:     isActive,
			Status:       status,
			StatusText:   statusText,
			Message:      msg,
			Scope:        scope,
			TargetUnitID: targetUnitID,
		}
	}

	// No unit specified (Admin global view)
	hasDate, isActive, status, statusText, msg := EvalDate(cfg.GlobalDate, now, "")
	actDate := cfg.GlobalDate
	if !hasDate || status == "tidak_valid" {
		actDate = ""
	}
	unitStatuses := make([]UnitBatalStatus, 0, len(cfg.Units))
	for _, u := range cfg.Units {
		uHasDate, uActive, uStatus, uStatusText, uMsg := EvalDate(u.Tanggal, now, u.NamaFakultasUnit)
		uDate := u.Tanggal
		if !uHasDate || uStatus == "tidak_valid" {
			uDate = ""
		}
		unitStatuses = append(unitStatuses, UnitBatalStatus{
			IdFakultasUnit:   u.IdFakultasUnit,
			NamaFakultasUnit: u.NamaFakultasUnit,
			Tanggal:          uDate,
			HasDate:          uHasDate,
			IsActive:         uActive,
			Status:           uStatus,
			StatusText:       uStatusText,
			Message:          uMsg,
			UpdatedAt:        u.UpdatedAt,
		})
	}

	scope := "global"
	if !hasDate && len(unitStatuses) == 0 {
		scope = "none"
	}

	return ProkerBatalSettingStatus{
		SettingKey:   ProkerBatalAktifTanggalKey,
		ActiveDate:   actDate,
		GlobalDate:   cfg.GlobalDate,
		HasDate:      hasDate,
		IsActive:     isActive,
		Status:       status,
		StatusText:   statusText,
		Message:      msg,
		Scope:        scope,
		Units:        unitStatuses,
	}
}

// EvalProkerBatalSetting evaluates whether the Proker Batal feature is active on `now` using global config.
func EvalProkerBatalSetting(activeDateStr string, now time.Time) ProkerBatalSettingStatus {
	return EvalProkerBatalSettingWithUnit(activeDateStr, now, nil)
}
