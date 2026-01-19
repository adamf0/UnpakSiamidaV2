package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/renstranilai/application/DeleteRenstraNilai"
	"UnpakSiamida/modules/renstranilai/domain"
	infra "UnpakSiamida/modules/renstranilai/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test handler sukses
func TestDeleteRenstraNilaiCommandHandler_Success(t *testing.T) {
	db, terminate := setupRenstraNilaiMySQL(t)
	defer terminate()

	repo := infra.NewRenstraNilaiRepository(db)
	handler := &app.DeleteRenstraNilaiCommandHandler{
		Repo: repo,
	}

	uuid, _ := uuid.Parse("7b59ffc3-f851-4e38-ba96-06cc168c8dd1")

	cmd := app.DeleteRenstraNilaiCommand{
		Uuid: uuid.String(),
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uuid.String(), deletedUUID)

	var saved domain.RenstraNilai
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err) // harus error karena sudah dihapus
}

// Test handler gagal karena UUID invalid
func TestDeleteRenstraNilaiCommandHandler_InvalidUUID(t *testing.T) {
	cmd := app.DeleteRenstraNilaiCommand{
		Uuid: "invalid-uuid",
	}

	err := app.DeleteRenstraNilaiCommandValidation(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID must be a valid UUIDv4 format")
}

// Test handler gagal karena data tidak ditemukan
func TestDeleteRenstraNilaiCommandHandler_NotFound(t *testing.T) {
	db, terminate := setupRenstraNilaiMySQL(t)
	defer terminate()

	repo := infra.NewRenstraNilaiRepository(db)
	handler := &app.DeleteRenstraNilaiCommandHandler{
		Repo: repo,
	}

	uuid := uuid.NewString()
	cmd := app.DeleteRenstraNilaiCommand{
		Uuid: uuid, // UUID valid tapi tidak ada di DB
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Contains(t, commonErr.Code, "RenstraNilai.NotFound")
	assert.Contains(t, commonErr.Description, fmt.Sprintf("RenstraNilai with identifier %s not found", uuid))
}
