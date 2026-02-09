package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/beritaacara/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBeritaAcaraErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "BeritaAcara.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidTanggal_ReturnsCorrectError",
			err:          domain.InvalidTanggal(),
			expectedCode: "BeritaAcara.InvalidTanggal",
			expectedDesc: "tanggal is invalid",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "BeritaAcara.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "BeritaAcara.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "InvalidFakultasUnit_ReturnsCorrectError",
			err:          domain.InvalidFakultasUnit(),
			expectedCode: "BeritaAcara.InvalidFakultasUnit",
			expectedDesc: "fakultas is invalid",
		},
		{
			name:         "InvalidAuditee_ReturnsCorrectError",
			err:          domain.InvalidAuditee(),
			expectedCode: "BeritaAcara.InvalidAuditee",
			expectedDesc: "auditee is invalid",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("XYZ99"),
			expectedCode: "BeritaAcara.NotFound",
			expectedDesc: "BeritaAcara with identifier XYZ99 not found",
		},
		{
			name:         "NotFoundFakultas_ReturnsCorrectError",
			err:          domain.NotFoundFakultas(),
			expectedCode: "BeritaAcara.NotFoundFakultas",
			expectedDesc: "fakultas not found",
		},
		{
			name:         "NotFoundAuditee_ReturnsCorrectError",
			err:          domain.NotFoundAuditee(),
			expectedCode: "BeritaAcara.NotFoundAuditee",
			expectedDesc: "auditee not found",
		},
		{
			name:         "NotFoundAuditor_ReturnsCorrectError",
			err:          domain.NotFoundAuditor(),
			expectedCode: "BeritaAcara.NotFoundAuditor",
			expectedDesc: "auditor not found",
		},
		{
			name:         "DuplicateAssigment_ReturnsCorrectError",
			err:          domain.DuplicateAssigment(),
			expectedCode: "BeritaAcara.DuplicateAssigment",
			expectedDesc: "auditee, auditee 1, and auditor 2 must not have the same target",
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
