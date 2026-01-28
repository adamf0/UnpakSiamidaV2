package domaintest

import (
	"testing"
	"time"

	domain "UnpakSiamida/modules/aktivitasproker/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// CREATE SUCCESS
// ====================
func TestNewAktivitasProker_Success(t *testing.T) {
	mataprogram := uint(1)
	fakultasunit := uint(2)
	aktivitas := "Aktivitas Test"
	pic := "PIC Test"
	tanggalAwal := "2026-02-01"
	tanggalAkhir := "2026-02-02"

	res := domain.NewAktivitasProker(
		mataprogram,
		fakultasunit,
		aktivitas,
		pic,
		tanggalAwal,
		tanggalAkhir,
	)

	require.True(t, res.IsSuccess)

	ap := res.Value
	require.NotNil(t, ap)

	assert.Equal(t, mataprogram, ap.MataProgram)
	assert.Equal(t, fakultasunit, ap.FakultasUnit)
	assert.Equal(t, aktivitas, ap.Aktivitas)
	assert.Equal(t, pic, ap.PIC)
	assert.NotEqual(t, uuid.Nil, ap.UUID)

	expectedAwal, _ := time.Parse("2006-01-02", "2026-02-01")
	expectedAkhir, _ := time.Parse("2006-01-02", "2026-02-02")

	assert.Equal(t, expectedAwal, ap.TanggalRKAwal)
	assert.Equal(t, expectedAkhir, ap.TanggalRKAkhir)
}

// ====================
// CREATE FAIL CASES
// ====================
func TestNewAktivitasProker_Fail(t *testing.T) {
	tests := []struct {
		name     string
		matpro   uint
		fakultas uint
		awal     string
		akhir    string
		errCode  string
	}{
		{
			name:     "InvalidFakultas",
			matpro:   1,
			fakultas: 0,
			awal:     "2026-02-01",
			akhir:    "2026-02-02",
			errCode:  domain.NotFoundFakultas().Code,
		},
		{
			name:     "InvalidMataProgram",
			matpro:   0,
			fakultas: 1,
			awal:     "2026-02-01",
			akhir:    "2026-02-02",
			errCode:  domain.NotFoundMataProgram().Code,
		},
		{
			name:     "InvalidTanggalAwal",
			matpro:   1,
			fakultas: 1,
			awal:     "2026-02-30",
			akhir:    "2026-03-01",
			errCode:  domain.InvalidDate("tanggal rk awal").Code,
		},
		{
			name:     "InvalidTanggalAkhir",
			matpro:   1,
			fakultas: 1,
			awal:     "2026-02-01",
			akhir:    "2026-02-30",
			errCode:  domain.InvalidDate("tanggal rk akhir").Code,
		},
		{
			name:     "TanggalOverlap",
			matpro:   1,
			fakultas: 1,
			awal:     "2026-02-02",
			akhir:    "2026-02-01",
			errCode:  domain.InvalidDateRange().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewAktivitasProker(
				tt.matpro,
				tt.fakultas,
				"Aktivitas",
				"PIC",
				tt.awal,
				tt.akhir,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.errCode, res.Error.Code)
		})
	}
}

// ====================
// UPDATE SUCCESS
// ====================
func TestUpdateAktivitasProker_Success(t *testing.T) {
	create := domain.NewAktivitasProker(
		1, 1, "Old", "Old PIC", "2026-02-01", "2026-02-02",
	)
	require.True(t, create.IsSuccess)

	prev := create.Value

	newMataprogram := uint(2)
	newFakultas := uint(3)
	newAktivitas := "New Aktivitas"
	newPIC := "New PIC"
	newAwal := "2026-03-01"
	newAkhir := "2026-03-02"

	update := domain.UpdateAktivitasProker(
		prev,
		prev.UUID,
		newMataprogram,
		newFakultas,
		newAktivitas,
		newPIC,
		newAwal,
		newAkhir,
	)

	require.True(t, update.IsSuccess)

	ap := update.Value
	assert.Equal(t, newMataprogram, ap.MataProgram)
	assert.Equal(t, newFakultas, ap.FakultasUnit)
	assert.Equal(t, newAktivitas, ap.Aktivitas)
	assert.Equal(t, newPIC, ap.PIC)
}

// ====================
// UPDATE FAIL CASES
// ====================
func TestUpdateAktivitasProker_Fail(t *testing.T) {
	prev := &domain.AktivitasProker{
		UUID: uuid.New(),
	}

	tests := []struct {
		name    string
		prev    *domain.AktivitasProker
		uid     uuid.UUID
		matpro  uint
		fak     uint
		awal    string
		akhir   string
		errDesc string
	}{
		{
			name:    "PrevNil",
			prev:    nil,
			uid:     uuid.New(),
			matpro:  1,
			fak:     1,
			awal:    "2026-02-01",
			akhir:   "2026-02-02",
			errDesc: domain.EmptyData().Code,
		},
		{
			name:    "UUIDMismatch",
			prev:    prev,
			uid:     uuid.New(),
			matpro:  1,
			fak:     1,
			awal:    "2026-02-01",
			akhir:   "2026-02-02",
			errDesc: domain.InvalidData().Code,
		},
		{
			name:    "InvalidTanggalAwal",
			prev:    prev,
			uid:     prev.UUID,
			matpro:  1,
			fak:     1,
			awal:    "2026-02-30",
			akhir:   "2026-03-01",
			errDesc: domain.InvalidDate("tanggal rk awal").Code,
		},
		{
			name:    "OverlapTanggal",
			prev:    prev,
			uid:     prev.UUID,
			matpro:  1,
			fak:     1,
			awal:    "2026-03-02",
			akhir:   "2026-03-01",
			errDesc: domain.InvalidDateRange().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateAktivitasProker(
				tt.prev,
				tt.uid,
				tt.matpro,
				tt.fak,
				"Aktivitas",
				"PIC",
				tt.awal,
				tt.akhir,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.errDesc, res.Error.Code)
		})
	}
}
