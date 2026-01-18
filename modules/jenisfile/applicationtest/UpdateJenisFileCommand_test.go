package applicationtest

import (
	"context"
	"testing"

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

	// Insert dulu record awal
	original := domain.JenisFile{
		UUID: uuid.New(),
		Nama: "Dokumen Lama",
	}
	err := repo.Create(context.Background(), &original)
	assert.NoError(t, err)

	// Update record
	cmd := app.UpdateJenisFileCommand{
		Uuid: original.UUID.String(),
		Nama: "Dokumen Baru",
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, original.UUID.String(), updatedUUID)

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

	// Insert record awal
	original := domain.JenisFile{
		UUID: uuid.New(),
		Nama: "Dokumen Edge",
	}
	err := repo.Create(context.Background(), &original)
	assert.NoError(t, err)

	// Update dengan nama yang sama
	cmdSame := app.UpdateJenisFileCommand{
		Uuid: original.UUID.String(),
		Nama: "Dokumen Edge",
	}
	updatedUUID, err := handler.Handle(context.Background(), cmdSame)
	assert.NoError(t, err)
	assert.Equal(t, original.UUID.String(), updatedUUID)

	// Update dengan nama sangat panjang (boundary)
	longName := "Dokumen " + string(make([]byte, 500))
	cmdLong := app.UpdateJenisFileCommand{
		Uuid: original.UUID.String(),
		Nama: longName,
	}
	updatedUUID2, err := handler.Handle(context.Background(), cmdLong)
	assert.NoError(t, err)
	assert.Equal(t, original.UUID.String(), updatedUUID2)
}

func TestUpdateJenisFileCommand_Fail(t *testing.T) {
	db, terminate := setupJenisFileMySQL(t)
	defer terminate()

	repo := infra.NewJenisFileRepository(db)
	handler := &app.UpdateJenisFileCommandHandler{Repo: repo}

	// UUID tidak valid
	cmdInvalidUUID := app.UpdateJenisFileCommand{
		Uuid: "invalid-uuid",
		Nama: "Nama Baru",
	}
	_, err := handler.Handle(context.Background(), cmdInvalidUUID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid uuid")

	// UUID valid tapi record tidak ada
	cmdNotFound := app.UpdateJenisFileCommand{
		Uuid: uuid.NewString(),
		Nama: "Nama Baru",
	}
	_, err = handler.Handle(context.Background(), cmdNotFound)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
