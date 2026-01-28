package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/jadwalproker/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJadwalProkerErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "JadwalProker.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "JadwalProker.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "JadwalProker.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "InvalidFakultas_ReturnsCorrectError",
			err:          domain.InvalidFakultas(),
			expectedCode: "JadwalProker.InvalidFakultas",
			expectedDesc: "fakultas is invalid",
		},
		{
			name:         "NotFoundFakultas_ReturnsCorrectError",
			err:          domain.NotFoundFakultas(),
			expectedCode: "JadwalProker.NotFoundFakultas",
			expectedDesc: "fakultas is not found",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("XYZ99"),
			expectedCode: "JadwalProker.NotFound",
			expectedDesc: "JadwalProker with identifier XYZ99 not found",
		},
		{
			name:         "InvalidDate_WithTarget_ReturnsCorrectError",
			err:          domain.InvalidDate("Periode"),
			expectedCode: "JadwalProker.InvalidDate",
			expectedDesc: "Periode period have wrong date format",
		},
		{
			name:         "InvalidDateRange_ReturnsCorrectError",
			err:          domain.InvalidDateRange(),
			expectedCode: "JadwalProker.InvalidDateRange",
			expectedDesc: "tanggal upload dokumen must not be earlier than tanggal input",
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
