package applicationtest

import (
	"context"
	"testing"

	app "UnpakSiamida/modules/standarrenstra/application/CreateStandarRenstra"
	domain "UnpakSiamida/modules/standarrenstra/domain"
	infra "UnpakSiamida/modules/standarrenstra/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateStandarRenstraCommandValidation_Success(t *testing.T) {
	validCmd := app.CreateStandarRenstraCommand{Nama: "Dokumen Valid"}
	err := app.CreateStandarRenstraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestCreateStandarRenstraCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.CreateStandarRenstraCommand{Nama: ""}
	err := app.CreateStandarRenstraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Nama cannot be blank")
}

func TestCreateStandarRenstraCommand_Success(t *testing.T) {
	db, terminate := setupStandarRenstraMySQL(t)
	defer terminate()

	repo := infra.NewStandarRenstraRepository(db)
	handler := &app.CreateStandarRenstraCommandHandler{Repo: repo}

	cmd := app.CreateStandarRenstraCommand{Nama: "Dokumen Baru"}
	uuidStr, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr)

	// Pastikan record tersimpan di DB
	var saved domain.StandarRenstra
	err = db.Where("uuid = ?", uuidStr).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, "Dokumen Baru", saved.Nama)
}

func TestCreateStandarRenstraCommand_Edge(t *testing.T) {
	db, terminate := setupStandarRenstraMySQL(t)
	defer terminate()

	repo := infra.NewStandarRenstraRepository(db)
	handler := &app.CreateStandarRenstraCommandHandler{Repo: repo}

	// Insert satu record dulu
	firstCmd := app.CreateStandarRenstraCommand{Nama: "Dokumen Sama"}
	_, err := handler.Handle(context.Background(), firstCmd)
	assert.NoError(t, err)

	// Insert lagi dengan nama yang sama
	secondCmd := app.CreateStandarRenstraCommand{Nama: "Dokumen Sama"}
	uuidStr2, err := handler.Handle(context.Background(), secondCmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr2)

	// UUID harus berbeda
	assert.NotEqual(t, uuidStr2, uuid.Nil)
}
