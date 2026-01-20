package applicationtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	common "UnpakSiamida/common/domain"
	infraJenisFile "UnpakSiamida/modules/jenisfile/infrastructure"
	app "UnpakSiamida/modules/templatedokumentambahan/application/UpdateTemplateDokumenTambahan"
	infra "UnpakSiamida/modules/templatedokumentambahan/infrastructure"
)

// -----------------------
// SUCCESS VALIDATION
// -----------------------
func TestUpdateTemplateDokumenTambahanCommandValidation_Success(t *testing.T) {
	cmd := app.UpdateTemplateDokumenTambahanCommand{
		Uuid:        "9b354f31-be71-4173-9e26-c319d163660d",
		Tahun:       "2031",
		JenisFile:   "14212231-792f-4935-bb1c-9a38695a4b6b",
		Pertanyaan:  "Apa yang harus dilakukan?",
		Klasifikasi: "minor",
		Kategori:    "fakultas#all",
		Tugas:       "auditor1",
	}

	err := app.UpdateTemplateDokumenTambahanCommandValidation(cmd)
	assert.NoError(t, err)
}

// -----------------------
// FAIL VALIDATION
// -----------------------
func TestUpdateTemplateDokumenTambahanCommandValidation_Fail(t *testing.T) {
	cmd := app.UpdateTemplateDokumenTambahanCommand{
		Uuid:        "",
		Tahun:       "",
		JenisFile:   "",
		Pertanyaan:  "",
		Klasifikasi: "",
		Kategori:    "",
		Tugas:       "",
	}

	err := app.UpdateTemplateDokumenTambahanCommandValidation(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "JenisFile cannot be blank")
	assert.Contains(t, err.Error(), "Pertanyaan cannot be blank")
	assert.Contains(t, err.Error(), "Klasifikasi cannot be blank")
	assert.Contains(t, err.Error(), "Kategori cannot be blank")
	assert.Contains(t, err.Error(), "Tugas cannot be blank")
}

// -----------------------
// SUCCESS HANDLER
// -----------------------
func TestUpdateTemplateDokumenTambahanCommandHandler_Success(t *testing.T) {
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	jenisFileRepo := infraJenisFile.NewJenisFileRepository(db)
	repo := infra.NewTemplateDokumenTambahanRepository(db)

	handler := &app.UpdateTemplateDokumenTambahanCommandHandler{
		Repo:          repo,
		JenisFileRepo: jenisFileRepo,
	}

	cmd := app.UpdateTemplateDokumenTambahanCommand{
		Uuid:        "9b354f31-be71-4173-9e26-c319d163660d",
		Tahun:       "2031",
		JenisFile:   "14212231-792f-4935-bb1c-9a38695a4b6b",
		Pertanyaan:  "Apa yang harus dilakukan?",
		Klasifikasi: "minor",
		Kategori:    "fakultas#all",
		Tugas:       "auditor1",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

// -----------------------
// FAIL HANDLER
// -----------------------
func TestUpdateTemplateDokumenTambahanCommandHandler_Fail(t *testing.T) {
	// Setup DB & repos
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	jenisFileRepo := infraJenisFile.NewJenisFileRepository(db)
	repo := infra.NewTemplateDokumenTambahanRepository(db)

	handler := &app.UpdateTemplateDokumenTambahanCommandHandler{
		Repo:          repo,
		JenisFileRepo: jenisFileRepo,
	}

	invalidUuid := uuid.NewString()
	invalidUuidJenisFile := uuid.NewString()

	tests := []struct {
		name          string
		cmd           app.UpdateTemplateDokumenTambahanCommand
		expectErrCode string
		expectErrText string
	}{
		{
			name: "JenisFileNotFound",
			cmd: app.UpdateTemplateDokumenTambahanCommand{
				Uuid:        "9b354f31-be71-4173-9e26-c319d163660d",
				Tahun:       "2031",
				JenisFile:   invalidUuidJenisFile, // random UUID → not found
				Pertanyaan:  "Apa yang harus dilakukan?",
				Klasifikasi: "minor",
				Kategori:    "fakultas#all",
				Tugas:       "auditor1",
			},
			expectErrCode: "TemplateDokumenTambahan.JenisFileNotFound",
			expectErrText: "JenisFile is not found",
		},
		{
			name: "NotFound",
			cmd: app.UpdateTemplateDokumenTambahanCommand{
				Uuid:        invalidUuid, // random template UUID → not found
				Tahun:       "2031",
				JenisFile:   "14212231-792f-4935-bb1c-9a38695a4b6b",
				Pertanyaan:  "Apa yang harus dilakukan?",
				Klasifikasi: "minor",
				Kategori:    "fakultas#all",
				Tugas:       "auditor1",
			},
			expectErrCode: "TemplateDokumenTambahan.NotFound",
			expectErrText: fmt.Sprintf("TemplateDokumenTambahan with identifier %s not found", invalidUuid),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Handle(context.Background(), tt.cmd)
			assert.Error(t, err)

			commonErr, ok := err.(common.Error)
			assert.True(t, ok, "error should be of type common.Error")
			assert.Equal(t, tt.expectErrCode, commonErr.Code)
			assert.Contains(t, commonErr.Description, tt.expectErrText)
		})
	}
}

func TestUpdateTemplateDokumenTambahanCommandHandler_Duplicate(t *testing.T) {
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	jenisFileRepo := infraJenisFile.NewJenisFileRepository(db)
	repo := infra.NewTemplateDokumenTambahanRepository(db)

	handler := &app.UpdateTemplateDokumenTambahanCommandHandler{
		Repo:          repo,
		JenisFileRepo: jenisFileRepo,
	}

	cmd := app.UpdateTemplateDokumenTambahanCommand{
		Uuid:        "9b354f31-be71-4173-9e26-c319d163660d",
		Tahun:       "2031",
		JenisFile:   "14212231-792f-4935-bb1c-9a38695a4b6b",
		Pertanyaan:  "Apa yang harus dilakukan?",
		Klasifikasi: "minor",
		Kategori:    "fakultas#all",
		Tugas:       "auditor1",
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)

	//

	cmd = app.UpdateTemplateDokumenTambahanCommand{
		Uuid:        "9cc38a80-d11b-44c9-9b81-f394816eaa96",
		Tahun:       "2031",
		JenisFile:   "14212231-792f-4935-bb1c-9a38695a4b6b",
		Pertanyaan:  "Apa yang harus dilakukan?",
		Klasifikasi: "minor",
		Kategori:    "fakultas#all",
		Tugas:       "auditor1",
	}

	res, err = handler.Handle(context.Background(), cmd)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Error(t, err)
	assert.Equal(t, "TemplateDokumenTambahan.DuplicateData", commonErr.Code)
	assert.Contains(t, "data not allowed duplicate", commonErr.Description)
}
