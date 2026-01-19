package applicationtest

import (
	"context"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/jenisfile/application/DeleteJenisFile"
	domain "UnpakSiamida/modules/jenisfile/domain"
	infra "UnpakSiamida/modules/jenisfile/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteJenisFileCommandValidation_Success(t *testing.T) {
	validCmd := app.DeleteJenisFileCommand{
		Uuid: uuid.NewString(),
	}
	err := app.DeleteJenisFileCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestDeleteJenisFileCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.DeleteJenisFileCommand{
		Uuid: "",
	}
	err := app.DeleteJenisFileCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

func TestDeleteJenisFileCommand_Success(t *testing.T) {
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

func TestDeleteJenisFileCommand_Edge(t *testing.T) {
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
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDeleteJenisFileCommand_Fail(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.DeleteJenisFileCommandHandler{Repo: repo}

	// UUID tidak valid
	cmdInvalidUUID := app.DeleteJenisFileCommand{
		Uuid: uuid.NewString(),
	}
	_, err := handler.Handle(context.Background(), cmdInvalidUUID)
	assert.Error(t, err)

	commonErr, _ := err.(common.Error)

	assert.Equal(t, "JenisFile.EmptyData", commonErr.Code)
	assert.Equal(t, "data is not found", commonErr.Description)
}
