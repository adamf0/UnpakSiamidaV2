package domaintest

import (
	"testing"

	domain "UnpakSiamida/modules/tahunproker/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// CREATE SUCCESS
// ====================
func TestNewTahunProker_Success(t *testing.T) {
	tahun := "2080"
	status := "aktif"
	res := domain.NewTahunProker(tahun, status)

	require.True(t, res.IsSuccess)
	tahunproker := res.Value
	require.NotNil(t, tahunproker)

	assert.Equal(t, tahun, tahunproker.Tahun)
	assert.Equal(t, status, tahunproker.Status)
	assert.NotEqual(t, uuid.Nil, tahunproker.UUID)
}

// ====================
// UPDATE SUCCESS
// ====================
func TestUpdateTahunProker_Success(t *testing.T) {
	// buat entity awal
	res := domain.NewTahunProker("2011", "aktif")
	require.True(t, res.IsSuccess)
	prev := res.Value

	newTahun := "2080"
	newStatus := "non-aktif"
	updateRes := domain.UpdateTahunProker(prev, prev.UUID, newTahun, newStatus)

	require.True(t, updateRes.IsSuccess)
	updated := updateRes.Value
	assert.Equal(t, newTahun, updated.Tahun)
	assert.Equal(t, prev.UUID, updated.UUID)
}

// ====================
// UPDATE FAIL CASES
// ====================
func TestUpdateTahunProker_Fail(t *testing.T) {
	prev := &domain.TahunProker{
		UUID:   uuid.New(),
		Tahun:  "2080",
		Status: "non-aktif",
	}

	tests := []struct {
		name    string
		prev    *domain.TahunProker
		uid     uuid.UUID
		tahun   string
		status  string
		wantErr string
	}{
		{
			name:    "PrevNil",
			prev:    nil,
			uid:     prev.UUID,
			tahun:   prev.Tahun,
			status:  prev.Status,
			wantErr: domain.EmptyData().Description,
		},
		{
			name:    "UUIDMismatch",
			prev:    prev,
			uid:     uuid.New(),
			tahun:   prev.Tahun,
			status:  prev.Status,
			wantErr: domain.InvalidData().Description,
		},
		{
			name:    "InvalidTahun",
			prev:    prev,
			uid:     prev.UUID,
			tahun:   "1990",
			status:  prev.Status,
			wantErr: domain.InvalidTahun().Description,
		},
		{
			name:    "InvalidStatus",
			prev:    prev,
			uid:     prev.UUID,
			tahun:   prev.Tahun,
			status:  "no-aktif",
			wantErr: domain.InvalidStatus().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateTahunProker(tt.prev, tt.uid, tt.tahun, tt.status)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.wantErr, res.Error.Description)
		})
	}
}

// ====================
// EDGE CASES
// ====================
// func TestTahunProker_EdgeCases(t *testing.T) {
// 	// 1. Nama kosong
// 	res := domain.NewTahunProker("")
// 	require.True(t, res.IsSuccess)
// 	assert.Equal(t, "", res.Value.Nama)

// 	// 2. UUID Nil saat update (edge, seharusnya fail UUIDMismatch)
// 	prev := res.Value
// 	resUpdate := domain.UpdateTahunProker(prev, uuid.Nil, "Nama Baru")
// 	require.False(t, resUpdate.IsSuccess)
// 	assert.Equal(t, domain.InvalidData().Code, resUpdate.Error.Code)
// }
