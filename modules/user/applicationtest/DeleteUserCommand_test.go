package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/user/application/DeleteUser"
	domain "UnpakSiamida/modules/user/domain"
	infra "UnpakSiamida/modules/user/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUserCommandValidation_Success(t *testing.T) {
	validCmd := app.DeleteUserCommand{
		Uuid: uuid.NewString(),
	}
	err := app.DeleteUserCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestDeleteUserCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.DeleteUserCommand{
		Uuid: "",
	}
	err := app.DeleteUserCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

func TestDeleteUserCommand_Success(t *testing.T) {
	db, terminate := setupUserMySQL(t)
	defer terminate()

	repo := infra.NewUserRepository(db)
	handler := &app.DeleteUserCommandHandler{Repo: repo}

	cmd := app.DeleteUserCommand{
		Uuid: "f03ba8b3-918d-4fd2-867d-6943dc14a5ac",
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "f03ba8b3-918d-4fd2-867d-6943dc14a5ac", deletedUUID)

	// Pastikan DB sudah terhapus
	var saved domain.User
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err) // harus error karena sudah dihapus
}

func TestDeleteUserCommand_Fail(t *testing.T) {
	db, terminate := setupUserMySQL(t)
	defer terminate()

	repo := infra.NewUserRepository(db)
	handler := &app.DeleteUserCommandHandler{Repo: repo}

	uuid := uuid.NewString()

	// UUID tidak valid
	cmdInvalidUUID := app.DeleteUserCommand{
		Uuid: uuid,
	}
	_, err := handler.Handle(context.Background(), cmdInvalidUUID)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "User.NotFound", commonErr.Code)
	assert.Contains(t, fmt.Sprintf("User with identifier %s not found", uuid), commonErr.Description)
}
