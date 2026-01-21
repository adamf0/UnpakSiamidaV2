package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/jenisfile/application/UpdateJenisFile"
	domain "UnpakSiamida/modules/jenisfile/domain"
	infra "UnpakSiamida/modules/jenisfile/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateJenisFileCommandValidation_Success(t *testing.T) {
	validCmd := app.UpdateJenisFileCommand{
		Uuid: uuid.NewString(),
		Nama: "Dokumen Valid",
	}
	err := app.UpdateJenisFileCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestUpdateJenisFileCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.UpdateJenisFileCommand{
		Uuid: "",
		Nama: "",
	}
	err := app.UpdateJenisFileCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
	assert.Contains(t, err.Error(), "Nama cannot be blank")
}

func TestUpdateJenisFileCommand_Success(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.UpdateJenisFileCommandHandler{Repo: repo}

	// Update record
	cmd := app.UpdateJenisFileCommand{
		Uuid: "14212231-792f-4935-bb1c-9a38695a4b6b",
		Nama: "Dokumen Baru",
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "14212231-792f-4935-bb1c-9a38695a4b6b", updatedUUID)

	// Pastikan DB sudah terupdate
	var saved domain.JenisFile
	err = db.Where("uuid = ?", updatedUUID).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, "Dokumen Baru", saved.Nama)
}

func TestUpdateJenisFileCommand_Edge(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.UpdateJenisFileCommandHandler{Repo: repo}

	// Update dengan nama sangat panjang (boundary)
	longName := "Dokumen " + string(make([]byte, 500))
	cmdLong := app.UpdateJenisFileCommand{
		Uuid: "14212231-792f-4935-bb1c-9a38695a4b6b",
		Nama: longName,
	}
	updatedUUID, err := handler.Handle(context.Background(), cmdLong)
	assert.NoError(t, err)
	assert.Equal(t, "14212231-792f-4935-bb1c-9a38695a4b6b", updatedUUID)
}

func TestUpdateJenisFileCommand_Fail(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.UpdateJenisFileCommandHandler{Repo: repo}

	// Insert record awal
	original := domain.JenisFile{
		UUID: uuid.New(),
		Nama: "Dokumen Edge",
	}
	err := repo.Create(context.Background(), &original)
	assert.NoError(t, err)

	uuid := uuid.NewString()
	// Update dengan nama yang sama
	cmdSame := app.UpdateJenisFileCommand{
		Uuid: uuid,
		Nama: "Dokumen Edge",
	}
	_, err = handler.Handle(context.Background(), cmdSame)
	assert.Error(t, err)

	commonErr, _ := err.(common.Error)

	assert.Equal(t, "JenisFile.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("JenisFile with identifier %s not found", uuid), commonErr.Description)
}
