package domaintest

import (
	"testing"
	"time"

	domain "UnpakSiamida/modules/jadwalproker/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// CREATE SUCCESS
// ====================
func TestNewJadwalProker_Success(t *testing.T) {
	fakultasUnit := uint(1)
	tutupEntry := "2026-02-01"
	tutupDokumen := "2026-02-02"

	res := domain.NewJadwalProker(fakultasUnit, tutupEntry, tutupDokumen)
	require.True(t, res.IsSuccess)
	jadwal := res.Value
	require.NotNil(t, jadwal)

	assert.Equal(t, fakultasUnit, jadwal.FakultasUnit)
	assert.NotEqual(t, uuid.Nil, jadwal.UUID)

	expectedEntry, _ := time.Parse("2006-01-02", tutupEntry)
	expectedDokumen, _ := time.Parse("2006-01-02", tutupDokumen)
	assert.Equal(t, expectedEntry, jadwal.TanggalTutupEntry)
	assert.Equal(t, expectedDokumen, jadwal.TanggalTutupDokumen)
}

// ====================
// CREATE FAIL CASES
// ====================
func TestNewJadwalProker_Fail(t *testing.T) {
	tests := []struct {
		name           string
		fakultasunit   uint
		entry          string
		dokumen        string
		expectedErrMsg string
	}{
		{
			name:           "InvalidTanggalEntry",
			fakultasunit:   1,
			entry:          "2026-02-31", // tanggal salah
			dokumen:        "2026-02-02",
			expectedErrMsg: domain.InvalidDate("tanggal input").Code,
		},
		{
			name:           "InvalidTanggalDokumen",
			fakultasunit:   1,
			entry:          "2026-02-01",
			dokumen:        "2026-02-30", // tanggal salah
			expectedErrMsg: domain.InvalidDate("tanggal upload dokumen").Code,
		},
		{
			name:           "TanggalOverlap",
			fakultasunit:   1,
			entry:          "2026-02-02",
			dokumen:        "2026-02-01", // dokumen sebelum entry
			expectedErrMsg: domain.InvalidDateRange().Code,
		},
		{
			name:           "NotFoundFakultas",
			fakultasunit:   0,
			entry:          "2026-02-01",
			dokumen:        "2026-02-02", // dokumen sebelum entry
			expectedErrMsg: domain.NotFoundFakultas().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewJadwalProker(tt.fakultasunit, tt.entry, tt.dokumen)
			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.expectedErrMsg, res.Error.Code)
		})
	}
}

// ====================
// UPDATE SUCCESS
// ====================
func TestUpdateJadwalProker_Success(t *testing.T) {
	res := domain.NewJadwalProker(1, "2026-02-01", "2026-02-02")
	require.True(t, res.IsSuccess)
	prev := res.Value

	newEntry := "2026-03-01"
	newDokumen := "2026-03-02"
	updateRes := domain.UpdateJadwalProker(prev, prev.UUID, uint(2), newEntry, newDokumen)
	require.True(t, updateRes.IsSuccess)

	updated := updateRes.Value
	assert.Equal(t, uint(2), updated.FakultasUnit)
	assert.Equal(t, prev.UUID, updated.UUID)

	expectedEntry, _ := time.Parse("2006-01-02", newEntry)
	expectedDokumen, _ := time.Parse("2006-01-02", newDokumen)
	assert.Equal(t, expectedEntry, updated.TanggalTutupEntry)
	assert.Equal(t, expectedDokumen, updated.TanggalTutupDokumen)
}

// ====================
// UPDATE FAIL CASES
// ====================
func TestUpdateJadwalProker_Fail(t *testing.T) {
	prev := &domain.JadwalProker{
		UUID: uuid.New(),
	}

	tests := []struct {
		name           string
		prev           *domain.JadwalProker
		uid            uuid.UUID
		fakultasunit   uint
		entry          string
		dokumen        string
		expectedErrMsg string
	}{
		{
			name:           "PrevNil",
			prev:           nil,
			uid:            uuid.New(),
			fakultasunit:   1,
			entry:          "2026-02-01",
			dokumen:        "2026-02-02",
			expectedErrMsg: domain.EmptyData().Code,
		},
		{
			name:           "UUIDMismatch",
			prev:           prev,
			uid:            uuid.New(),
			fakultasunit:   1,
			entry:          "2026-02-01",
			dokumen:        "2026-02-02",
			expectedErrMsg: domain.InvalidData().Code,
		},
		{
			name:           "InvalidTanggalEntry",
			prev:           prev,
			uid:            prev.UUID,
			fakultasunit:   1,
			entry:          "2026-02-30",
			dokumen:        "2026-03-01",
			expectedErrMsg: domain.InvalidDate("tanggal input").Code,
		},
		{
			name:           "InvalidTanggalDokumen",
			prev:           prev,
			uid:            prev.UUID,
			fakultasunit:   1,
			entry:          "2026-02-01",
			dokumen:        "2026-02-30",
			expectedErrMsg: domain.InvalidDate("tanggal upload dokumen").Code,
		},
		{
			name:           "TanggalOverlap",
			prev:           prev,
			uid:            prev.UUID,
			fakultasunit:   1,
			entry:          "2026-03-02",
			dokumen:        "2026-03-01",
			expectedErrMsg: domain.InvalidDateRange().Code,
		},
		{
			name:           "NotFoundFakultas",
			prev:           prev,
			uid:            prev.UUID,
			fakultasunit:   0,
			entry:          "2026-03-01",
			dokumen:        "2026-03-02",
			expectedErrMsg: domain.NotFoundFakultas().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateJadwalProker(tt.prev, tt.uid, tt.fakultasunit, tt.entry, tt.dokumen)
			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.expectedErrMsg, res.Error.Code)
		})
	}
}
