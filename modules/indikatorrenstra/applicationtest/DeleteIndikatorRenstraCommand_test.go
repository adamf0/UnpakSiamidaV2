package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"
	domain "UnpakSiamida/modules/indikatorrenstra/domain"
	infra "UnpakSiamida/modules/indikatorrenstra/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteIndikatorRenstraCommandValidation_Success(t *testing.T) {
	validCmd := app.DeleteIndikatorRenstraCommand{
		Uuid: uuid.NewString(),
	}
	err := app.DeleteIndikatorRenstraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestDeleteIndikatorRenstraCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.DeleteIndikatorRenstraCommand{
		Uuid: "",
	}
	err := app.DeleteIndikatorRenstraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

func TestDeleteIndikatorRenstraCommand_Success(t *testing.T) {
	db, terminate := setupIndikatorRenstraMySQL(t)
	defer terminate()

	repo := infra.NewIndikatorRenstraRepository(db)
	handler := &app.DeleteIndikatorRenstraCommandHandler{Repo: repo}

	cmd := app.DeleteIndikatorRenstraCommand{
		Uuid: "b763b5b3-a18e-416c-9d0d-a0c23aa6076c",
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "b763b5b3-a18e-416c-9d0d-a0c23aa6076c", deletedUUID)

	// Pastikan DB sudah terhapus
	var saved domain.IndikatorRenstra
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err) // harus error karena sudah dihapus
}

func TestDeleteIndikatorRenstraCommand_Edge(t *testing.T) {
	db, terminate := setupIndikatorRenstraMySQL(t)
	defer terminate()

	repo := infra.NewIndikatorRenstraRepository(db)
	handler := &app.DeleteIndikatorRenstraCommandHandler{Repo: repo}

	cmd := app.DeleteIndikatorRenstraCommand{
		Uuid: "b763b5b3-a18e-416c-9d0d-a0c23aa6076c",
	}

	// Delete pertama → sukses
	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "b763b5b3-a18e-416c-9d0d-a0c23aa6076c", deletedUUID)

	// Delete kedua → harus not found
	_, err = handler.Handle(context.Background(), cmd)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Error(t, err)
	assert.Equal(t, "IndikatorRenstra.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("IndikatorRenstra with identifier %s not found", "b763b5b3-a18e-416c-9d0d-a0c23aa6076c"), commonErr.Description)
}

func TestDeleteIndikatorRenstraCommand_Fail(t *testing.T) {
	db, terminate := setupIndikatorRenstraMySQL(t)
	defer terminate()

	repo := infra.NewIndikatorRenstraRepository(db)
	handler := &app.DeleteIndikatorRenstraCommandHandler{Repo: repo}

	uuid := uuid.NewString()

	// UUID tidak valid
	cmdInvalidUUID := app.DeleteIndikatorRenstraCommand{
		Uuid: uuid,
	}
	_, err := handler.Handle(context.Background(), cmdInvalidUUID)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "IndikatorRenstra.NotFound", commonErr.Code)
	assert.Contains(t, fmt.Sprintf("IndikatorRenstra with identifier %s not found", uuid), commonErr.Description)
}
