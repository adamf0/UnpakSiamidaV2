package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/indikatorrenstra/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func IndikatorRenstraErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "IndikatorRenstra.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "IndikatorRenstra.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidStandar_ReturnsCorrectError",
			err:          domain.InvalidStandar(),
			expectedCode: "IndikatorRenstra.InvalidStandar",
			expectedDesc: "standar renstra is invalid",
		},
		{
			name:         "InvalidParent_ReturnsCorrectError",
			err:          domain.InvalidParent(),
			expectedCode: "IndikatorRenstra.InvalidParent",
			expectedDesc: "parent is invalid",
		},
		{
			name:         "NotUniqueIndikator_ReturnsCorrectError",
			err:          domain.NotUniqueIndikator(),
			expectedCode: "IndikatorRenstra.NotUniqueIndikator",
			expectedDesc: "indikator is not unique",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "IndikatorRenstra.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("ABC123"),
			expectedCode: "IndikatorRenstra.NotFound",
			expectedDesc: "IndikatorRenstra with identifier ABC123 not found",
		},
		{
			name:         "NotFoundParent_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFoundParent("PARENT01"),
			expectedCode: "IndikatorRenstra.InvalidParent",
			expectedDesc: "Parent with identifier PARENT01 not found",
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
