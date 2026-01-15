package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/dokumentambahan/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDokumenTambahanErrors(t *testing.T) {

	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "DokumenTambahan.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "InvalidUuid_ReturnsCorrectError",
			err:          domain.InvalidUuid(),
			expectedCode: "DokumenTambahan.InvalidUuid",
			expectedDesc: "uuid is invalid",
		},
		{
			name:         "InvalidRenstra_ReturnsCorrectError",
			err:          domain.InvalidRenstra(),
			expectedCode: "DokumenTambahan.InvalidRenstra",
			expectedDesc: "renstra is invalid",
		},
		{
			name:         "RejectAction_ReturnsCorrectError",
			err:          domain.RejectAction(),
			expectedCode: "DokumenTambahan.RejectAction",
			expectedDesc: "your action was rejected",
		},
		{
			name:         "NotGranted_ReturnsCorrectError",
			err:          domain.NotGranted(),
			expectedCode: "DokumenTambahan.NotGranted",
			expectedDesc: "you are not granted permission in this action",
		},
		{
			name:         "InvalidData_ReturnsCorrectError",
			err:          domain.InvalidData(),
			expectedCode: "DokumenTambahan.InvalidData",
			expectedDesc: "data is invalid",
		},
		{
			name:         "NotFound_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFound("DOC123"),
			expectedCode: "DokumenTambahan.NotFound",
			expectedDesc: "DokumenTambahan with identifier DOC123 not found",
		},
		{
			name:         "NotFoundRenstra_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFoundRenstra("REN456"),
			expectedCode: "DokumenTambahan.NotFoundRenstra",
			expectedDesc: "Renstra with identifier REN456 not found",
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
