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
	tahun := "2080"
	fakultasUnit := uint(1)
	tanggal := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	auditee := uint(1)
	auditor1 := uint(2)
	auditor2 := uint(3)

	res := domain.NewBeritaAcara(
		tahun,
		fakultasUnit,
		tanggal,
		auditee,
		&auditor1,
		&auditor2,
	)

	require.True(t, res.IsSuccess)
	ba := res.Value
	require.NotNil(t, ba)

	assert.Equal(t, tahun, ba.Tahun)
	assert.Equal(t, fakultasUnit, ba.FakultasUnit)
	assert.Equal(t, tanggal, ba.Tanggal)
	assert.Equal(t, auditee, ba.Auditee)
	assert.NotEqual(t, uuid.Nil, ba.UUID)
}

// ====================
// CREATE FAIL CASES
// ====================
func TestNewBeritaAcara_Fail(t *testing.T) {
	tahun := "2080"
	tanggal := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)

	validFakultas := uint(1)
	validAuditee := uint(1)
	validAuditor1 := uint(2)
	validAuditor2 := uint(3)

	tests := []struct {
		name         string
		fakultasUnit uint
		auditee      uint
		auditor1     *uint
		auditor2     *uint
		wantErr      string
	}{
		{
			name:         "NotFoundFakultas",
			fakultasUnit: 0,
			auditee:      validAuditee,
			auditor1:     &validAuditor1,
			auditor2:     &validAuditor2,
			wantErr:      domain.NotFoundFakultas().Code,
		},
		{
			name:         "NotFoundAuditee",
			fakultasUnit: validFakultas,
			auditee:      0,
			auditor1:     &validAuditor1,
			auditor2:     &validAuditor2,
			wantErr:      domain.NotFoundAuditee().Code,
		},
		{
			name:         "NotFoundAuditor",
			fakultasUnit: validFakultas,
			auditee:      validAuditee,
			auditor1:     func() *uint { v := uint(0); return &v }(),
			auditor2:     &validAuditor2,
			wantErr:      domain.NotFoundAuditor().Code,
		},
		{
			name:         "DuplicateAssigment",
			fakultasUnit: validFakultas,
			auditee:      validAuditee,
			auditor1:     &validAuditee,
			auditor2:     &validAuditor2,
			wantErr:      domain.DuplicateAssigment().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewBeritaAcara(
				tahun,
				tt.fakultasUnit,
				tanggal,
				tt.auditee,
				tt.auditor1,
				tt.auditor2,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.wantErr, res.Error.Code)
		})
	}
}

// ====================
// UPDATE SUCCESS
// ====================
func TestUpdateBeritaAcara_Success(t *testing.T) {
	tahun := "2080"
	fakultasUnit := uint(1)
	tanggal := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	auditee := uint(1)
	auditor1 := uint(2)
	auditor2 := uint(3)

	createRes := domain.NewBeritaAcara(
		tahun,
		fakultasUnit,
		tanggal,
		auditee,
		&auditor1,
		&auditor2,
	)
	require.True(t, createRes.IsSuccess)

	prev := createRes.Value

	newTahun := "2090"
	newFakultasUnit := uint(2)
	newTanggal := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	newAuditee := uint(10)
	newAuditor1 := uint(20)
	newAuditor2 := uint(30)

	updateRes := domain.UpdateBeritaAcara(
		prev,
		prev.UUID,
		newTahun,
		newFakultasUnit,
		newTanggal,
		newAuditee,
		&newAuditor1,
		&newAuditor2,
	)

	require.True(t, updateRes.IsSuccess)
	updated := updateRes.Value

	assert.Equal(t, newTahun, updated.Tahun)
	assert.Equal(t, newFakultasUnit, updated.FakultasUnit)
	assert.Equal(t, newTanggal, updated.Tanggal)
	assert.Equal(t, newAuditee, updated.Auditee)
	assert.Equal(t, prev.UUID, updated.UUID)
}

// ====================
// UPDATE FAIL CASES
// ====================
func TestUpdateBeritaAcara_Fail(t *testing.T) {
	tahun := "2080"
	fakultasUnit := uint(1)
	tanggal := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	auditee := uint(1)
	auditor1 := uint(2)
	auditor2 := uint(3)

	prev := &domain.BeritaAcara{
		UUID:         uuid.New(),
		Tahun:        tahun,
		FakultasUnit: fakultasUnit,
		Tanggal:      tanggal,
		Auditee:      auditee,
		Auditor1:     &auditor1,
		Auditor2:     &auditor2,
	}

	tests := []struct {
		name         string
		prev         *domain.BeritaAcara
		uid          uuid.UUID
		fakultasUnit uint
		auditee      uint
		auditor1     *uint
		auditor2     *uint
		wantErr      string
	}{
		{
			name:         "EmptyData",
			prev:         nil,
			uid:          uuid.New(),
			fakultasUnit: fakultasUnit,
			auditee:      auditee,
			auditor1:     &auditor1,
			auditor2:     &auditor2,
			wantErr:      domain.EmptyData().Code,
		},
		{
			name:         "InvalidData",
			prev:         prev,
			uid:          uuid.New(),
			fakultasUnit: fakultasUnit,
			auditee:      auditee,
			auditor1:     &auditor1,
			auditor2:     &auditor2,
			wantErr:      domain.InvalidData().Code,
		},
		{
			name:         "NotFoundFakultas",
			prev:         prev,
			uid:          prev.UUID,
			fakultasUnit: 0,
			auditee:      auditee,
			auditor1:     &auditor1,
			auditor2:     &auditor2,
			wantErr:      domain.NotFoundFakultas().Code,
		},
		{
			name:         "NotFoundAuditee",
			prev:         prev,
			uid:          prev.UUID,
			fakultasUnit: fakultasUnit,
			auditee:      0,
			auditor1:     &auditor1,
			auditor2:     &auditor2,
			wantErr:      domain.NotFoundAuditee().Code,
		},
		{
			name:         "NotFoundAuditor",
			prev:         prev,
			uid:          prev.UUID,
			fakultasUnit: fakultasUnit,
			auditee:      auditee,
			auditor1:     func() *uint { v := uint(0); return &v }(),
			auditor2:     &auditor2,
			wantErr:      domain.NotFoundAuditor().Code,
		},
		{
			name:         "DuplicateAssigment",
			prev:         prev,
			uid:          prev.UUID,
			fakultasUnit: fakultasUnit,
			auditee:      auditee,
			auditor1:     &auditee,
			auditor2:     &auditor2,
			wantErr:      domain.DuplicateAssigment().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateBeritaAcara(
				tt.prev,
				tt.uid,
				tahun,
				tt.fakultasUnit,
				tanggal,
				tt.auditee,
				tt.auditor1,
				tt.auditor2,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.wantErr, res.Error.Code)
		})
	}
}
