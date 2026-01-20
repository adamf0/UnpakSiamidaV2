package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/kts/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// KTS ERROR TESTS
// ====================
func TestKtsErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "Kts.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "ExistData_ReturnsCorrectError",
			err:          domain.ExistData(),
			expectedCode: "Kts.ExistData",
			expectedDesc: "data is exist",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "Kts.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "Kts.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "InvalidTahun_ReturnsCorrectError",
			err:          domain.InvalidTahun(),
			expectedCode: "Kts.InvalidTahun",
			expectedDesc: "tahun is invalid",
		},
		{
			name:         "InvalidAuditor_ReturnsCorrectError",
			err:          domain.InvalidAuditor(),
			expectedCode: "Kts.InvalidAuditor",
			expectedDesc: "auditor is invalid",
		},
		{
			name:         "InvalidAuditee_ReturnsCorrectError",
			err:          domain.InvalidAuditee(),
			expectedCode: "Kts.InvalidAuditee",
			expectedDesc: "auditee is invalid",
		},
		{
			name:         "InvalidStatusAcc_ReturnsCorrectError",
			err:          domain.InvalidStatusAcc(),
			expectedCode: "Kts.InvalidStatusAcc",
			expectedDesc: "auditor is invalid",
		},
		{
			name:         "InvalidTanggal_ReturnsCorrectError",
			err:          domain.InvalidTanggal(),
			expectedCode: "Kts.InvalidTanggal",
			expectedDesc: "tanggal is invalid",
		},
		{
			name:         "InvalidStep_ReturnsCorrectError",
			err:          domain.InvalidStep(),
			expectedCode: "Kts.InvalidStep",
			expectedDesc: "tanggal is invalid",
		},
		{
			name:         "RequiredNomorLaporan_ReturnsCorrectError",
			err:          domain.RequiredNomorLaporan(),
			expectedCode: "Kts.RequiredNomorLaporan",
			expectedDesc: "nomor laporan is required",
		},
		{
			name:         "RequiredKeteranganTolak_ReturnsCorrectError",
			err:          domain.RequiredKeteranganTolak(),
			expectedCode: "Kts.RequiredKeteranganTolak",
			expectedDesc: "keterangan tolak is required",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("ABC123"),
			expectedCode: "Kts.NotFound",
			expectedDesc: "Kts with identifier ABC123 not found",
		},
		{
			name:         "NotFoundUser_ReturnsCorrectError",
			err:          domain.NotFound("ABC123"),
			expectedCode: "Kts.NotFoundUser",
			expectedDesc: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.err)
			assert.Equal(t, tt.expectedCode, tt.err.Code)
			assert.Equal(t, tt.expectedDesc, tt.err.Description)
		})
	}
}
