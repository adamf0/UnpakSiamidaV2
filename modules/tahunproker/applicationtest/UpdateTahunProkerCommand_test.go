package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/tahunproker/application/UpdateTahunProker"
	domain "UnpakSiamida/modules/tahunproker/domain"
	infra "UnpakSiamida/modules/tahunproker/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTahunProkerCommandValidation_Success(t *testing.T) {
	validCmd := app.UpdateTahunProkerCommand{
		Uuid:   uuid.NewString(),
		Tahun:  "2080",
		Status: "aktif",
	}
	err := app.UpdateTahunProkerCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestUpdateTahunProkerCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.UpdateTahunProkerCommand{
		Uuid:   "",
		Tahun:  "",
		Status: "",
	}
	err := app.UpdateTahunProkerCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "Status cannot be blank")
}

func TestUpdateTahunProkerCommand_Success(t *testing.T) {
	db, terminate := setupTahunProkerMySQL(t)
	defer terminate()

	repo := infra.NewTahunProkerRepository(db)
	handler := &app.UpdateTahunProkerCommandHandler{Repo: repo}

	// Update record
	cmd := app.UpdateTahunProkerCommand{
		Uuid:   "666a6b72-d2b4-481f-adb8-298d807e9e20",
		Tahun:  "2080",
		Status: "aktif",
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "666a6b72-d2b4-481f-adb8-298d807e9e20", updatedUUID)

	// Pastikan DB sudah terupdate
	var saved domain.TahunProker
	err = db.Where("uuid = ?", updatedUUID).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, "2080", saved.Tahun)
}

func TestUpdateTahunProkerCommand_Edge(t *testing.T) {
	db, terminate := setupTahunProkerMySQL(t)
	defer terminate()

	repo := infra.NewTahunProkerRepository(db)
	handler := &app.UpdateTahunProkerCommandHandler{Repo: repo}

	// Update dengan nama sangat panjang (boundary)
	longName := "Dokumen " + string(make([]byte, 500))
	cmdLong := app.UpdateTahunProkerCommand{
		Uuid:  "666a6b72-d2b4-481f-adb8-298d807e9e20",
		Tahun: longName,
	}
	updatedUUID, err := handler.Handle(context.Background(), cmdLong)
	assert.NoError(t, err)
	assert.Equal(t, "666a6b72-d2b4-481f-adb8-298d807e9e20", updatedUUID)
}

func TestUpdateTahunProkerCommand_Fail(t *testing.T) {
	db, terminate := setupTahunProkerMySQL(t)
	defer terminate()

	repo := infra.NewTahunProkerRepository(db)
	handler := &app.UpdateTahunProkerCommandHandler{Repo: repo}

	// Insert record awal
	original := domain.TahunProker{
		UUID:  uuid.New(),
		Tahun: "Dokumen Edge",
	}
	err := repo.Create(context.Background(), &original)
	assert.NoError(t, err)

	uuid := uuid.NewString()
	// Update dengan nama yang sama
	cmdSame := app.UpdateTahunProkerCommand{
		Uuid:  uuid,
		Tahun: "Dokumen Edge",
	}
	_, err = handler.Handle(context.Background(), cmdSame)
	assert.Error(t, err)

	commonErr, _ := err.(common.Error)

	assert.Equal(t, "TahunProker.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("TahunProker with identifier %s not found", uuid), commonErr.Description)
}
