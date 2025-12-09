package domaintest

import (
	"testing"

	"UnpakSiamida/modules/templatedokumentambahan/domain"
	common "UnpakSiamida/common/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateDokumenTambahanError(t *testing.T) {

	tests := []struct {
		name        string
		err         common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData",
			err:          domain.EmptyData(),
			expectedCode: "TemplateDokumenTambahan.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid",
			err:          domain.InvalidUuid(),
			expectedCode: "TemplateDokumenTambahan.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidData",
			err:          domain.InvalidData(),
			expectedCode: "TemplateDokumenTambahan.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFound_with_dynamic_id",
			err:          domain.NotFound("ABC123"),
			expectedCode: "TemplateDokumenTambahan.NotFound",
			expectedDesc: "TemplateDokumenTambahan with identifier ABC123 not found",
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