package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/user/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
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
			name:         "InvalidFakultasUnit_ReturnsCorrectError",
			err:          domain.InvalidFakultasUnit(),
			expectedCode: "User.InvalidFakultasUnit",
			expectedDesc: "fakultas unit tidak valid",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("USR-001"),
			expectedCode: "User.NotFound",
			expectedDesc: "user with identifier USR-001 not found",
		},
		{
			name:         "InvalidParsing_WithTarget_ReturnsCorrectError",
			err:          domain.InvalidParsing("user_id"),
			expectedCode: "User.InvalidParsing",
			expectedDesc: "failed parsing user_id to UUID",
		},
		{
			name:         "NotFoundFakultasUnit_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFoundFakultasUnit("FK-01"),
			expectedCode: "User.NotFoundFakultasUnit",
			expectedDesc: "fakultas unit with identifier FK-01 not found",
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
