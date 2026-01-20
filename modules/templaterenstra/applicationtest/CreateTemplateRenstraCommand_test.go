package applicationtest

import (
	"context"
	"testing"

	common "UnpakSiamida/common/domain"
	infrafakultas "UnpakSiamida/modules/fakultasunit/infrastructure"
	infraindikator "UnpakSiamida/modules/indikatorrenstra/infrastructure"
	app "UnpakSiamida/modules/templaterenstra/application/CreateTemplateRenstra"
	infra "UnpakSiamida/modules/templaterenstra/infrastructure"

	"github.com/goforj/godump"
	"github.com/stretchr/testify/assert"
)

func TestCreateTemplateRenstraCommandHandler_Success(t *testing.T) {
	db, cleanup := setupTemplateRenstraMySQL(t)
	defer cleanup()

	templateRepo := infra.NewTemplateRenstraRepository(db)
	indikatorRepo := infraindikator.NewIndikatorRenstraRepository(db)
	fakultasRepo := infrafakultas.NewFakultasUnitRepository(db)

	handler := app.CreateTemplateRenstraCommandHandler{
		Repo:                 templateRepo,
		IndikatorRenstraRepo: indikatorRepo,
		FakultasUnitRepo:     fakultasRepo,
	}

	satuan := "% Lulusan"
	target := "15"

	cmd := app.CreateTemplateRenstraCommand{
		Tahun:        "2024",
		Indikator:    "b763b5b3-a18e-416c-9d0d-a0c23aa6076c",
		IsPertanyaan: "1",
		FakultasUnit: "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Kategori:     "fakultas#all",
		Klasifikasi:  "minor",
		Satuan:       &satuan,
		Target:       &target,
		TargetMin:    nil,
		TargetMax:    nil,
		Tugas:        "auditor1",
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

func TestCreateTemplateRenstraCommandValidation_Success(t *testing.T) {
	satuan := "% Lulusan"
	target := "15"

	cmd := app.CreateTemplateRenstraCommand{
		Tahun:        "2024",
		Indikator:    "b763b5b3-a18e-416c-9d0d-a0c23aa6076c",
		IsPertanyaan: "1",
		FakultasUnit: "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Kategori:     "fakultas#all",
		Klasifikasi:  "minor",
		Satuan:       &satuan,
		Target:       &target,
		TargetMin:    nil,
		TargetMax:    nil,
		Tugas:        "auditor1",
	}

	err := app.CreateTemplateRenstraCommandValidation(cmd)
	assert.NoError(t, err)
}

func TestCreateTemplateRenstraCommandValidation_Fail(t *testing.T) {
	cmd := app.CreateTemplateRenstraCommand{
		Tahun:        "",
		Indikator:    "",
		IsPertanyaan: "",
		FakultasUnit: "",
		Kategori:     "",
		Klasifikasi:  "",
		Tugas:        "",
	}

	err := app.CreateTemplateRenstraCommandValidation(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "Indikator cannot be blank")
	assert.Contains(t, err.Error(), "IsPertanyaan cannot be blank")
	assert.Contains(t, err.Error(), "FakultasUnit cannot be blank")
	assert.Contains(t, err.Error(), "Kategori cannot be blank")
	assert.Contains(t, err.Error(), "Klasifikasi cannot be blank")
	assert.Contains(t, err.Error(), "Tugas cannot be blank")
}

func TestCreateTemplateRenstraCommandHandler_Fail(t *testing.T) {
	tests := []struct {
		name          string
		indikatorUUID string
		fakultasUUID  string
		target        *string
		targetMin     *string
		targetMax     *string
		expectedCode  string
	}{
		{
			name:          "indikator not found",
			indikatorUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", // tidak ada
			fakultasUUID:  "0d2fa3f8-6df3-45b8-8985-654cb49d5d03", // valid
			target:        nil,
			targetMin:     nil,
			targetMax:     nil,
			expectedCode:  "TemplateRenstra.IndikatorNotFound",
		},
		{
			name:          "fakultas unit not found",
			indikatorUUID: "b763b5b3-a18e-416c-9d0d-a0c23aa6076c", // valid
			fakultasUUID:  "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", // tidak ada
			target:        nil,
			targetMin:     nil,
			targetMax:     nil,
			expectedCode:  "TemplateRenstra.FakultasUnitNotFound",
		},
		{
			name:          "InvalidValueTarget",
			indikatorUUID: "b763b5b3-a18e-416c-9d0d-a0c23aa6076c", // valid
			fakultasUUID:  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0", // tidak ada
			target:        nil,
			targetMin:     strPtr("80"),
			targetMax:     strPtr("120"),
			expectedCode:  "TemplateRenstra.InvalidValueTarget",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := setupTemplateRenstraMySQL(t)
			defer cleanup()

			handler := app.CreateTemplateRenstraCommandHandler{
				Repo:                 infra.NewTemplateRenstraRepository(db),
				IndikatorRenstraRepo: infraindikator.NewIndikatorRenstraRepository(db),
				FakultasUnitRepo:     infrafakultas.NewFakultasUnitRepository(db),
			}

			cmd := app.CreateTemplateRenstraCommand{
				Tahun:        "2025",
				Indikator:    tt.indikatorUUID,
				IsPertanyaan: "1",
				FakultasUnit: tt.fakultasUUID,
				Kategori:     "akademik",
				Klasifikasi:  "minor",
				Tugas:        "auditor1",
				Target:       tt.target,
				TargetMin:    tt.targetMin,
				TargetMax:    tt.targetMax,
			}

			_, err := handler.Handle(context.Background(), cmd)

			assert.Error(t, err)

			domainErr, ok := err.(common.Error)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, domainErr.Code)
		})
	}
}

func TestCreateTemplateRenstra_ContextTimeout(t *testing.T) {
	db, cleanup := setupTemplateRenstraMySQL(t)
	defer cleanup()

	templateRepo := infra.NewTemplateRenstraRepository(db)
	indikatorRepo := infraindikator.NewIndikatorRenstraRepository(db)
	fakultasRepo := infrafakultas.NewFakultasUnitRepository(db)

	handler := app.CreateTemplateRenstraCommandHandler{
		Repo:                 templateRepo,
		IndikatorRenstraRepo: indikatorRepo,
		FakultasUnitRepo:     fakultasRepo,
	}

	satuan := "% Lulusan"
	target := "15"

	cmd := app.CreateTemplateRenstraCommand{
		Tahun:        "2024",
		Indikator:    "b763b5b3-a18e-416c-9d0d-a0c23aa6076c",
		IsPertanyaan: "1",
		FakultasUnit: "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Kategori:     "fakultas#all",
		Klasifikasi:  "minor",
		Satuan:       &satuan,
		Target:       &target,
		TargetMin:    nil,
		TargetMax:    nil,
		Tugas:        "auditor1",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := handler.Handle(ctx, cmd)
	assert.NoError(t, err)
	assert.True(t, err == context.Canceled || err == context.DeadlineExceeded, "expected context canceled or timeout error")
}

func TestCreateTemplateRenstraCommandValidation_Duplicate(t *testing.T) {
	satuan := "% Lulusan"
	target := "15"

	cmd := app.CreateTemplateRenstraCommand{
		Tahun:        "2024",
		Indikator:    "b763b5b3-a18e-416c-9d0d-a0c23aa6076c",
		IsPertanyaan: "1",
		FakultasUnit: "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Kategori:     "fakultas#all",
		Klasifikasi:  "minor",
		Satuan:       &satuan,
		Target:       &target,
		TargetMin:    nil,
		TargetMax:    nil,
		Tugas:        "auditor1",
	}

	err := app.CreateTemplateRenstraCommandValidation(cmd)
	assert.NoError(t, err)

	//

	cmd = app.CreateTemplateRenstraCommand{
		Tahun:        "2024",
		Indikator:    "b763b5b3-a18e-416c-9d0d-a0c23aa6076c",
		IsPertanyaan: "1",
		FakultasUnit: "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Kategori:     "fakultas#all",
		Klasifikasi:  "minor",
		Satuan:       &satuan,
		Target:       &target,
		TargetMin:    nil,
		TargetMax:    nil,
		Tugas:        "auditor1",
	}

	err = app.CreateTemplateRenstraCommandValidation(cmd)
	godump.Dump(err)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Error(t, err)
	assert.Equal(t, "TemplateRenstra.DuplicateData", commonErr.Code)
	assert.Equal(t, "data not allowed duplicate", commonErr.Description)
}

func strPtr(v string) *string {
	return &v
}
