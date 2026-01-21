package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/beritaacara/application/DeleteBeritaAcara"
	domain "UnpakSiamida/modules/beritaacara/domain"
	infra "UnpakSiamida/modules/beritaacara/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBeritaAcaraCommandValidation_Success(t *testing.T) {
	validCmd := app.DeleteBeritaAcaraCommand{
		Uuid: uuid.NewString(),
	}
	err := app.DeleteBeritaAcaraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestDeleteBeritaAcaraCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.DeleteBeritaAcaraCommand{
		Uuid: "",
	}
	err := app.DeleteBeritaAcaraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

func TestDeleteBeritaAcaraCommand_Success(t *testing.T) {
	db, terminate := setupBeritaAcaraMySQL(t)
	defer terminate()

	repo := infra.NewBeritaAcaraRepository(db)
	handler := &app.DeleteBeritaAcaraCommandHandler{Repo: repo}

	cmd := app.DeleteBeritaAcaraCommand{
		Uuid: "14212231-792f-4935-bb1c-9a38695a4b6b",
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "14212231-792f-4935-bb1c-9a38695a4b6b", deletedUUID)

	// Pastikan DB sudah terhapus
	var saved domain.BeritaAcara
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err) // harus error karena sudah dihapus
}

func TestDeleteBeritaAcaraCommand_Edge(t *testing.T) {
	db, terminate := setupBeritaAcaraMySQL(t)
	defer terminate()

	repo := infra.NewBeritaAcaraRepository(db)
	handler := &app.DeleteBeritaAcaraCommandHandler{Repo: repo}

	cmd := app.DeleteBeritaAcaraCommand{
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
	assert.Equal(t, "BeritaAcara.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("BeritaAcara with identifier %s not found", "14212231-792f-4935-bb1c-9a38695a4b6b"), commonErr.Description)
}

func TestDeleteBeritaAcaraCommand_Fail(t *testing.T) {
	db, terminate := setupBeritaAcaraMySQL(t)
	defer terminate()

	repo := infra.NewBeritaAcaraRepository(db)
	handler := &app.DeleteBeritaAcaraCommandHandler{Repo: repo}

	uuid := uuid.NewString()

	// UUID tidak valid
	cmdInvalidUUID := app.DeleteBeritaAcaraCommand{
		Uuid: uuid,
	}
	_, err := handler.Handle(context.Background(), cmdInvalidUUID)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "BeritaAcara.NotFound", commonErr.Code)
	assert.Contains(t, fmt.Sprintf("BeritaAcara with identifier %s not found", uuid), commonErr.Description)
}
