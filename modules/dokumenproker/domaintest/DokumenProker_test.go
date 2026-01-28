package domaintest

import (
	"testing"

	domain "UnpakSiamida/modules/dokumenproker/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// CREATE SUCCESS
// ====================
func TestNewDokumenProker_Success(t *testing.T) {
	mataprogram := uint(1)
	fakultas := uint(2)
	jenisDokumen := "PROPOSAL/TOR"
	file := "file.pdf"
	status := "belum_terverifikasi"
	catatan := "catatan test"

	res := domain.NewDokumenProker(
		mataprogram,
		fakultas,
		jenisDokumen,
		file,
		status,
		&catatan,
	)

	require.True(t, res.IsSuccess)

	dp := res.Value
	require.NotNil(t, dp)

	assert.NotEqual(t, uuid.Nil, dp.UUID)
	assert.Equal(t, mataprogram, dp.MataProgram)
	assert.Equal(t, fakultas, dp.FakultasUnit)
	assert.Equal(t, jenisDokumen, dp.JenisDokumen)
	assert.Equal(t, file, dp.File)
	assert.Equal(t, status, dp.Status)
	assert.Equal(t, &catatan, dp.Catatan)
}

// ====================
// CREATE FAIL CASES
// ====================
func TestNewDokumenProker_Fail(t *testing.T) {
	validJenis := "PROPOSAL/TOR"
	validStatus := "belum_terverifikasi"
	file := "file.pdf"

	tests := []struct {
		name    string
		matpro  uint
		fak     uint
		jenis   string
		status  string
		errCode string
	}{
		{
			name:    "InvalidFakultas",
			matpro:  1,
			fak:     0,
			jenis:   validJenis,
			status:  validStatus,
			errCode: domain.NotFoundFakultas().Code,
		},
		{
			name:    "InvalidMataProgram",
			matpro:  0,
			fak:     1,
			jenis:   validJenis,
			status:  validStatus,
			errCode: domain.NotFoundMataProgram().Code,
		},
		{
			name:    "InvalidJenisDokumen",
			matpro:  1,
			fak:     1,
			jenis:   "PDF",
			status:  validStatus,
			errCode: domain.InvalidJenisDokumen().Code,
		},
		{
			name:    "InvalidStatus",
			matpro:  1,
			fak:     1,
			jenis:   validJenis,
			status:  "draft",
			errCode: domain.InvalidStatus().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewDokumenProker(
				tt.matpro,
				tt.fak,
				tt.jenis,
				file,
				tt.status,
				nil,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.errCode, res.Error.Code)
		})
	}
}

// ====================
// UPDATE SUCCESS
// ====================
func TestUpdateDokumenProker_Success(t *testing.T) {
	createMatpro := uint(1)
	createFak := uint(1)
	createJenis := "PROPOSAL/TOR"
	createFile := "old.pdf"
	createStatus := "belum_terverifikasi"

	create := domain.NewDokumenProker(
		createMatpro,
		createFak,
		createJenis,
		createFile,
		createStatus,
		nil,
	)
	require.True(t, create.IsSuccess)

	prev := create.Value

	newMatpro := uint(2)
	newFak := uint(3)
	newJenis := "LAPORAN"
	newFile := "new.pdf"
	newStatus := "terverifikasi"
	newCatatan := "updated"

	update := domain.UpdateDokumenProker(
		prev,
		prev.UUID,
		newMatpro,
		newFak,
		newJenis,
		newFile,
		newStatus,
		&newCatatan,
	)

	require.True(t, update.IsSuccess)

	dp := update.Value
	assert.Equal(t, newMatpro, dp.MataProgram)
	assert.Equal(t, newFak, dp.FakultasUnit)
	assert.Equal(t, newJenis, dp.JenisDokumen)
	assert.Equal(t, newFile, dp.File)
	assert.Equal(t, newStatus, dp.Status)
	assert.Equal(t, &newCatatan, dp.Catatan)
}

// ====================
// UPDATE FAIL CASES
// ====================
func TestUpdateDokumenProker_Fail(t *testing.T) {
	validJenis := "PROPOSAL/TOR"
	validStatus := "belum_terverifikasi"
	file := "file.pdf"

	prev := &domain.DokumenProker{
		UUID: uuid.New(),
	}

	tests := []struct {
		name    string
		prev    *domain.DokumenProker
		uid     uuid.UUID
		jenis   string
		status  string
		errCode string
	}{
		{
			name:    "PrevNil",
			prev:    nil,
			uid:     uuid.New(),
			jenis:   validJenis,
			status:  validStatus,
			errCode: domain.EmptyData().Code,
		},
		{
			name:    "UUIDMismatch",
			prev:    prev,
			uid:     uuid.New(),
			jenis:   validJenis,
			status:  validStatus,
			errCode: domain.InvalidData().Code,
		},
		{
			name:    "InvalidJenisDokumen",
			prev:    prev,
			uid:     prev.UUID,
			jenis:   "PDF",
			status:  validStatus,
			errCode: domain.InvalidJenisDokumen().Code,
		},
		{
			name:    "InvalidStatus",
			prev:    prev,
			uid:     prev.UUID,
			jenis:   validJenis,
			status:  "draft",
			errCode: domain.InvalidStatus().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateDokumenProker(
				tt.prev,
				tt.uid,
				1,
				1,
				tt.jenis,
				file,
				tt.status,
				nil,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.errCode, res.Error.Code)
		})
	}
}
