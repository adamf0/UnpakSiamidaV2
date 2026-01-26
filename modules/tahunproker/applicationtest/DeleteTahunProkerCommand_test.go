package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/tahunproker/application/DeleteTahunProker"
	domain "UnpakSiamida/modules/tahunproker/domain"
	infra "UnpakSiamida/modules/tahunproker/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTahunProkerCommandValidation_Success(t *testing.T) {
	validCmd := app.DeleteTahunProkerCommand{
		Uuid: uuid.NewString(),
	}
	err := app.DeleteTahunProkerCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestDeleteTahunProkerCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.DeleteTahunProkerCommand{
		Uuid: "",
	}
	err := app.DeleteTahunProkerCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

func TestDeleteTahunProkerCommand_Success(t *testing.T) {
	db, terminate := setupTahunProkerMySQL(t)
	defer terminate()

	repo := infra.NewTahunProkerRepository(db)
	handler := &app.DeleteTahunProkerCommandHandler{Repo: repo}

	cmd := app.DeleteTahunProkerCommand{
		Uuid: "666a6b72-d2b4-481f-adb8-298d807e9e20",
	}

	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "666a6b72-d2b4-481f-adb8-298d807e9e20", deletedUUID)

	// Pastikan DB sudah terhapus
	var saved domain.TahunProker
	err = db.Where("uuid = ?", deletedUUID).First(&saved).Error
	assert.Error(t, err) // harus error karena sudah dihapus
}

func TestDeleteTahunProkerCommand_Edge(t *testing.T) {
	db, terminate := setupTahunProkerMySQL(t)
	defer terminate()

	repo := infra.NewTahunProkerRepository(db)
	handler := &app.DeleteTahunProkerCommandHandler{Repo: repo}

	cmd := app.DeleteTahunProkerCommand{
		Uuid: "666a6b72-d2b4-481f-adb8-298d807e9e20",
	}

	// Delete pertama → sukses
	deletedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "666a6b72-d2b4-481f-adb8-298d807e9e20", deletedUUID)

	// Delete kedua → harus not found
	_, err = handler.Handle(context.Background(), cmd)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Error(t, err)
	assert.Equal(t, "TahunProker.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("TahunProker with identifier %s not found", "666a6b72-d2b4-481f-adb8-298d807e9e20"), commonErr.Description)
}

func TestDeleteTahunProkerCommand_Fail(t *testing.T) {
	db, terminate := setupTahunProkerMySQL(t)
	defer terminate()

	repo := infra.NewTahunProkerRepository(db)
	handler := &app.DeleteTahunProkerCommandHandler{Repo: repo}

	uuid := uuid.NewString()

	// UUID tidak valid
	cmdInvalidUUID := app.DeleteTahunProkerCommand{
		Uuid: uuid,
	}
	_, err := handler.Handle(context.Background(), cmdInvalidUUID)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "TahunProker.NotFound", commonErr.Code)
	assert.Contains(t, fmt.Sprintf("TahunProker with identifier %s not found", uuid), commonErr.Description)
}
