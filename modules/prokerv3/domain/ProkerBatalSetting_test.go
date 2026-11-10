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
