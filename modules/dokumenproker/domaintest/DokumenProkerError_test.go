package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/dokumenproker/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDokumenProkerErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "DokumenProker.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "DokumenProker.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "DokumenProker.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "InvalidFakultas_ReturnsCorrectError",
			err:          domain.InvalidFakultas(),
			expectedCode: "DokumenProker.InvalidFakultas",
			expectedDesc: "fakultas is invalid",
		},
		{
			name:         "InvalidMataProgram_ReturnsCorrectError",
			err:          domain.InvalidMataProgram(),
			expectedCode: "DokumenProker.InvalidMataProgram",
			expectedDesc: "mata program is invalid",
		},
		{
			name:         "InvalidJenisDokumen_ReturnsCorrectError",
			err:          domain.InvalidJenisDokumen(),
			expectedCode: "DokumenProker.InvalidJenisDokumen",
			expectedDesc: "jenis dokumen is invalid",
		},
		{
			name:         "InvalidStatus_ReturnsCorrectError",
			err:          domain.InvalidStatus(),
			expectedCode: "DokumenProker.InvalidStatus",
			expectedDesc: "status is invalid",
		},
		{
			name:         "NotFoundFakultas_ReturnsCorrectError",
			err:          domain.NotFoundFakultas(),
			expectedCode: "DokumenProker.NotFoundFakultas",
			expectedDesc: "fakultas is not found",
		},
		{
			name:         "NotFoundMataProgram_ReturnsCorrectError",
			err:          domain.NotFoundMataProgram(),
			expectedCode: "DokumenProker.NotFoundMataProgram",
			expectedDesc: "mata program is not found",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("XYZ99"),
			expectedCode: "DokumenProker.NotFound",
			expectedDesc: "DokumenProker with identifier XYZ99 not found",
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
