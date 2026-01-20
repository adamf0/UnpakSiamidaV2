package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/standarrenstra/application/DeleteStandarRenstra"
	domain "UnpakSiamida/modules/standarrenstra/domain"
	infra "UnpakSiamida/modules/standarrenstra/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteStandarRenstraCommandValidation_Success(t *testing.T) {
	validCmd := app.DeleteStandarRenstraCommand{
		Uuid: uuid.NewString(),
	}
	err := app.DeleteStandarRenstraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestDeleteStandarRenstraCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.DeleteStandarRenstraCommand{
		Uuid: "",
	}
	err := app.DeleteStandarRenstraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

func TestDeleteStandarRenstraCommand_Success(t *testing.T) {
	db, terminate := setupStandarRenstraMySQL(t)
	defer terminate()

	repo := infra.NewStandarRenstraRepository(db)
	handler := &app.DeleteStandarRenstraCommandHandler{Repo: repo}

	// Insert record dulu
	record := domain.StandarRenstra{
		UUID: uuid.New(),
		Nama: "Dokumen Untuk Delete",
	}
	err := repo.Create(context.Background(), &record)
	assert.NoError(t, err)

	cmd := app.DeleteStandarRenstraCommand{
		Uuid: record.UUID.String(),
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, record.UUID.String(), deletedUUID)

	// Pastikan DB sudah terhapus
	var saved domain.StandarRenstra
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err) // harus error karena sudah dihapus
}

func TestDeleteStandarRenstraCommand_Edge(t *testing.T) {
	db, terminate := setupStandarRenstraMySQL(t)
	defer terminate()

	repo := infra.NewStandarRenstraRepository(db)
	handler := &app.DeleteStandarRenstraCommandHandler{Repo: repo}

	// Insert record
	record := domain.StandarRenstra{
		UUID: uuid.New(),
		Nama: "Dokumen Edge Delete",
	}
	err := repo.Create(context.Background(), &record)
	assert.NoError(t, err)

	cmd := app.DeleteStandarRenstraCommand{
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
	assert.Equal(t, "StandarRenstra.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("StandarRenstra with identifier %s not found", record.UUID.String()), commonErr.Description)
}

func TestDeleteStandarRenstraCommand_Fail(t *testing.T) {
	db, terminate := setupStandarRenstraMySQL(t)
	defer terminate()

	repo := infra.NewStandarRenstraRepository(db)
	handler := &app.DeleteStandarRenstraCommandHandler{Repo: repo}

	uuid := uuid.NewString()

	// UUID tidak valid
	cmdInvalidUUID := app.DeleteStandarRenstraCommand{
		Uuid: uuid,
	}
	_, err := handler.Handle(context.Background(), cmdInvalidUUID)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "StandarRenstra.NotFound", commonErr.Code)
	assert.Contains(t, fmt.Sprintf("StandarRenstra with identifier %s not found", uuid), commonErr.Description)
}
