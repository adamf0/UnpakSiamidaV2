package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/aktivitasproker/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAktivitasErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "Aktivitas.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "Aktivitas.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "Aktivitas.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "InvalidFakultas_ReturnsCorrectError",
			err:          domain.InvalidFakultas(),
			expectedCode: "Aktivitas.InvalidFakultas",
			expectedDesc: "fakultas is invalid",
		},
		{
			name:         "InvalidMataProgram_ReturnsCorrectError",
			err:          domain.InvalidMataProgram(),
			expectedCode: "Aktivitas.InvalidMataProgram",
			expectedDesc: "mata program is invalid",
		},
		{
			name:         "NotFoundFakultas_ReturnsCorrectError",
			err:          domain.NotFoundFakultas(),
			expectedCode: "Aktivitas.NotFoundFakultas",
			expectedDesc: "fakultas is not found",
		},
		{
			name:         "NotFoundMataProgram_ReturnsCorrectError",
			err:          domain.NotFoundMataProgram(),
			expectedCode: "Aktivitas.NotFoundMataProgram",
			expectedDesc: "mata program is not found",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("XYZ99"),
			expectedCode: "Aktivitas.NotFound",
			expectedDesc: "Aktivitas with identifier XYZ99 not found",
		},
		{
			name:         "InvalidDate_WithTarget_ReturnsCorrectError",
			err:          domain.InvalidDate("Periode"),
			expectedCode: "Aktivitas.InvalidDate",
			expectedDesc: "Periode period have wrong date format",
		},
		{
			name:         "InvalidDateRange_ReturnsCorrectError",
			err:          domain.InvalidDateRange(),
			expectedCode: "Aktivitas.InvalidDateRange",
			expectedDesc: "tanggal rk akhir must not be earlier than tanggal rk awal",
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
