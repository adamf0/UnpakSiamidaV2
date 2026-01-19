package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/templatedokumentambahan/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// TEMPLATE DOKUMEN TAMBAHAN ERROR TESTS
// ========================================
func TemplateDokumenTambahanErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "TemplateDokumenTambahan.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "TemplateDokumenTambahan.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "JenisFileNotFound_ReturnsCorrectError",
			err:          domain.JenisFileNotFound(),
			expectedCode: "TemplateDokumenTambahan.JenisFileNotFound",
			expectedDesc: "fakultas unit not found",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "TemplateDokumenTambahan.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("TD-001"),
			expectedCode: "TemplateDokumenTambahan.NotFound",
			expectedDesc: "TemplateDokumenTambahan with identifier TD-001 not found",
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
