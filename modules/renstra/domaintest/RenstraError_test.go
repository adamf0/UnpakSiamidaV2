package domaintest

import (
	"testing"

	"UnpakSiamida/modules/renstra/domain"
	common "UnpakSiamida/common/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenstraErrors(t *testing.T) {

	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "Renstra.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "Renstra.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "Renstra.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFound_ReturnsCorrectError",
			err:          domain.NotFound("ABC123"),
			expectedCode: "Renstra.NotFound",
			expectedDesc: "Renstra with identifier ABC123 not found",
		},
		{
			name:         "InvalidFakultasUnit_ReturnsCorrectError",
			err:          domain.InvalidFakultasUnit(),
			expectedCode: "Renstra.InvalidFakultasUnit",
			expectedDesc: "fakultas unit is invalid",
		},
		{
			name:         "DataExisting_ReturnsCorrectError",
			err:          domain.DataExisting(),
			expectedCode: "Renstra.DataExisting",
			expectedDesc: "the audit data form already existed before",
		},
		{
			name:         "MissingAuditee_ReturnsCorrectError",
			err:          domain.MissingAuditee(),
			expectedCode: "Renstra.MissingAuditee",
			expectedDesc: "auditee have not been assigned",
		},
		{
			name:         "MissingAuditor1_ReturnsCorrectError",
			err:          domain.MissingAuditor1(),
			expectedCode: "Renstra.MissingAuditor1",
			expectedDesc: "auditor1 have not been assigned",
		},
		{
			name:         "MissingAuditor2_ReturnsCorrectError",
			err:          domain.MissingAuditor2(),
			expectedCode: "Renstra.MissingAuditor2",
			expectedDesc: "auditor2 have not been assigned",
		},
		{
			name:         "InvalidParsing_ReturnsCorrectError",
			err:          domain.InvalidParsing("tahun"),
			expectedCode: "Renstra.IvalidParsing", // sesuai domain-mu (ada typo "IvalidParsing")
			expectedDesc: "failed parsing tahun to UUID",
		},
		{
			name:         "DuplicateAssigment_ReturnsCorrectError",
			err:          domain.DuplicateAssigment(),
			expectedCode: "Renstra.DuplicateAssigment",
			expectedDesc: "auditee, auditee 1, and auditor 2 must not have the same target",
		},
		{
			name:         "InvalidDate_ReturnsCorrectError",
			err:          domain.InvalidDate("upload"),
			expectedCode: "Renstra.InvalidDate",
			expectedDesc: "upload period have wrong date format",
		},
		{
			name:         "PeriodOverlapUploadDokumen_ReturnsCorrectError",
			err:          domain.PeriodOverlapUploadDokumen(),
			expectedCode: "Renstra.PeriodOverlapUploadDokumen",
			expectedDesc: "upload period overlaps with document period",
		},
		{
			name:         "PeriodOverlapUploadLapangan_ReturnsCorrectError",
			err:          domain.PeriodOverlapUploadLapangan(),
			expectedCode: "Renstra.PeriodOverlapUploadLapangan",
			expectedDesc: "upload period overlaps with AL period",
		},
		{
			name:         "PeriodOverlapDokumenLapangan_ReturnsCorrectError",
			err:          domain.PeriodOverlapDokumenLapangan(),
			expectedCode: "Renstra.PeriodOverlapDokumenLapangan",
			expectedDesc: "document period overlaps with AL period",
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
