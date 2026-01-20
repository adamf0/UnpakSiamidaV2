package applicationtest

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	common "UnpakSiamida/common/domain"
	repoFakultas "UnpakSiamida/modules/fakultasunit/infrastructure"
	repoIndikator "UnpakSiamida/modules/indikatorrenstra/infrastructure"
	app "UnpakSiamida/modules/previewtemplate/application/GetPreviewTemplate"
	repoPreview "UnpakSiamida/modules/previewtemplate/infrastructure"
)

func newHandler(db *gorm.DB) *app.GetPreviewTemplateByTahunFakultasUnitQueryHandler {
	return &app.GetPreviewTemplateByTahunFakultasUnitQueryHandler{
		Repo:             repoPreview.NewPreviewTemplateRepository(db),
		RepoIndikator:    repoIndikator.NewIndikatorRenstraRepository(db),
		RepoFakultasUnit: repoFakultas.NewFakultasUnitRepository(db),
	}
}

// ======================================================
// SUCCESS CASES
// ======================================================
func Test_GetPreviewTemplateByTahunFakultasUnit_Success(t *testing.T) {
	db, cleanup := setupPreviewTemplateMySQL(t)
	defer cleanup()

	ctx := context.Background()
	handler := newHandler(db)

	tests := []struct {
		name         string
		tahun        string
		fakultasUUID string
		tipe         string
	}{
		{
			name:         "Success Preview Renstra",
			tahun:        "2024",
			fakultasUUID: "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
			tipe:         "renstra",
		},
		{
			name:         "Success Preview DokumenTambahan",
			tahun:        "2024",
			fakultasUUID: "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
			tipe:         "dokumen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := handler.Handle(ctx, app.GetPreviewTemplateByTahunFakultasUnitQuery{
				Tahun:        tt.tahun,
				FakultasUnit: tt.fakultasUUID,
				Tipe:         tt.tipe,
			})

			assert.NoError(t, err)
			assert.NotNil(t, res)
		})
	}
}

// ======================================================
// FAIL CASES
// ======================================================
func Test_GetPreviewTemplateByTahunFakultasUnit_Fail(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupPreviewTemplateMySQL(t)
	defer cleanup()

	handler := newHandler(db)

	tests := []struct {
		name         string
		tahun        string
		FakultasUnit string
		tipe         string
		expectedErr  string
	}{
		{
			name:         "NotFoundFakultasUnit",
			tahun:        "2024",
			FakultasUnit: uuid.NewString(), // TIDAK ADA
			tipe:         "renstra",
			expectedErr:  "PreviewTemplate.NotFoundFakultasUnit",
		},
		{
			name:         "NotFoundTreeIndikator",
			tahun:        "2024",
			FakultasUnit: "e00f35c9-f679-4584-a0c5-8f4ae4e69b0f", // ADA tapi tidak punya tree
			tipe:         "renstra",
			expectedErr:  "PreviewTemplate.NotFoundTreeIndikator",
		},
		{
			name:         "NotFound",
			tahun:        "2080", // TAHUN BELUM ADA DATA
			FakultasUnit: "e00f35c9-f679-4584-a0c5-8f4ae4e69b0f",
			tipe:         "renstra",
			expectedErr:  "PreviewTemplate.NotFound",
		},
		{
			name:         "NotFound",
			tahun:        "2080", // TAHUN BELUM ADA DATA
			FakultasUnit: "fakultas",
			tipe:         "tambahan",
			expectedErr:  "PreviewTemplate.NotFound",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := handler.Handle(ctx, app.GetPreviewTemplateByTahunFakultasUnitQuery{
				Tahun:        tt.tahun,
				FakultasUnit: tt.FakultasUnit,
				Tipe:         tt.tipe,
			})

			assert.Nil(t, res)
			assert.Error(t, err)

			commonErr, ok := err.(common.Error)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedErr, commonErr.Code)
		})
	}
}
