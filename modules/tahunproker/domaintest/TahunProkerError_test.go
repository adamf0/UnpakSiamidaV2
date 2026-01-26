package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/tahunproker/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTahunProkerErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "TahunProker.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "TahunProker.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "TahunProker.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "InvalidTahun_ReturnsCorrectError",
			err:          domain.InvalidTahun(),
			expectedCode: "TahunProker.InvalidTahun",
			expectedDesc: "tahun is invalid",
		},
		{
			name:         "InvalidStatus_ReturnsCorrectError",
			err:          domain.InvalidStatus(),
			expectedCode: "TahunProker.InvalidStatus",
			expectedDesc: "status is invalid",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("XYZ99"),
			expectedCode: "TahunProker.NotFound",
			expectedDesc: "TahunProker with identifier XYZ99 not found",
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
