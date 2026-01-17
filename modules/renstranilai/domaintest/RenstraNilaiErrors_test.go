package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/renstranilai/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =========================
// RENSTRA NILAI ERROR TESTS
// =========================
func TestRenstraNilaiErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "RenstraNilai.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "RenstraNilai.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidRenstra_ReturnsCorrectError",
			err:          domain.InvalidRenstra(),
			expectedCode: "RenstraNilai.InvalidRenstra",
			expectedDesc: "renstra is invalid",
		},
		{
			name:         "RejectAction_ReturnsCorrectError",
			err:          domain.RejectAction(),
			expectedCode: "RenstraNilai.RejectAction",
			expectedDesc: "your action was rejected",
		},
		{
			name:         "NotGranted_ReturnsCorrectError",
			err:          domain.NotGranted(),
			expectedCode: "RenstraNilai.NotGranted",
			expectedDesc: "you are not granted permission in this action",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "RenstraNilai.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("RN-001"),
			expectedCode: "RenstraNilai.NotFound",
			expectedDesc: "RenstraNilai with identifier RN-001 not found",
		},
		{
			name:         "NotFoundRenstra_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFoundRenstra("REN-2026"),
			expectedCode: "RenstraNilai.NotFoundRenstra",
			expectedDesc: "Renstra with identifier REN-2026 not found",
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
