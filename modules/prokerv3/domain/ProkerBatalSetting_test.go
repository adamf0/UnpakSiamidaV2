package domain

import (
	"testing"
	"time"
)

func TestEvalProkerBatalSetting(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)
	now := time.Date(2026, 11, 10, 14, 30, 0, 0, loc)

	tests := []struct {
		name           string
		activeDate     string
		wantActive     bool
		wantHasDate    bool
		wantStatus     string
		wantActiveDate string
	}{
		{
			name:        "empty date is inactive",
			activeDate:  "",
			wantActive:  false,
			wantHasDate: false,
			wantStatus:  "belum_diatur",
		},
		{
			name:        "whitespace only is inactive",
			activeDate:  "   ",
			wantActive:  false,
			wantHasDate: false,
			wantStatus:  "belum_diatur",
		},
		{
			name:        "invalid date format is inactive",
			activeDate:  "10-11-2026",
			wantActive:  false,
			wantHasDate: false,
			wantStatus:  "tidak_valid",
		},
		{
			name:           "future date is not yet active",
			activeDate:     "2026-11-11",
			wantActive:     false,
			wantHasDate:    true,
			wantStatus:     "belum_aktif",
			wantActiveDate: "2026-11-11",
		},
		{
			name:           "same day is active",
			activeDate:     "2026-11-10",
			wantActive:     true,
			wantHasDate:    true,
			wantStatus:     "aktif",
			wantActiveDate: "2026-11-10",
		},
		{
			name:           "past date is active",
			activeDate:     "2026-11-01",
			wantActive:     true,
			wantHasDate:    true,
			wantStatus:     "aktif",
			wantActiveDate: "2026-11-01",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := EvalProkerBatalSetting(tc.activeDate, now)
			if got.IsActive != tc.wantActive {
				t.Fatalf("IsActive = %v, want %v (msg=%q)", got.IsActive, tc.wantActive, got.Message)
			}
			if got.HasDate != tc.wantHasDate {
				t.Fatalf("HasDate = %v, want %v", got.HasDate, tc.wantHasDate)
			}
			if got.Status != tc.wantStatus {
				t.Fatalf("Status = %q, want %q", got.Status, tc.wantStatus)
			}
			if got.ActiveDate != tc.wantActiveDate {
				t.Fatalf("ActiveDate = %q, want %q", got.ActiveDate, tc.wantActiveDate)
			}
		})
	}
}

func TestEvalProkerBatalSettingWithUnit(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)
	now := time.Date(2026, 11, 10, 14, 30, 0, 0, loc) // 10 Nov 2026

	jsonConfig := `{
		"global_date": "2026-11-10",
		"units": [
			{ "id_fakultas_unit": 6, "nama_fakultas_unit": "Fakultas MIPA", "tanggal": "2026-11-15" },
			{ "id_fakultas_unit": 8, "nama_fakultas_unit": "Sekolah Vokasi", "tanggal": "2026-11-05" }
		]
	}`

	t.Run("admin view returns global date and all unit statuses", func(t *testing.T) {
		got := EvalProkerBatalSettingWithUnit(jsonConfig, now, nil)
		if !got.IsActive {
			t.Fatalf("expected global to be active")
		}
		if got.ActiveDate != "2026-11-10" {
			t.Fatalf("expected active date 2026-11-10, got %s", got.ActiveDate)
		}
		if len(got.Units) != 2 {
			t.Fatalf("expected 2 units, got %d", len(got.Units))
		}
		if got.Units[0].IdFakultasUnit != 6 || got.Units[0].IsActive {
			t.Fatalf("expected unit 6 to be inactive (date 2026-11-15 > now 2026-11-10)")
		}
		if got.Units[1].IdFakultasUnit != 8 || !got.Units[1].IsActive {
			t.Fatalf("expected unit 8 to be active (date 2026-11-05 <= now 2026-11-10)")
		}
	})

	t.Run("unit with future override is inactive", func(t *testing.T) {
		unitID := 6
		got := EvalProkerBatalSettingWithUnit(jsonConfig, now, &unitID)
		if got.IsActive {
			t.Fatalf("expected unit 6 to be inactive")
		}
		if got.Scope != "unit" {
			t.Fatalf("expected scope 'unit', got %s", got.Scope)
		}
		if got.ActiveDate != "2026-11-15" {
			t.Fatalf("expected active date 2026-11-15, got %s", got.ActiveDate)
		}
	})

	t.Run("unit with past override is active", func(t *testing.T) {
		unitID := 8
		got := EvalProkerBatalSettingWithUnit(jsonConfig, now, &unitID)
		if !got.IsActive {
			t.Fatalf("expected unit 8 to be active")
		}
		if got.Scope != "unit" {
			t.Fatalf("expected scope 'unit', got %s", got.Scope)
		}
		if got.ActiveDate != "2026-11-05" {
			t.Fatalf("expected active date 2026-11-05, got %s", got.ActiveDate)
		}
	})

	t.Run("unit without override falls back to global", func(t *testing.T) {
		unitID := 99
		got := EvalProkerBatalSettingWithUnit(jsonConfig, now, &unitID)
		if !got.IsActive {
			t.Fatalf("expected unit 99 to fall back to active global")
		}
		if got.Scope != "global" {
			t.Fatalf("expected scope 'global', got %s", got.Scope)
		}
		if got.ActiveDate != "2026-11-10" {
			t.Fatalf("expected active date 2026-11-10, got %s", got.ActiveDate)
		}
	})
}
