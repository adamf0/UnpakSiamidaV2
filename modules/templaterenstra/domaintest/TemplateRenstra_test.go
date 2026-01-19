package domaintest

import (
	"testing"

	"UnpakSiamida/modules/templaterenstra/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===== CREATE =====
func NewTemplateRenstra_Success(t *testing.T) {
	tahun := "2025"
	indikatorID := uint(1)
	isPertanyaan := true
	fakultasUnit := uint(1)
	kategori := "C"
	klasifikasi := "S1"
	tugas := "Unit"

	result := domain.NewTemplateRenstra(
		tahun,
		indikatorID,
		isPertanyaan,
		fakultasUnit,
		kategori,
		klasifikasi,
		nil,       // satuan
		nil,       // target
		ptr("10"), // targetMin
		ptr("20"), // targetMax
		tugas,
	)

	require.True(t, result.IsSuccess)

	v := result.Value
	assert.Equal(t, tahun, v.Tahun)
	assert.Equal(t, indikatorID, v.IndikatorRenstraID)
	assert.NotEqual(t, uuid.Nil, v.UUID)
}

// ===== CREATE (NEGATIVE) ====
func NewTemplateRenstra_NegativeCases(t *testing.T) {
	tests := []struct {
		name            string
		tahun           string
		indikatorID     uint
		fakultasUnit    uint
		target          *string
		targetMin       *string
		targetMax       *string
		expectedErrDesc string
	}{
		{
			name:  "Invalid target combination (Mode A fail)",
			tahun: "2025", indikatorID: 1, fakultasUnit: 1,
			target: ptr("10"), targetMin: ptr("5"), targetMax: nil,
			expectedErrDesc: domain.InvalidValueTarget().Description,
		},
		{
			name:  "Invalid target combination (Mode B fail)",
			tahun: "2025", indikatorID: 1, fakultasUnit: 1,
			target: nil, targetMin: nil, targetMax: ptr("20"),
			expectedErrDesc: domain.InvalidValueTarget().Description,
		},
		{
			name:  "Invalid indikatorRenstraID",
			tahun: "2025", indikatorID: 0, fakultasUnit: 1,
			target: ptr("10"), targetMin: nil, targetMax: nil,
			expectedErrDesc: domain.IndikatorNotFound().Description,
		},
		{
			name:  "Invalid fakultasUnit",
			tahun: "2025", indikatorID: 1, fakultasUnit: 0,
			target: ptr("10"), targetMin: nil, targetMax: nil,
			expectedErrDesc: domain.FakultasUnitNotFound().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewTemplateRenstra(
				tt.tahun,
				tt.indikatorID,
				true,
				tt.fakultasUnit,
				"kategori",
				"klasifikasi",
				nil, // satuan
				tt.target,
				tt.targetMin,
				tt.targetMax,
				"tugas",
			)
			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.expectedErrDesc, res.Error.Description)
		})
	}
}

// ===== UPDATE (POSITIVE) =====
func UpdateTemplateRenstra_Success(t *testing.T) {
	prev := domain.NewTemplateRenstra(
		"2025", 1, true, 1, "C", "S1",
		nil, nil, ptr("10"), ptr("20"), "Unit",
	).Value

	newTahun := "2026"
	newKategori := "D"
	newKlasifikasi := "S2"

	result := domain.UpdateTemplateRenstra(
		prev,
		prev.UUID,
		newTahun,
		2,     // indikatorRenstraID baru
		false, // isPertanyaan
		1,     // fakultasUnit
		newKategori,
		newKlasifikasi,
		nil,        // satuan
		nil,        // target
		ptr("100"), // targetMin
		ptr("200"), // targetMax
		"Unit",
	)

	require.True(t, result.IsSuccess)
	assert.Equal(t, newTahun, result.Value.Tahun)
	assert.Equal(t, newKategori, result.Value.Kategori)
	assert.Equal(t, newKlasifikasi, result.Value.Klasifikasi)
}

// ===== UPDATE (NEGATIVE) =====
func UpdateTemplateRenstra_NegativeCases(t *testing.T) {
	prev := domain.NewTemplateRenstra(
		"2025", 1, true, 1, "C", "S1",
		nil, nil, ptr("10"), ptr("20"), "Unit",
	).Value

	tests := []struct {
		name            string
		uid             uuid.UUID
		target          *string
		targetMin       *string
		targetMax       *string
		indikatorID     uint
		fakultasUnit    uint
		expectedErrDesc string
	}{
		{
			name: "Invalid target combination",
			uid:  prev.UUID, target: ptr("10"), targetMin: ptr("5"), targetMax: nil,
			indikatorID: 1, fakultasUnit: 1,
			expectedErrDesc: domain.InvalidValueTarget().Description,
		},
		{
			name: "Invalid indikatorRenstraID",
			uid:  prev.UUID, target: ptr("10"), targetMin: nil, targetMax: nil,
			indikatorID: 0, fakultasUnit: 1,
			expectedErrDesc: domain.IndikatorNotFound().Description,
		},
		{
			name: "Invalid fakultasUnit",
			uid:  prev.UUID, target: ptr("10"), targetMin: nil, targetMax: nil,
			indikatorID: 1, fakultasUnit: 0,
			expectedErrDesc: domain.FakultasUnitNotFound().Description,
		},
		{
			name: "UUID mismatch",
			uid:  uuid.New(), target: ptr("10"), targetMin: nil, targetMax: nil,
			indikatorID: 1, fakultasUnit: 1,
			expectedErrDesc: domain.InvalidData().Description,
		},
		{
			name: "Prev is nil",
			uid:  uuid.New(), target: ptr("10"), targetMin: nil, targetMax: nil,
			indikatorID: 1, fakultasUnit: 1,
			expectedErrDesc: domain.EmptyData().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var toUpdate *domain.TemplateRenstra
			if tt.name != "Prev is nil" {
				toUpdate = prev
			} else {
				toUpdate = nil
			}

			res := domain.UpdateTemplateRenstra(
				toUpdate,
				tt.uid,
				"2026",
				tt.indikatorID,
				false,
				tt.fakultasUnit,
				"kategori",
				"klasifikasi",
				nil, // satuan
				tt.target,
				tt.targetMin,
				tt.targetMax,
				"tugas",
			)
			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.expectedErrDesc, res.Error.Description)
		})
	}
}

// ===== HELPER =====
func ptr(s string) *string { return &s }
