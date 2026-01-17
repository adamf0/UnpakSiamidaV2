package domaintest

import (
	"testing"

	domain "UnpakSiamida/modules/jenisfile/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// CREATE SUCCESS
// ====================
func TestNewJenisFile_Success(t *testing.T) {
	nama := "Dokumen Pendukung"
	res := domain.NewJenisFile(nama)

	require.True(t, res.IsSuccess)
	jenisfile := res.Value
	require.NotNil(t, jenisfile)

	assert.Equal(t, nama, jenisfile.Nama)
	assert.NotEqual(t, uuid.Nil, jenisfile.UUID)
}

// ====================
// UPDATE SUCCESS
// ====================
func TestUpdateJenisFile_Success(t *testing.T) {
	// buat entity awal
	res := domain.NewJenisFile("Awal")
	require.True(t, res.IsSuccess)
	prev := res.Value

	newNama := "Update Nama"
	updateRes := domain.UpdateJenisFile(prev, prev.UUID, newNama)

	require.True(t, updateRes.IsSuccess)
	updated := updateRes.Value
	assert.Equal(t, newNama, updated.Nama)
	assert.Equal(t, prev.UUID, updated.UUID)
}

// ====================
// UPDATE FAIL CASES
// ====================
func TestUpdateJenisFile_Fail(t *testing.T) {
	prev := &domain.JenisFile{
		UUID: uuid.New(),
		Nama: "Awal",
	}

	tests := []struct {
		name    string
		prev    *domain.JenisFile
		uid     uuid.UUID
		nama    string
		wantErr string
	}{
		{
			name:    "PrevNil",
			prev:    nil,
			uid:     uuid.New(),
			nama:    "Test",
			wantErr: domain.EmptyData().Description,
		},
		{
			name:    "UUIDMismatch",
			prev:    prev,
			uid:     uuid.New(),
			nama:    "Test",
			wantErr: domain.InvalidData().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateJenisFile(tt.prev, tt.uid, tt.nama)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.wantErr, res.Error.Description)
		})
	}
}

// ====================
// EDGE CASES
// ====================
func TestJenisFile_EdgeCases(t *testing.T) {
	// 1. Nama kosong
	res := domain.NewJenisFile("")
	require.True(t, res.IsSuccess)
	assert.Equal(t, "", res.Value.Nama)

	// 2. UUID Nil saat update (edge, seharusnya fail UUIDMismatch)
	prev := res.Value
	resUpdate := domain.UpdateJenisFile(prev, uuid.Nil, "Nama Baru")
	require.False(t, resUpdate.IsSuccess)
	assert.Equal(t, domain.InvalidData().Code, resUpdate.Error.Code)
}
