package applicationtest

import (
	"context"
	"testing"

	app "UnpakSiamida/modules/tahunproker/application/CreateTahunProker"
	domain "UnpakSiamida/modules/tahunproker/domain"
	infra "UnpakSiamida/modules/tahunproker/infrastructure"

	"github.com/stretchr/testify/assert"
)

func TestCreateTahunProkerCommandValidation_Success(t *testing.T) {
	validCmd := app.CreateTahunProkerCommand{
		Tahun:  "2080",
		Status: "aktif",
	}
	err := app.CreateTahunProkerCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestCreateTahunProkerCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.CreateTahunProkerCommand{
		Tahun:  "",
		Status: "",
	}
	err := app.CreateTahunProkerCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
}

func TestCreateTahunProkerCommand_Success(t *testing.T) {
	db, terminate := setupTahunProkerMySQL(t)
	defer terminate()

	repo := infra.NewTahunProkerRepository(db)
	handler := &app.CreateTahunProkerCommandHandler{Repo: repo}

	cmd := app.CreateTahunProkerCommand{Tahun: "2080"}
	uuidStr, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr)

	// Pastikan record tersimpan di DB
	var saved domain.TahunProker
	err = db.Where("uuid = ?", uuidStr).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, "2080", saved.Tahun)
}

// func TestCreateTahunProkerCommand_Edge(t *testing.T) {
// 	db, terminate := setupTahunProkerMySQL(t)
// 	defer terminate()

// 	repo := infra.NewTahunProkerRepository(db)
// 	handler := &app.CreateTahunProkerCommandHandler{Repo: repo}

// 	// Insert satu record dulu
// 	firstCmd := app.CreateTahunProkerCommand{
// 		Tahun:  "2080",
// 		Status: "aktif",
// 	}
// 	_, err := handler.Handle(context.Background(), firstCmd)
// 	assert.NoError(t, err)

// 	// Insert lagi dengan nama yang sama
// 	secondCmd := app.CreateTahunProkerCommand{
// 		Tahun:  "2080",
// 		Status: "aktif",
// 	}
// 	uuidStr2, err := handler.Handle(context.Background(), secondCmd)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, uuidStr2)

// 	// UUID harus berbeda
// 	assert.NotEqual(t, uuidStr2, uuid.Nil)
// }
