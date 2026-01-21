package domaintest

import (
	"testing"
	"time"

	domain "UnpakSiamida/modules/beritaacara/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// CREATE SUCCESS
// ====================
func TestNewBeritaAcara_Success(t *testing.T) {
	Tahun := "2080"
	FakultasUnit := 1
	Tanggal := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	res := domain.NewBeritaAcara(
		Tahun,
		FakultasUnit,
		Tanggal,
		&Auditee,
		&Auditor1,
		&Auditor2,
	)

	require.True(t, res.IsSuccess)
	beritaacara := res.Value
	require.NotNil(t, beritaacara)

	assert.Equal(t, Tahun, beritaacara.Tahun)
	assert.Equal(t, FakultasUnit, beritaacara.FakultasUnit)
	assert.Equal(t, Tanggal, beritaacara.Tanggal)
	assert.NotEqual(t, uuid.Nil, beritaacara.UUID)
}

// ====================
// UPDATE SUCCESS
// ====================
func TestUpdateBeritaAcara_Success(t *testing.T) {
	// buat entity awal
	Tahun := "2080"
	FakultasUnit := 1
	Tanggal := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	res := domain.NewBeritaAcara(
		Tahun,
		FakultasUnit,
		Tanggal,
		&Auditee,
		&Auditor1,
		&Auditor2,
	)
	require.True(t, res.IsSuccess)
	prev := res.Value

	newTahun := "2080"
	newFakultasUnit := 2
	newTanggal := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	newAuditee := 10
	newAuditor1 := 20
	newAuditor2 := 30
	updateRes := domain.UpdateBeritaAcara(
		prev,
		prev.UUID,
		newTahun,
		newFakultasUnit,
		newTanggal,
		&newAuditee,
		&newAuditor1,
		&newAuditor2,
	)

	require.True(t, updateRes.IsSuccess)
	updated := updateRes.Value
	assert.Equal(t, newTahun, updated.Tahun)
	assert.Equal(t, newFakultasUnit, updated.FakultasUnit)
	assert.Equal(t, newTanggal, updated.Tanggal)
	assert.Equal(t, prev.UUID, updated.UUID)
}

// ====================
// UPDATE FAIL CASES
// ====================
func TestUpdateBeritaAcara_Fail(t *testing.T) {
	Tahun := "2080"
	FakultasUnit := 1
	Tanggal := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	prev := &domain.BeritaAcara{
		UUID:         uuid.New(),
		Tahun:        Tahun,
		FakultasUnit: FakultasUnit,
		Tanggal:      Tanggal,
		Auditee:      &Auditee,
		Auditor1:     &Auditor1,
		Auditor2:     &Auditor2,
	}

	tests := []struct {
		name         string
		prev         *domain.BeritaAcara
		uid          uuid.UUID
		tahun        string
		fakultasUnit int
		tanggal      time.Time
		auditee      *int
		auditor1     *int
		auditor2     *int
		wantErr      string
	}{
		{
			name:         "PrevNil",
			prev:         nil,
			uid:          uuid.New(),
			tahun:        Tahun,
			fakultasUnit: FakultasUnit,
			tanggal:      Tanggal,
			auditee:      &Auditee,
			auditor1:     &Auditor1,
			auditor2:     &Auditor2,
			wantErr:      domain.EmptyData().Description,
		},
		{
			name:         "UUIDMismatch",
			prev:         prev,
			uid:          uuid.New(),
			tahun:        Tahun,
			fakultasUnit: FakultasUnit,
			tanggal:      Tanggal,
			auditee:      &Auditee,
			auditor1:     &Auditor1,
			auditor2:     &Auditor2,
			wantErr:      domain.InvalidData().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateBeritaAcara(
				tt.prev,
				tt.uid,
				tt.tahun,
				tt.fakultasUnit,
				tt.tanggal,
				tt.auditee,
				tt.auditor1,
				tt.auditor2,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.wantErr, res.Error.Description)
		})
	}
}
