package applicationtest

import (
	"context"
	"testing"

	common "UnpakSiamida/common/domain"
	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
	app "UnpakSiamida/modules/user/application/CreateUser"
	domain "UnpakSiamida/modules/user/domain"
	infra "UnpakSiamida/modules/user/infrastructure"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserCommandValidation_Success(t *testing.T) {
	uuidFakultas := "dea9a83f-70b3-4295-85ed-459eb1a9f6a0"

	validCmd := app.CreateUserCommand{
		Name:             "adminf",
		Username:         "adminf",
		Password:         "123",
		Email:            "adminf@unpak.ac.id",
		Level:            "fakultas",
		UuidFakultasUnit: &uuidFakultas,
	}
	err := app.CreateUserCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestCreateUserCommandValidation_Fail(t *testing.T) {
	uuidFakultas := ""

	invalidCmd := app.CreateUserCommand{
		Name:             "",
		Username:         "",
		Password:         "",
		Email:            "",
		Level:            "",
		UuidFakultasUnit: &uuidFakultas,
	}

	err := app.CreateUserCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Username cannot be blank")
	assert.Contains(t, err.Error(), "Password cannot be blank")
	assert.Contains(t, err.Error(), "Name cannot be blank")
	assert.Contains(t, err.Error(), "Email cannot be blank")
	assert.Contains(t, err.Error(), "Level cannot be blank")
}

func TestCreateUserCommand_Success(t *testing.T) {
	db, terminate := setupUserMySQL(t)
	defer terminate()

	repo := infra.NewUserRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)

	handler := &app.CreateUserCommandHandler{Repo: repo, RepoFakultasUnit: repoFakultasUnit}
	uuidFakultas := "dea9a83f-70b3-4295-85ed-459eb1a9f6a0"

	cmd := app.CreateUserCommand{
		Name:             "adminf",
		Username:         "adminf",
		Password:         "123",
		Email:            "adminf@unpak.ac.id",
		Level:            "fakultas",
		UuidFakultasUnit: &uuidFakultas,
	}
	uuidStr, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr)

	// Pastikan record tersimpan di DB
	var saved domain.User
	err = db.Where("uuid = ?", uuidStr).First(&saved).Error
	assert.NoError(t, err)
}

func TestCreateUserCommand_FailEmail(t *testing.T) {
	db, terminate := setupUserMySQL(t)
	defer terminate()

	repo := infra.NewUserRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)

	handler := &app.CreateUserCommandHandler{Repo: repo, RepoFakultasUnit: repoFakultasUnit}
	uuidFakultas := "dea9a83f-70b3-4295-85ed-459eb1a9f6a0"

	cmdSame := app.CreateUserCommand{
		Name:             "adminf",
		Username:         "adminf",
		Password:         "123",
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

func TestCreateUserCommand_FailFakultasUnit(t *testing.T) {
	db, terminate := setupUserMySQL(t)
	defer terminate()

	repo := infra.NewUserRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)

	handler := &app.CreateUserCommandHandler{Repo: repo, RepoFakultasUnit: repoFakultasUnit}
	uuidFakultas := "dea9a83f-70b3-4295-85ed-000000000000"

	cmdSame := app.CreateUserCommand{
		Name:             "adminf",
		Username:         "adminf",
		Password:         "123",
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
