package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/renstra/application/DeleteRenstra"
	"UnpakSiamida/modules/renstra/domain"
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

	uuid, _ := uuid.Parse("0025699d-d69b-41e4-b712-f437aa15d3b1")

	cmd := app.DeleteRenstraCommand{
		Uuid: uuid.String(),
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uuid.String(), deletedUUID)

	var saved domain.Renstra
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err) // harus error karena sudah dihapus
}

// Test handler gagal karena UUID invalid
func TestDeleteRenstraCommandHandler_InvalidUUID(t *testing.T) {
	cmd := app.DeleteRenstraCommand{
		Uuid: "invalid-uuid",
	}

	err := app.DeleteRenstraCommandValidation(cmd)
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
