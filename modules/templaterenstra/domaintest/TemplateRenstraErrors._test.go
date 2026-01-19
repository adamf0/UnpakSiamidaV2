package domaintest

import (
	"testing"

	"UnpakSiamida/modules/templaterenstra/domain"
	common "UnpakSiamida/common/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateRenstraError(t *testing.T) {

	tests := []struct {
		name        string
		err         common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData",
			err:          domain.EmptyData(),
			expectedCode: "TemplateRenstra.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid",
			err:          domain.InvalidUuid(),
			expectedCode: "TemplateRenstra.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "IndikatorNotFound",
			err:          domain.IndikatorNotFound(),
			expectedCode: "TemplateRenstra.IndikatorNotFound",
			expectedDesc: "indikator not found",
		},
		{
			name:         "FakultasUnitNotFound",
			err:          domain.FakultasUnitNotFound(),
			expectedCode: "TemplateRenstra.FakultasUnitNotFound",
			expectedDesc: "fakultas unit not found",
		},
		{
			name:         "InvalidValueTarget",
			err:          domain.InvalidValueTarget(),
			expectedCode: "TemplateRenstra.InvalidValueTarget",
			expectedDesc: "invalid target combination, either provide Target only or provide both TargetMin and TargetMax",
		},
		{
			name:         "InvalidData",
			err:          domain.InvalidData(),
			expectedCode: "TemplateRenstra.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFound_with_dynamic_id",
			err:          domain.NotFound("ABC123"),
			expectedCode: "TemplateRenstra.NotFound",
			expectedDesc: "TemplateRenstra with identifier ABC123 not found",
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