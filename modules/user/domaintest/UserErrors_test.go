package domaintest

import (
	"testing"

	domain "UnpakSiamida/modules/user/domain"
	common "UnpakSiamida/common/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserErrors_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		err           common.Error
		expectedCode  string
		expectedDesc  string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "User.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "User.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "User.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "InvalidEmail_ReturnsCorrectError",
			err:          domain.InvalidEmail(),
			expectedCode: "User.InvalidEmail",
			expectedDesc: "email tidak valid atau tidak diperbolehkan",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("12345"),
			expectedCode: "User.NotFound",
			expectedDesc: "User with identifier 12345 not found",
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
