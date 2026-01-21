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
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("XYZ99"),
			expectedCode: "BeritaAcara.NotFound",
			expectedDesc: "BeritaAcara with identifier XYZ99 not found",
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
