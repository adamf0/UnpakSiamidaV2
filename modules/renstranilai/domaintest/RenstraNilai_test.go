package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	domain "UnpakSiamida/modules/renstranilai/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// UPDATE SUCCESS
// ====================
func TestUpdateRenstraNilai_Success_Auditee(t *testing.T) {
	renstraUUID := uuid.New()
	nilaiUUID := uuid.New()

	renstra := &domainrenstra.Renstra{
		UUID:  renstraUUID,
		Tahun: "2026",
	}

	prev := &domain.RenstraNilai{
		UUID:    nilaiUUID,
		Renstra: 1,
	}

	capaian := "Capaian auditee"
	catatan := "Catatan auditee"
	link := "http://bukti"

	res := domain.UpdateRenstraNilai(
		prev,
		renstra,
		nilaiUUID,
		renstraUUID,
		"2026",
		"auditee",
		"2026#auditee",
		&capaian,
		&catatan,
		&link,
		nil,
		nil,
	)

	require.True(t, res.IsSuccess)
	require.NotNil(t, res.Value)

	assert.Equal(t, capaian, *res.Value.Capaian)
	assert.Equal(t, catatan, *res.Value.Catatan)
	assert.Equal(t, link, *res.Value.LinkBukti)
}

// ====================
// UPDATE SUCCESS AUDITOR
// ====================
func TestUpdateRenstraNilai_Success_Auditor(t *testing.T) {
	renstraUUID := uuid.New()
	nilaiUUID := uuid.New()

	renstra := &domainrenstra.Renstra{
		UUID:  renstraUUID,
		Tahun: "2026",
	}

	prev := &domain.RenstraNilai{
		UUID: nilaiUUID,
	}

	capaianAuditor := "Capaian auditor"
	catatanAuditor := "Catatan auditor"

	res := domain.UpdateRenstraNilai(
		prev,
		renstra,
		nilaiUUID,
		renstraUUID,
		"2026",
		"auditor1",
		"2026#auditor1",
		nil,
		nil,
		nil,
		&capaianAuditor,
		&catatanAuditor,
	)

	require.True(t, res.IsSuccess)
	assert.Equal(t, capaianAuditor, *res.Value.CapaianAuditor)
	assert.Equal(t, catatanAuditor, *res.Value.CatatanAuditor)
}

// ====================
// UPDATE FAIL CASES
// ====================
func TestUpdateRenstraNilai_Fail(t *testing.T) {
	validUUID := uuid.New()
	otherUUID := uuid.New()

	validRenstra := &domainrenstra.Renstra{
		UUID:  validUUID,
		Tahun: "2026",
	}

	validPrev := &domain.RenstraNilai{
		UUID: validUUID,
	}

	tests := []struct {
		name    string
		prev    *domain.RenstraNilai
		renstra *domainrenstra.Renstra
		uuid    uuid.UUID
		rUUID   uuid.UUID
		tahun   string
		mode    string
		granted string
		wantErr common.Error
	}{
		{
			name:    "RenstraNil",
			prev:    validPrev,
			renstra: nil,
			wantErr: domain.InvalidRenstra(),
		},
		{
			name:    "PrevNil",
			prev:    nil,
			renstra: validRenstra,
			wantErr: domain.EmptyData(),
		},
		{
			name:    "UUIDMismatch",
			prev:    validPrev,
			renstra: validRenstra,
			uuid:    otherUUID,
			rUUID:   validUUID,
			tahun:   "2026",
			wantErr: domain.InvalidData(),
		},
		{
			name:    "TahunMismatch",
			prev:    validPrev,
			renstra: validRenstra,
			uuid:    validUUID,
			rUUID:   validUUID,
			tahun:   "2025",
			wantErr: domain.InvalidData(),
		},
		{
			name:    "InvalidMode",
			prev:    validPrev,
			renstra: validRenstra,
			uuid:    validUUID,
			rUUID:   validUUID,
			tahun:   "2026",
			mode:    "admin",
			wantErr: domain.RejectAction(),
		},
		{
			name:    "NotGranted",
			prev:    validPrev,
			renstra: validRenstra,
			uuid:    validUUID,
			rUUID:   validUUID,
			tahun:   "2026",
			mode:    "auditee",
			granted: "2026#auditor1",
			wantErr: domain.NotGranted(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateRenstraNilai(
				tt.prev,
				tt.renstra,
				tt.uuid,
				tt.rUUID,
				tt.tahun,
				tt.mode,
				tt.granted,
				nil,
				nil,
				nil,
				nil,
				nil,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.wantErr.Code, res.Error.Code)
			assert.Equal(t, tt.wantErr.Description, res.Error.Description)
		})
	}
}

// ====================
// EDGE CASES
// ====================
func TestUpdateRenstraNilai_EdgeCases(t *testing.T) {
	renstraUUID := uuid.New()
	prevUUID := uuid.New()

	renstra := &domainrenstra.Renstra{
		UUID:  renstraUUID,
		Tahun: "2026",
	}

	prev := &domain.RenstraNilai{
		UUID: prevUUID,
	}

	// granted kosong
	res := domain.UpdateRenstraNilai(
		prev,
		renstra,
		prevUUID,
		renstraUUID,
		"2026",
		"auditee",
		"",
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	require.False(t, res.IsSuccess)
	assert.Equal(t, domain.NotGranted().Code, res.Error.Code)
}
