package domaintest

import (
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/generaterenstra/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRenstraErrors_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData",
			err:          domain.EmptyData(),
			expectedCode: "GenerateRenstra.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid",
			err:          domain.InvalidUuid(),
			expectedCode: "GenerateRenstra.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData",
			err:          domain.InvalidData(),
			expectedCode: "GenerateRenstra.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFoundDynamicId",
			err:          domain.NotFound("12345"),
			expectedCode: "GenerateRenstra.NotFound",
			expectedDesc: "GenerateRenstra with identifier 12345 not found",
		},
		{
			name:         "NotFoundFakultasUnit",
			err:          domain.NotFoundFakultasUnit("F01"),
			expectedCode: "GenerateRenstra.NotFoundFakultasUnit",
			expectedDesc: "Fakultas unit with identifier F01 not found",
		},
		{
			name:         "NotFoundRenstra",
			err:          domain.NotFoundRenstra("REN123"),
			expectedCode: "GenerateRenstra.NotFoundRenstra",
			expectedDesc: "renstra with identifier REN123 not found",
		},
		{
			name:         "NotFoundTemplate",
			err:          domain.NotFoundTemplate("2026", "F01"),
			expectedCode: "GenerateRenstra.NotFoundTemplate",
			expectedDesc: "template not found with identifier tahun 2026 & fakultas F01",
		},
		{
			name:         "NotFoundAudit",
			err:          domain.NotFoundAudit("2026", "F01"),
			expectedCode: "GenerateRenstra.NotFoundAudit",
			expectedDesc: "previous audit not found with identifier tahun 2026 & fakultas F01",
		},
		{
			name:         "InvalidRenstra",
			err:          domain.InvalidRenstra(),
			expectedCode: "GenerateRenstra.InvalidRenstra",
			expectedDesc: "renstra is invalid",
		},
		{
			name:         "InvalidTemplate",
			err:          domain.InvalidTemplate(),
			expectedCode: "GenerateRenstra.InvalidTemplate",
			expectedDesc: "template is invalid",
		},
		{
			name:         "InvalidType",
			err:          domain.InvalidType("unknown"),
			expectedCode: "GenerateRenstra.InvalidType",
			expectedDesc: fmt.Sprintf("type %s is invalid", "unknown"),
		},
		{
			name:         "InvalidTahunRenstra",
			err:          domain.InvalidTahunRenstra("TMP1", "IND1", "2025", "2026", "insert"),
			expectedCode: "GenerateRenstra.InvalidTahunRenstra",
			expectedDesc: "template TMP1 with indicator 'IND1 (2025)' cannot be generated because the indicator year does not match the audit year (2026). Please check the year of the indicator question again, is it correct?",
		},
		{
			name:         "InvalidTahunDokumenTambahan",
			err:          domain.InvalidTahunDokumenTambahan("TMP1", "FileX", "2025", "2026", "insert"),
			expectedCode: "GenerateRenstra.InvalidTahunDokumenTambahan",
			expectedDesc: "template TMP1 with jenis file 'FileX (2025)' cannot be generated because the jenis file year does not match the audit year (2026). Please check the year of the jenis file question again, is it correct?",
		},
		{
			name:         "InvalidTugas",
			err:          domain.InvalidTugas(),
			expectedCode: "GenerateRenstra.InvalidTugas",
			expectedDesc: "tugas is invalid",
		},
		{
			name:         "InvalidParsing",
			err:          domain.InvalidParsing("abc"),
			expectedCode: "GenerateRenstra.InvalidParsing",
			expectedDesc: "failed parsing abc to UUID",
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
