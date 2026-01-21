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

	cmd := app.DeleteJenisFileCommand{
		Uuid: "14212231-792f-4935-bb1c-9a38695a4b6b",
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "14212231-792f-4935-bb1c-9a38695a4b6b", deletedUUID)

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

	cmd := app.DeleteJenisFileCommand{
		Uuid: "14212231-792f-4935-bb1c-9a38695a4b6b",
	}

	// Delete pertama → sukses
	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "14212231-792f-4935-bb1c-9a38695a4b6b", deletedUUID)

	// Delete kedua → harus not found
	_, err = handler.Handle(context.Background(), cmd)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Error(t, err)
	assert.Equal(t, "JenisFile.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("JenisFile with identifier %s not found", "14212231-792f-4935-bb1c-9a38695a4b6b"), commonErr.Description)
}

func TestDeleteJenisFileCommand_Fail(t *testing.T) {
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
