package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/generaterenstra/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRenstraErrors_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "GenerateRenstra.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "GenerateRenstra.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "GenerateRenstra.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("12345"),
			expectedCode: "GenerateRenstra.NotFound",
			expectedDesc: "GenerateRenstra with identifier 12345 not found",
		},
		{
			name:         "NotFoundFakultasUnit_ReturnsCorrectError",
			err:          domain.NotFoundFakultasUnit("F01"),
			expectedCode: "GenerateRenstra.NotFoundFakultasUnit",
			expectedDesc: "Fakultas unit with identifier F01 not found",
		},
		{
			name:         "InvalidRenstra_ReturnsCorrectError",
			err:          domain.InvalidRenstra(),
			expectedCode: "GenerateRenstra.InvalidRenstra",
			expectedDesc: "renstra is invalid",
		},
		{
			name:         "InvalidTemplate_ReturnsCorrectError",
			err:          domain.InvalidTemplate(),
			expectedCode: "GenerateRenstra.InvalidTemplate",
			expectedDesc: "template is invalid",
		},
		{
			name:         "InvalidParsing_ReturnsCorrectError",
			err:          domain.InvalidParsing("abc"),
			expectedCode: "GenerateRenstra.InvalidParsing",
			expectedDesc: "failed parsing abc to UUID",
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
