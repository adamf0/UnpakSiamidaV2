package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/standarrenstra/application/UpdateStandarRenstra"
	domain "UnpakSiamida/modules/standarrenstra/domain"
	infra "UnpakSiamida/modules/standarrenstra/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateStandarRenstraCommandValidation_Success(t *testing.T) {
	validCmd := app.UpdateStandarRenstraCommand{
		Uuid: uuid.NewString(),
		Nama: "Dokumen Valid",
	}
	err := app.UpdateStandarRenstraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestUpdateStandarRenstraCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.UpdateStandarRenstraCommand{
		Uuid: "",
		Nama: "",
	}
	err := app.UpdateStandarRenstraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
	assert.Contains(t, err.Error(), "Nama cannot be blank")
}

func TestUpdateStandarRenstraCommand_Success(t *testing.T) {
	db, terminate := setupStandarRenstraMySQL(t)
	defer terminate()

	repo := infra.NewStandarRenstraRepository(db)
	handler := &app.UpdateStandarRenstraCommandHandler{Repo: repo}

	// Insert dulu record awal
	original := domain.StandarRenstra{
		UUID: uuid.New(),
		Nama: "Dokumen Lama",
	}
	err := repo.Create(context.Background(), &original)
	assert.NoError(t, err)

	// Update record
	cmd := app.UpdateStandarRenstraCommand{
		Uuid: original.UUID.String(),
		Nama: "Dokumen Baru",
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, original.UUID.String(), updatedUUID)

	// Pastikan DB sudah terupdate
	var saved domain.StandarRenstra
	err = db.Where("uuid = ?", updatedUUID).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, "Dokumen Baru", saved.Nama)
}

func TestUpdateStandarRenstraCommand_Edge(t *testing.T) {
	db, terminate := setupStandarRenstraMySQL(t)
	defer terminate()

	repo := infra.NewStandarRenstraRepository(db)
	handler := &app.UpdateStandarRenstraCommandHandler{Repo: repo}

	// Insert record awal
	original := domain.StandarRenstra{
		UUID: uuid.New(),
		Nama: "Dokumen Edge",
	}
	err := repo.Create(context.Background(), &original)
	assert.NoError(t, err)

	// Update dengan nama sangat panjang (boundary)
	longName := "Dokumen " + string(make([]byte, 500))
	cmdLong := app.UpdateStandarRenstraCommand{
		Uuid: original.UUID.String(),
		Nama: longName,
	}
	updatedUUID, err := handler.Handle(context.Background(), cmdLong)
	assert.NoError(t, err)
	assert.Equal(t, original.UUID.String(), updatedUUID)
}

func TestUpdateStandarRenstraCommand_Fail(t *testing.T) {
	db, terminate := setupStandarRenstraMySQL(t)
	defer terminate()

	repo := infra.NewStandarRenstraRepository(db)
	handler := &app.UpdateStandarRenstraCommandHandler{Repo: repo}

	// Insert record awal
	original := domain.StandarRenstra{
		UUID: uuid.New(),
		Nama: "Dokumen Edge",
	}
	err := repo.Create(context.Background(), &original)
	assert.NoError(t, err)

	uuid := uuid.NewString()
	// Update dengan nama yang sama
	cmdSame := app.UpdateStandarRenstraCommand{
		Uuid: uuid,
		Nama: "Dokumen Edge",
	}
	_, err = handler.Handle(context.Background(), cmdSame)
	assert.Error(t, err)

	commonErr, _ := err.(common.Error)

	assert.Equal(t, "StandarRenstra.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("StandarRenstra with identifier %s not found", uuid), commonErr.Description)
}
