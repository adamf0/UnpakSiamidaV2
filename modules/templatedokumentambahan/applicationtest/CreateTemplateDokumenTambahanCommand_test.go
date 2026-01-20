package applicationtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	common "UnpakSiamida/common/domain"
	infraJenisFile "UnpakSiamida/modules/jenisfile/infrastructure"
	app "UnpakSiamida/modules/templatedokumentambahan/application/CreateTemplateDokumenTambahan"
	infra "UnpakSiamida/modules/templatedokumentambahan/infrastructure"
)

// -----------------------
// SUCCESS VALIDATION
// -----------------------
func TestCreateTemplateDokumenTambahanCommandValidation_Success(t *testing.T) {
	cmd := app.CreateTemplateDokumenTambahanCommand{
		Tahun:       "2031",
		JenisFile:   "14212231-792f-4935-bb1c-9a38695a4b6b",
		Pertanyaan:  "Apa yang harus dilakukan?",
		Klasifikasi: "minor",
		Kategori:    "fakultas#all",
		Tugas:       "auditor1",
	}

	err := app.CreateTemplateDokumenTambahanCommandValidation(cmd)
	assert.NoError(t, err)
}

// -----------------------
// FAIL VALIDATION
// -----------------------
func TestCreateTemplateDokumenTambahanCommandValidation_Fail(t *testing.T) {
	cmd := app.CreateTemplateDokumenTambahanCommand{
		Tahun:       "",
		JenisFile:   "",
		Pertanyaan:  "",
		Klasifikasi: "",
		Kategori:    "",
		Tugas:       "",
	}

	err := app.CreateTemplateDokumenTambahanCommandValidation(cmd)
	assert.Error(t, err)
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
func TestCreateTemplateDokumenTambahanCommandHandler_Success(t *testing.T) {
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	jenisFileRepo := infraJenisFile.NewJenisFileRepository(db)
	repo := infra.NewTemplateDokumenTambahanRepository(db)

	handler := &app.CreateTemplateDokumenTambahanCommandHandler{
		Repo:          repo,
		JenisFileRepo: jenisFileRepo,
	}

	cmd := app.CreateTemplateDokumenTambahanCommand{
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
// FAIL HANDLER (JenisFileNotFound)
// -----------------------
func TestCreateTemplateDokumenTambahanCommandHandler_FailJenisFileNotFound(t *testing.T) {
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	jenisFileRepo := infraJenisFile.NewJenisFileRepository(db)
	repo := infra.NewTemplateDokumenTambahanRepository(db)

	handler := &app.CreateTemplateDokumenTambahanCommandHandler{
		Repo:          repo,
		JenisFileRepo: jenisFileRepo,
	}

	uuid := uuid.NewString()

	cmd := app.CreateTemplateDokumenTambahanCommand{
		Tahun:       "2031",
		JenisFile:   uuid,
		Pertanyaan:  "Apa yang harus dilakukan?",
		Klasifikasi: "minor",
		Kategori:    "fakultas#all",
		Tugas:       "auditor1",
	}

	_, err := handler.Handle(context.Background(), cmd)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Error(t, err)
	assert.Equal(t, "TemplateDokumenTambahan.JenisFileNotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("JenisFile with identifier %s not found", uuid), commonErr.Description)
}

//[note] mustihal masuk duplicate karena simpannya saja sudah pasang onconflic to update
// func TestCreateTemplateDokumenTambahanCommandHandler_Duplicate(t *testing.T) {
// 	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
// 	defer cleanup()

// 	jenisFileRepo := infraJenisFile.NewJenisFileRepository(db)
// 	repo := infra.NewTemplateDokumenTambahanRepository(db)

// 	handler := &app.CreateTemplateDokumenTambahanCommandHandler{
// 		Repo:          repo,
// 		JenisFileRepo: jenisFileRepo,
// 	}

// 	cmd := app.CreateTemplateDokumenTambahanCommand{
// 		Tahun:       "2031",
// 		JenisFile:   "14212231-792f-4935-bb1c-9a38695a4b6b",
// 		Pertanyaan:  "Apa yang harus dilakukan?",
// 		Klasifikasi: "minor",
// 		Kategori:    "fakultas#all",
// 		Tugas:       "auditor1",
// 	}

// 	res, err := handler.Handle(context.Background(), cmd)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, res)

// 	//

// 	res, err = handler.Handle(context.Background(), cmd)
// 	godump.Dump(res, err)
// 	assert.Error(t, err)

// 	commonErr, ok := err.(common.Error)
// 	assert.True(t, ok)
// 	assert.Error(t, err)
// 	assert.Equal(t, "TemplateDokumenTambahan.DuplicateData", commonErr.Code)
// 	assert.Equal(t, "data not allowed duplicate", commonErr.Description)
// }
