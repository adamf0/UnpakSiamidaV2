package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/account/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountErrors(t *testing.T) {

	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "InvalidCredential_ReturnsCorrectError",
			err:          domain.InvalidCredential(),
			expectedCode: "Account.InvalidCredential",
			expectedDesc: "invalid credentials",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("ABC123"),
			expectedCode: "Account.NotFound",
			expectedDesc: "Account with identifier ABC123 not found",
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
