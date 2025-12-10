package domaintest

import (
	"testing"

	"UnpakSiamida/modules/tahunrenstra/domain"
	common "UnpakSiamida/common/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTahunRenstraErrors(t *testing.T) {

	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "TahunRenstra.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "NotFound_ReturnsCorrectError",
			err:          domain.NotFound("XYZ99"),
			expectedCode: "TahunRenstra.NotFound",
			expectedDesc: "TahunRenstra with identifier XYZ99 not found",
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
