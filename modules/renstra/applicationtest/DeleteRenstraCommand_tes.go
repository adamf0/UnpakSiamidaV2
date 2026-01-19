package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/renstra/application/DeleteRenstra"
	infra "UnpakSiamida/modules/renstra/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test handler sukses
func TestDeleteRenstraCommandHandler_Success(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	repo := infra.NewRenstraRepository(db)
	handler := &app.DeleteRenstraCommandHandler{
		Repo: repo,
	}

	uuid, _ := uuid.Parse("c67a37c3-7f25-43de-835d-e4bece0eb308")

	cmd := app.DeleteRenstraCommand{
		Uuid: uuid.String(),
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uuid.String(), deletedUUID)

	// Pastikan sudah terhapus
	found, err := repo.GetByUuid(context.Background(), uuid)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

// Test handler gagal karena UUID invalid
func TestDeleteRenstraCommandHandler_InvalidUUID(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	repo := infra.NewRenstraRepository(db)
	handler := &app.DeleteRenstraCommandHandler{
		Repo: repo,
	}

	cmd := app.DeleteRenstraCommand{
		Uuid: "invalid-uuid",
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID must be a valid UUIDv4 format")
}

// Test handler gagal karena data tidak ditemukan
func TestDeleteRenstraCommandHandler_NotFound(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	repo := infra.NewRenstraRepository(db)
	handler := &app.DeleteRenstraCommandHandler{
		Repo: repo,
	}

	uuid := uuid.NewString()
	cmd := app.DeleteRenstraCommand{
		Uuid: uuid, // UUID valid tapi tidak ada di DB
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Contains(t, commonErr.Code, "Renstra.NotFound")
	assert.Contains(t, commonErr.Description, fmt.Sprintf("Renstra with identifier %s not found", uuid))
}
