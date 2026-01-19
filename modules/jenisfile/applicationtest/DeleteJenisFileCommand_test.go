package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/jenisfile/application/DeleteJenisFile"
	domain "UnpakSiamida/modules/jenisfile/domain"
	infra "UnpakSiamida/modules/jenisfile/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func DeleteJenisFileCommandValidation_Success(t *testing.T) {
	validCmd := app.DeleteJenisFileCommand{
		Uuid: uuid.NewString(),
	}
	err := app.DeleteJenisFileCommandValidation(validCmd)
	assert.NoError(t, err)
}

func DeleteJenisFileCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.DeleteJenisFileCommand{
		Uuid: "",
	}
	err := app.DeleteJenisFileCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

func DeleteJenisFileCommand_Success(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.DeleteJenisFileCommandHandler{Repo: repo}

	// Insert record dulu
	record := domain.JenisFile{
		UUID: uuid.New(),
		Nama: "Dokumen Untuk Delete",
	}
	err := repo.Create(context.Background(), &record)
	assert.NoError(t, err)

	cmd := app.DeleteJenisFileCommand{
		Uuid: record.UUID.String(),
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, record.UUID.String(), deletedUUID)

	// Pastikan DB sudah terhapus
	var saved domain.JenisFile
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err) // harus error karena sudah dihapus
}

func DeleteJenisFileCommand_Edge(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.DeleteJenisFileCommandHandler{Repo: repo}

	// Insert record
	record := domain.JenisFile{
		UUID: uuid.New(),
		Nama: "Dokumen Edge Delete",
	}
	err := repo.Create(context.Background(), &record)
	assert.NoError(t, err)

	cmd := app.DeleteJenisFileCommand{
		Uuid: record.UUID.String(),
	}

	// Delete pertama → sukses
	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, record.UUID.String(), deletedUUID)

	// Delete kedua → harus not found
	_, err = handler.Handle(context.Background(), cmd)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Error(t, err)
	assert.Equal(t, "JenisFile.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("JenisFile with identifier %s not found", record.UUID.String()), commonErr.Description)
}

func DeleteJenisFileCommand_Fail(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.DeleteJenisFileCommandHandler{Repo: repo}

	uuid := uuid.NewString()

	// UUID tidak valid
	cmdInvalidUUID := app.DeleteJenisFileCommand{
		Uuid: uuid,
	}
	_, err := handler.Handle(context.Background(), cmdInvalidUUID)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "JenisFile.NotFound", commonErr.Code)
	assert.Contains(t, fmt.Sprintf("JenisFile with identifier %s not found", uuid), commonErr.Description)
}
