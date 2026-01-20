package applicationtest

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

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
			tipe:         "tag",
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
			assert.NotEmpty(t, res)

			for _, r := range res {
				assert.NotEmpty(t, r.Pointing)
			}
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
		fakultasUUID string
		tipe         string
		expectedErr  string
	}{
		{
			name:         "NotFoundFakultasUnit",
			tahun:        "2024",
			fakultasUUID: uuid.NewString(), // TIDAK ADA
			tipe:         "renstra",
			expectedErr:  "PreviewTemplate.NotFoundFakultasUnit",
		},
		{
			name:         "NotFoundTreeIndikator",
			tahun:        "2024",
			fakultasUUID: "11111111-1111-1111-1111-111111111111", // ADA tapi tidak punya tree
			tipe:         "renstra",
			expectedErr:  "PreviewTemplate.NotFoundTreeIndikator",
		},
		{
			name:         "NotFoundTemplate",
			tahun:        "2080", // TAHUN BELUM ADA DATA
			fakultasUUID: "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
			tipe:         "renstra",
			expectedErr:  "PreviewTemplate.NotFoundTemplate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := handler.Handle(ctx, app.GetPreviewTemplateByTahunFakultasUnitQuery{
				Tahun:        tt.tahun,
				FakultasUnit: tt.fakultasUUID,
				Tipe:         tt.tipe,
			})

			assert.Nil(t, res)
			assert.Error(t, err)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
