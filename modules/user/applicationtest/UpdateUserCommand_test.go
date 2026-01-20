package applicationtest

import (
	"context"
	"testing"

	common "UnpakSiamida/common/domain"
	appAdd "UnpakSiamida/modules/user/application/CreateUser"
	app "UnpakSiamida/modules/user/application/UpdateUser"
	infra "UnpakSiamida/modules/user/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserCommandValidation_Success(t *testing.T) {
	uuidFakultas := "dea9a83f-70b3-4295-85ed-459eb1a9f6a0"
	password := "123"

	validCmd := app.UpdateUserCommand{
		Uuid:             uuid.NewString(),
		Name:             "adminf",
		Username:         "adminf",
		Password:         &password,
		Email:            "adminf@unpak.ac.id",
		Level:            "fakultas",
		UuidFakultasUnit: &uuidFakultas,
	}
	err := app.UpdateUserCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestUpdateUserCommandValidation_Fail(t *testing.T) {
	password := ""

	invalidCmd := app.UpdateUserCommand{
		Uuid:     "",
		Name:     "",
		Username: "",
		Password: &password,
		Email:    "",
		Level:    "",
	}
	err := app.UpdateUserCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
	assert.Contains(t, err.Error(), "Username cannot be blank")
	assert.Contains(t, err.Error(), "Password cannot be blank")
	assert.Contains(t, err.Error(), "Name cannot be blank")
	assert.Contains(t, err.Error(), "Email cannot be blank")
	assert.Contains(t, err.Error(), "Level cannot be blank")
}

func TestUpdateUserCommand_Success(t *testing.T) {
	db, terminate := setupUserMySQL(t)
	defer terminate()

	repo := infra.NewUserRepository(db)
	handler := &app.UpdateUserCommandHandler{Repo: repo}
	handlerAdd := &appAdd.CreateUserCommandHandler{Repo: repo}

	uuidFakultas := "dea9a83f-70b3-4295-85ed-459eb1a9f6a0"
	password := "123"

	// Insert dulu record awal
	cmdAdd := appAdd.CreateUserCommand{
		Name:             "adminf",
		Username:         "adminf",
		Password:         "123",
		Email:            "adminf@unpak.ac.id",
		Level:            "fakultas",
		UuidFakultasUnit: &uuidFakultas,
	}
	uuidStr, err := handlerAdd.Handle(context.Background(), cmdAdd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr)

	// Update record
	cmd := app.UpdateUserCommand{
		Uuid:             uuidStr,
		Name:             "adminfh",
		Username:         "adminf",
		Password:         &password,
		Email:            "adminf@unpak.ac.id",
		Level:            "fakultas",
		UuidFakultasUnit: &uuidFakultas,
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uuidStr, updatedUUID)
}

func TestUpdateUserCommand_FailEmail(t *testing.T) {
	db, terminate := setupUserMySQL(t)
	defer terminate()

	repo := infra.NewUserRepository(db)
	handler := &app.UpdateUserCommandHandler{Repo: repo}
	uuidFakultas := "dea9a83f-70b3-4295-85ed-459eb1a9f6a0"
	password := "123"

	uuid := uuid.NewString()
	cmdSame := app.UpdateUserCommand{
		Uuid:             uuid,
		Name:             "adminf",
		Username:         "adminf",
		Password:         &password,
		Email:            "adminf@gmail.com",
		Level:            "fakultas",
		UuidFakultasUnit: &uuidFakultas,
	}
	_, err := handler.Handle(context.Background(), cmdSame)
	assert.Error(t, err)

	commonErr, _ := err.(common.Error)

	assert.Equal(t, "User.InvalidEmail", commonErr.Code)
	assert.Equal(t, "email tidak valid atau tidak diperbolehkan", commonErr.Description)
}

func TestUpdateUserCommand_Fail2FakultasUnit(t *testing.T) {
	db, terminate := setupUserMySQL(t)
	defer terminate()

	repo := infra.NewUserRepository(db)
	handler := &app.UpdateUserCommandHandler{Repo: repo}
	uuidFakultas := "dea9a83f-70b3-4295-85ed-000000000000"
	password := "123"

	uuid := uuid.NewString()
	cmdSame := app.UpdateUserCommand{
		Uuid:             uuid,
		Name:             "adminf",
		Username:         "adminf",
		Password:         &password,
		Email:            "adminf@unpak.ac.id",
		Level:            "fakultas",
		UuidFakultasUnit: &uuidFakultas,
	}
	_, err := handler.Handle(context.Background(), cmdSame)
	assert.Error(t, err)

	commonErr, _ := err.(common.Error)

	assert.Equal(t, "User.InvalidFakultasUnit", commonErr.Code)
	assert.Equal(t, "fakultas unit tidak valid", commonErr.Description)
}
