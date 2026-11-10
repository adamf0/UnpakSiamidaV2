package domain

import (
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

type ProkerBatalSettingStatus struct {
	SettingKey string `json:"setting_key"`
	ActiveDate string `json:"proker_batal_aktif_tanggal"`
	HasDate    bool   `json:"has_date"`
	IsActive   bool   `json:"is_active"`
	Status     string `json:"status"`
	StatusText string `json:"status_label"`
	Message    string `json:"message"`
	UpdatedAt  string `json:"updated_at,omitempty"`
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

// EvalProkerBatalSetting evaluates whether the Proker Batal feature is active for unit users on `now`.
// Active starts on the configured date at 00:00:00 local time (`today >= activeDate`).
func EvalProkerBatalSetting(activeDateStr string, now time.Time) ProkerBatalSettingStatus {
	trimmed := strings.TrimSpace(activeDateStr)
	if trimmed == "" {
		return ProkerBatalSettingStatus{
			SettingKey: ProkerBatalAktifTanggalKey,
			ActiveDate: "",
			HasDate:    false,
			IsActive:   false,
			Status:     "belum_diatur",
			StatusText: "Belum Diatur",
			Message:    "Tanggal aktif fitur Proker Batal belum diatur oleh Admin.",
		}
	}

	loc := now.Location()
	parsed, err := time.ParseInLocation("2006-01-02", trimmed, loc)
	if err != nil {
		return ProkerBatalSettingStatus{
			SettingKey: ProkerBatalAktifTanggalKey,
			ActiveDate: "",
			HasDate:    false,
			IsActive:   false,
			Status:     "tidak_valid",
			StatusText: "Tanggal Tidak Valid",
			Message:    "Konfigurasi tanggal aktif Proker Batal tidak valid.",
		}
	}

	startOfActive := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	formatted := formatTanggalSettingIndo(parsed)

	if !today.Before(startOfActive) {
		return ProkerBatalSettingStatus{
			SettingKey: ProkerBatalAktifTanggalKey,
			ActiveDate: parsed.Format("2006-01-02"),
			HasDate:    true,
			IsActive:   true,
			Status:     "aktif",
			StatusText: "Aktif",
			Message:    fmt.Sprintf("Fitur Proker Batal aktif mulai %s", formatted),
		}
	}

	return ProkerBatalSettingStatus{
		SettingKey: ProkerBatalAktifTanggalKey,
		ActiveDate: parsed.Format("2006-01-02"),
		HasDate:    true,
		IsActive:   false,
		Status:     "belum_aktif",
		StatusText: "Belum Aktif",
		Message:    fmt.Sprintf("Fitur Proker Batal akan aktif pada tanggal %s", formatted),
	}
}
