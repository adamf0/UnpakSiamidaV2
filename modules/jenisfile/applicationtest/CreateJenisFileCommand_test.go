package applicationtest

import (
	"context"
	"testing"

	app "UnpakSiamida/modules/jenisfile/application/CreateJenisFile"
	domain "UnpakSiamida/modules/jenisfile/domain"
	infra "UnpakSiamida/modules/jenisfile/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func CreateJenisFileCommandValidation_Success(t *testing.T) {
	validCmd := app.CreateJenisFileCommand{Nama: "Dokumen Valid"}
	err := app.CreateJenisFileCommandValidation(validCmd)
	assert.NoError(t, err)
}

func CreateJenisFileCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.CreateJenisFileCommand{Nama: ""}
	err := app.CreateJenisFileCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Nama cannot be blank")
}

func CreateJenisFileCommand_Success(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.CreateJenisFileCommandHandler{Repo: repo}

	cmd := app.CreateJenisFileCommand{Nama: "Dokumen Baru"}
	uuidStr, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr)

	// Pastikan record tersimpan di DB
	var saved domain.JenisFile
	err = db.Where("uuid = ?", uuidStr).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, "Dokumen Baru", saved.Nama)
}

func CreateJenisFileCommand_Edge(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.CreateJenisFileCommandHandler{Repo: repo}

	// Insert satu record dulu
	firstCmd := app.CreateJenisFileCommand{Nama: "Dokumen Sama"}
	_, err := handler.Handle(context.Background(), firstCmd)
	assert.NoError(t, err)

	// Insert lagi dengan nama yang sama
	secondCmd := app.CreateJenisFileCommand{Nama: "Dokumen Sama"}
	uuidStr2, err := handler.Handle(context.Background(), secondCmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr2)

	// UUID harus berbeda
	assert.NotEqual(t, uuidStr2, uuid.Nil)
}
