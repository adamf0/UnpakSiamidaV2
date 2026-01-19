package domaintest

import (
	"testing"

	common "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/previewtemplate/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ====================
// PREVIEW TEMPLATE ERROR TESTS
// ====================
func PreviewTemplateErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          common.Error
		expectedCode string
		expectedDesc string
	}{
		{
			name:         "EmptyData_ReturnsCorrectError",
			err:          domain.EmptyData(),
			expectedCode: "PreviewTemplate.EmptyData",
			expectedDesc: "data is not found",
		},
		{
			name:         "NotFoundTreeIndikator_ReturnsCorrectError",
			err:          domain.NotFoundTreeIndikator(),
			expectedCode: "PreviewTemplate.NotFoundTreeIndikator",
			expectedDesc: "Tree Indikator not found",
		},
		{
			name:         "NotFound_ReturnsCorrectError",
			err:          domain.NotFound(),
			expectedCode: "PreviewTemplate.NotFound",
			expectedDesc: "Preview Template not found",
		},
		{
			name:         "NotFoundFakultasUnit_WithDynamicId_ReturnsCorrectError",
			err:          domain.NotFoundFakultasUnit("FK-01"),
			expectedCode: "PreviewTemplate.NotFoundFakultasUnit",
			expectedDesc: "FakultasUnit with identifier FK-01 not found",
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
