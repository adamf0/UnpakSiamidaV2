package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/renstra/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Helper valid data
// -----------------------------------------------------------------------------
const (
	validDate1 = "2025-01-01"
	validDate2 = "2025-01-10"
	validDate3 = "2025-02-01"
	validDate4 = "2025-02-10"
	validDate5 = "2025-03-01"
	validDate6 = "2025-03-10"
)

// -----------------------------------------------------------------------------
// Test NewRenstra - Full Failure Scenarios
// -----------------------------------------------------------------------------
func NewRenstraFailures(t *testing.T) {
	tests := []struct {
		name        string
		fakultas    uint
		unique      bool
		auditee     uint
		aud1        uint
		aud2        uint
		date1       string
		date2       string
		date3       string
		date4       string
		date5       string
		date6       string
		expectedErr common.Error
	}{
		{
			name:        "InvalidFakultas",
			fakultas:    0,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			expectedErr: domain.InvalidFakultasUnit(),
		},
		{
			name:        "DataExistingFailsUniqueCheck",
			fakultas:    10,
			unique:      false,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			expectedErr: domain.DataExisting(),
		},
		{
			name:        "AuditeeMissing",
			fakultas:    10,
			unique:      true,
			auditee:     0,
			aud1:        2,
			aud2:        3,
			expectedErr: domain.MissingAuditee(),
		},
		{
			name:        "Auditor1Missing",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        0,
			aud2:        3,
			expectedErr: domain.MissingAuditor1(),
		},
		{
			name:        "Auditor2Missing",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        0,
			expectedErr: domain.MissingAuditor2(),
		},
		{
			name:        "DuplicateAssignment",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        1,
			aud2:        2,
			expectedErr: domain.DuplicateAssigment(),
		},
		{
			name:        "InvalidDateUpload",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       "invalid",
			expectedErr: domain.InvalidDate("upload"),
		},
		{
			name:        "InvalidDateUpload",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       validDate1,
			date2:       "invalid",
			expectedErr: domain.InvalidDate("upload"),
		},
		{
			name:        "InvalidDateDokumen",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       validDate1,
			date2:       validDate2,
			date3:       "invalid",
			expectedErr: domain.InvalidDate("assessment dokumen"),
		},
		{
			name:        "InvalidDateDokumen",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       validDate1,
			date2:       validDate2,
			date3:       validDate3,
			date4:       "invalid",
			expectedErr: domain.InvalidDate("assessment dokumen"),
		},
		{
			name:        "InvalidDateLapangan",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       validDate1,
			date2:       validDate2,
			date3:       validDate3,
			date4:       validDate4,
			date5:       "invalid",
			expectedErr: domain.InvalidDate("assessment lapangan"),
		},
		{
			name:        "InvalidDateLapangan",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       validDate1,
			date2:       validDate2,
			date3:       validDate3,
			date4:       validDate4,
			date5:       validDate5,
			date6:       "invalid",
			expectedErr: domain.InvalidDate("assessment lapangan"),
		},
		{
			name:        "OverlapUploadDokumen",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       validDate1,
			date2:       validDate6, // overlaps everything
			date3:       validDate3,
			date4:       validDate4,
			date5:       validDate5,
			date6:       validDate6,
			expectedErr: domain.PeriodOverlapUploadDokumen(),
		},
		{
			name:        "OverlapUploadLapangan",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       validDate1, // upload start
			date2:       validDate4, // upload end
			date3:       validDate5, // dokumen start
			date4:       validDate6, // dokumen end
			date5:       validDate3, // lapangan start
			date6:       validDate5, // lapangan end
			expectedErr: domain.PeriodOverlapUploadLapangan(),
		},

		{
			name:        "OverlapDokumenLapangan",
			fakultas:    10,
			unique:      true,
			auditee:     1,
			aud1:        2,
			aud2:        3,
			date1:       validDate1,
			date2:       validDate2,
			date3:       validDate3,
			date4:       validDate6, // overlap with lapangan
			date5:       validDate5,
			date6:       validDate6,
			expectedErr: domain.PeriodOverlapDokumenLapangan(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := domain.NewRenstra(
				"2025",
				tt.fakultas,
				tt.date1, tt.date2,
				tt.date3, tt.date4,
				tt.date5, tt.date6,
				tt.auditee, tt.aud1, tt.aud2,
				tt.unique,
			)

			require.False(t, r.IsSuccess)
			assert.Equal(t, tt.expectedErr.Code, r.Error.Code)
		})
	}
}

// -----------------------------------------------------------------------------
// Test NewRenstra Success
// -----------------------------------------------------------------------------
func NewRenstraSuccess(t *testing.T) {
	res := domain.NewRenstra(
		"2025",
		10,
		validDate1, validDate2,
		validDate3, validDate4,
		validDate5, validDate6,
		1, 2, 3,
		true, // unique
	)

	require.True(t, res.IsSuccess)

	r := res.Value

	assert.Equal(t, "2025", r.Tahun)
	assert.Equal(t, uint(10), r.FakultasUnit)
	assert.Equal(t, validDate1, r.PeriodeUploadMulai)
	assert.Equal(t, validDate6, r.PeriodeAssesmentLapanganAkhir)
	assert.Equal(t, uint(1), r.Auditee)
	assert.Equal(t, uint(2), r.Auditor1)
	assert.Equal(t, uint(3), r.Auditor2)
}

// -----------------------------------------------------------------------------
// GiveCodeAccessRenstra Tests
// -----------------------------------------------------------------------------
func GiveCodeAccess(t *testing.T) {
	r := &domain.Renstra{
		UUID: uuid.New(),
	}

	// nil prev → error
	r1 := domain.GiveCodeAccessRenstra(nil, uuid.New(), "ABC")
	require.False(t, r1.IsSuccess)
	assert.Equal(t, domain.EmptyData().Code, r1.Error.Code)

	// wrong UUID → error
	r2 := domain.GiveCodeAccessRenstra(r, uuid.New(), "XXX")
	require.False(t, r2.IsSuccess)
	assert.Equal(t, domain.InvalidData().Code, r2.Error.Code)

	// success
	newUUID := r.UUID
	resSuccess := domain.GiveCodeAccessRenstra(r, newUUID, "SECRET")
	require.True(t, resSuccess.IsSuccess)
	assert.NotNil(t, resSuccess.Value.KodeAkses)
	assert.Equal(t, "SECRET", *resSuccess.Value.KodeAkses)
}

// -----------------------------------------------------------------------------
// Test UpdateRenstra - Full Failure Scenarios
// -----------------------------------------------------------------------------
func UpdateRenstra_Errors(t *testing.T) {

	baseUUID := uuid.New()

	// Entity valid sebagai baseline
	prev := &domain.Renstra{
		UUID:                          baseUUID,
		Tahun:                         "2023",
		FakultasUnit:                  1,
		PeriodeUploadMulai:            "2023-01-01",
		PeriodeUploadAkhir:            "2023-01-10",
		PeriodeAssesmentDokumenMulai:  "2023-01-11",
		PeriodeAssesmentDokumenAkhir:  "2023-01-20",
		PeriodeAssesmentLapanganMulai: "2023-01-21",
		PeriodeAssesmentLapanganAkhir: "2023-01-30",
		Auditee:                       10,
		Auditor1:                      20,
		Auditor2:                      30,
	}

	tests := []struct {
		name    string
		call    func() common.ResultValue[*domain.Renstra]
		wantErr common.Error
	}{
		// 1. prev == nil
		{
			name: "PrevNil_ReturnsEmptyData",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					nil, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					11, 22, 33,
				)
			},
			wantErr: domain.EmptyData(),
		},

		// 2. UUID mismatch
		{
			name: "UUIDMismatch_ReturnsInvalidData",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, uuid.New(), "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					11, 22, 33,
				)
			},
			wantErr: domain.InvalidData(),
		},

		// 3. fakultasUnit <= 0
		{
			name: "InvalidFakultasUnit_Zero_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 0,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					11, 22, 33,
				)
			},
			wantErr: domain.InvalidFakultasUnit(),
		},

		// 4. MissingAuditee
		{
			name: "MissingAuditee_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					0, 22, 33,
				)
			},
			wantErr: domain.MissingAuditee(),
		},

		// 5. MissingAuditor1
		{
			name: "MissingAuditor1_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					11, 0, 33,
				)
			},
			wantErr: domain.MissingAuditor1(),
		},

		// 6. MissingAuditor2
		{
			name: "MissingAuditor2_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					11, 22, 0,
				)
			},
			wantErr: domain.MissingAuditor2(),
		},

		// 7. Duplicate assignments
		{
			name: "DuplicateAssignment_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					55, 55, 55,
				)
			},
			wantErr: domain.DuplicateAssigment(),
		},

		// 8. Invalid date parse → upload
		{
			name: "InvalidUploadDate_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"invalid", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					11, 22, 33,
				)
			},
			wantErr: domain.InvalidDate("upload"),
		},
		{
			name: "InvalidUploadDate_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-05", "invalid",
					"2025-01-06", "2025-01-10",
					"2025-01-11", "2025-01-15",
					11, 22, 33,
				)
			},
			wantErr: domain.InvalidDate("upload"),
		},

		// 9. Invalid assessment dokumen
		{
			name: "InvalidDokumenDate_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"invalid", "2025-01-10",
					"2025-01-11", "2025-01-15",
					11, 22, 33,
				)
			},
			wantErr: domain.InvalidDate("assessment dokumen"),
		},
		{
			name: "InvalidDokumenDate_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-10", "invalid",
					"2025-01-11", "2025-01-15",
					11, 22, 33,
				)
			},
			wantErr: domain.InvalidDate("assessment dokumen"),
		},

		// 10. Invalid assessment lapangan
		{
			name: "InvalidLapanganDate_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"invalid", "2025-01-15",
					11, 22, 33,
				)
			},
			wantErr: domain.InvalidDate("assessment lapangan"),
		},
		{
			name: "InvalidLapanganDate_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-06", "2025-01-10",
					"2025-01-15", "invalid",
					11, 22, 33,
				)
			},
			wantErr: domain.InvalidDate("assessment lapangan"),
		},

		// 11. Overlap: Upload overlaps Dokumen
		{
			name: "OverlapUploadDokumen_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-15", // overlap dengan dokumen
					"2025-01-10", "2025-01-20",
					"2025-01-21", "2025-01-30",
					11, 22, 33,
				)
			},
			wantErr: domain.PeriodOverlapUploadDokumen(),
		},

		// 12. Overlap: Upload overlaps Lapangan
		{
			name: "OverlapUploadLapangan_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-10", "2025-01-25",
					"2025-01-01", "2025-01-09",
					"2025-01-20", "2025-01-30",
					11, 22, 33,
				)
			},
			wantErr: domain.PeriodOverlapUploadLapangan(),
		},

		// 13. Overlap: Dokumen overlaps Lapangan
		{
			name: "OverlapDokumenLapangan_ReturnsError",
			call: func() common.ResultValue[*domain.Renstra] {
				return domain.UpdateRenstra(
					prev, baseUUID, "2025", 1,
					"2025-01-01", "2025-01-05",
					"2025-01-10", "2025-01-25",
					"2025-01-20", "2025-01-30",
					11, 22, 33,
				)
			},
			wantErr: domain.PeriodOverlapDokumenLapangan(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res := tt.call()

			require.False(t, res.IsSuccess)

			assert.Equal(t,
				tt.wantErr.Code,
				res.Error.Code,
				"Error code mismatch",
			)
		})
	}
}

// -----------------------------------------------------------------------------
// Test UpdateRenstra Success
// -----------------------------------------------------------------------------
func UpdateRenstra_Success(t *testing.T) {

	// PREPARE EXISTING ENTITY
	prevUUID := uuid.New()
	prev := &domain.Renstra{
		UUID:                          prevUUID,
		Tahun:                         "2023",
		FakultasUnit:                  1,
		PeriodeUploadMulai:            "2023-01-01",
		PeriodeUploadAkhir:            "2023-01-10",
		PeriodeAssesmentDokumenMulai:  "2023-01-11",
		PeriodeAssesmentDokumenAkhir:  "2023-01-20",
		PeriodeAssesmentLapanganMulai: "2023-01-21",
		PeriodeAssesmentLapanganAkhir: "2023-01-30",
		Auditee:                       10,
		Auditor1:                      20,
		Auditor2:                      30,
	}

	// INPUT BARU (VALID SEMUA)
	res := domain.UpdateRenstra(
		prev,
		prevUUID,
		"2025",
		2,
		"2025-02-01", "2025-02-10",
		"2025-02-11", "2025-02-20",
		"2025-02-21", "2025-02-28",
		100, 200, 300,
	)

	// === ASSERT RESULT SUCCESS ===
	if !res.IsSuccess {
		t.Fatalf("expected success, got error: %v", res.Error)
	}

	updated := res.Value

	// === ASSERT FIELDS UPDATED ===
	if updated.Tahun != "2025" {
		t.Errorf("expected tahun=2025, got %s", updated.Tahun)
	}
	if updated.FakultasUnit != 2 {
		t.Errorf("expected fakultas_unit=2, got %d", updated.FakultasUnit)
	}
	if updated.PeriodeUploadMulai != "2025-02-01" {
		t.Errorf("expected PeriodeUploadMulai updated")
	}
	if updated.PeriodeAssesmentLapanganAkhir != "2025-02-28" {
		t.Errorf("expected PeriodeAssesmentLapanganAkhir updated")
	}

	if updated.Auditee != 100 || updated.Auditor1 != 200 || updated.Auditor2 != 300 {
		t.Errorf("auditee/auditor fields not updated")
	}
}
