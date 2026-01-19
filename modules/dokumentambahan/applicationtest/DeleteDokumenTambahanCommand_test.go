package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/dokumentambahan/application/DeleteDokumenTambahan"
	domain "UnpakSiamida/modules/dokumentambahan/domain"
	infra "UnpakSiamida/modules/dokumentambahan/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteDokumenTambahanCommandValidation_Success(t *testing.T) {
	validCmd := app.DeleteDokumenTambahanCommand{
		Uuid: uuid.NewString(),
	}
	err := app.DeleteDokumenTambahanCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestDeleteDokumenTambahanCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.DeleteDokumenTambahanCommand{
		Uuid: "",
	}
	err := app.DeleteDokumenTambahanCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

func TestDeleteDokumenTambahanCommandHandler_Success(t *testing.T) {
	db, terminate := setupDokumenTambahanMySQL(t)
	defer terminate()

	repo := infra.NewDokumenTambahanRepository(db)
	handler := &app.DeleteDokumenTambahanCommandHandler{Repo: repo}

	cmd := app.DeleteDokumenTambahanCommand{
		Uuid: "802b0732-e5b5-4852-a770-a834a8b70746",
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "802b0732-e5b5-4852-a770-a834a8b70746", deletedUUID)

	// Pastikan DB sudah terhapus
	var saved domain.DokumenTambahan
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err)
}

func TestDeleteDokumenTambahanCommandHandler_Edge(t *testing.T) {
	db, terminate := setupDokumenTambahanMySQL(t)
	defer terminate()

	repo := infra.NewDokumenTambahanRepository(db)
	handler := &app.DeleteDokumenTambahanCommandHandler{Repo: repo}

	cmd := app.DeleteDokumenTambahanCommand{
		Uuid: "864285ba-3b78-4aaa-bbb3-02b162af12a6",
	}

	// Delete pertama → sukses
	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "864285ba-3b78-4aaa-bbb3-02b162af12a6", deletedUUID)

	// Delete kedua → harus not found
	_, err = handler.Handle(context.Background(), cmd)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, common.NotFound, commonErr.Type)
	assert.Contains(t, commonErr.Description, "DokumenTambahan with identifier 864285ba-3b78-4aaa-bbb3-02b162af12a6 not found")
}

func TestDeleteDokumenTambahanCommandHandler_Fail(t *testing.T) {
	db, terminate := setupDokumenTambahanMySQL(t)
	defer terminate()

	repo := infra.NewDokumenTambahanRepository(db)
	handler := &app.DeleteDokumenTambahanCommandHandler{Repo: repo}

	uuid := uuid.NewString()
	cmdInvalid := app.DeleteDokumenTambahanCommand{
		Uuid: uuid,
	}
	_, err := handler.Handle(context.Background(), cmdInvalid)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "DokumenTambahan.NotFound", commonErr.Code)
	assert.Contains(t, commonErr.Description, fmt.Sprintf("Renstra with identifier %s not found", uuid))
}
