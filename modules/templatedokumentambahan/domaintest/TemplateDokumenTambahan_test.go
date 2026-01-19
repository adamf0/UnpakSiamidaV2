package domaintest

import (
	"testing"

	"UnpakSiamida/modules/templatedokumentambahan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper untuk pointer string
func ptr(s string) *string {
	return &s
}

// ====================
// CREATE
// ====================
func TestNewTemplateDokumenTambahan_Success(t *testing.T) {
	result := domain.NewTemplateDokumenTambahan(
		"2025",       // tahun
		1,            // jenisFileID
		"Pertanyaan", // pertanyaan
		"Klasifikasi",
		"Kategori",
		"Tugas",
	)

	require.True(t, result.IsSuccess)

	v := result.Value
	assert.Equal(t, "2025", v.Tahun)
	assert.Equal(t, uint(1), v.JenisFileID)
	assert.Equal(t, "Pertanyaan", v.Pertanyaan)
	assert.Equal(t, "Klasifikasi", v.Klasifikasi)
	assert.Equal(t, "Kategori", v.Kategori)
	assert.Equal(t, "Tugas", v.Tugas)
	assert.NotEqual(t, uuid.Nil, v.UUID)
}

// ====================
// CREATE NEGATIVE (jenisFileID <=0)
// ====================
func TestNewTemplateDokumenTambahan_Fail_InvalidJenisFileID(t *testing.T) {
	result := domain.NewTemplateDokumenTambahan(
		"2025", 0, "Pertanyaan", "Klasifikasi", "Kategori", "Tugas",
	)
	require.False(t, result.IsSuccess)
	assert.Equal(t, domain.JenisFileNotFound().Description, result.Error.Description)
}

// ====================
// UPDATE POSITIVE
// ====================
func TestUpdateTemplateDokumenTambahan_Success(t *testing.T) {
	prev := domain.NewTemplateDokumenTambahan(
		"2025", 1, "Pertanyaan Lama", "Klasifikasi", "Kategori", "Tugas",
	).Value

	result := domain.UpdateTemplateDokumenTambahan(
		prev,
		prev.UUID,
		"2026",       // tahun
		2,            // jenisFileID
		"Pertanyaan Baru",
		"Klasifikasi Baru",
		"Kategori Baru",
		"Tugas Baru",
	)

	require.True(t, result.IsSuccess)
	v := result.Value
	assert.Equal(t, "2026", v.Tahun)
	assert.Equal(t, uint(2), v.JenisFileID)
	assert.Equal(t, "Pertanyaan Baru", v.Pertanyaan)
	assert.Equal(t, "Klasifikasi Baru", v.Klasifikasi)
	assert.Equal(t, "Kategori Baru", v.Kategori)
	assert.Equal(t, "Tugas Baru", v.Tugas)
}

// ====================
// UPDATE NEGATIVE
// ====================
func TestUpdateTemplateDokumenTambahan_NegativeCases(t *testing.T) {
	valid := domain.NewTemplateDokumenTambahan(
		"2025", 1, "Pertanyaan", "Klasifikasi", "Kategori", "Tugas",
	).Value

	tests := []struct {
		name            string
		prev            *domain.TemplateDokumenTambahan
		uid             uuid.UUID
		jenisFileID     uint
		expectedErrDesc string
	}{
		{
			name:            "prev is nil → EmptyData",
			prev:            nil,
			uid:             uuid.New(),
			jenisFileID:     1,
			expectedErrDesc: domain.EmptyData().Description,
		},
		{
			name:            "uuid mismatch → InvalidData",
			prev:            valid,
			uid:             uuid.New(), // beda UUID
			jenisFileID:     1,
			expectedErrDesc: domain.InvalidData().Description,
		},
		{
			name:            "invalid jenisFileID → JenisFileNotFound",
			prev:            valid,
			uid:             valid.UUID,
			jenisFileID:     0,
			expectedErrDesc: domain.JenisFileNotFound().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domain.UpdateTemplateDokumenTambahan(
				tt.prev,
				tt.uid,
				"2026",
				tt.jenisFileID,
				"Pertanyaan Baru",
				"Klasifikasi Baru",
				"Kategori Baru",
				"Tugas Baru",
			)
			require.False(t, result.IsSuccess)
			assert.Equal(t, tt.expectedErrDesc, result.Error.Description)
		})
	}
}
