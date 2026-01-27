package applicationtest

import (
	"context"
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

func TestUpdateTahunProkerCommand_Fail(t *testing.T) {
	tests := []struct {
		name      string
		uuid      string
		tahun     string
		status    string
		errorCode string
	}{
		{
			name:      "not found",
			uuid:      uuid.NewString(),
			tahun:     "2080",
			status:    "non-aktif",
			errorCode: "TahunProker.NotFound",
		},
		{
			name:      "invalid tahun",
			uuid:      "666a6b72-d2b4-481f-adb8-298d807e9e20",
			tahun:     "1900",
			status:    "non-aktif",
			errorCode: "TahunProker.InvalidTahun",
		},
		{
			name:      "invalid status",
			uuid:      "666a6b72-d2b4-481f-adb8-298d807e9e20",
			tahun:     "2001",
			status:    "no-aktif",
			errorCode: "TahunProker.InvalidStatus",
		},
		{
			name:      "invalid status",
			uuid:      "666a6b72-d2b4-481f-adb8-298d807e9e20",
			tahun:     "2024",
			status:    "non-aktif",
			errorCode: "TahunProker.DuplicateData",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, terminate := setupTahunProkerMySQL(t)
			defer terminate()

			repo := infra.NewTahunProkerRepository(db)
			handler := &app.UpdateTahunProkerCommandHandler{
				Repo: repo,
			}

			cmd := app.UpdateTahunProkerCommand{
				Uuid:   tt.uuid,
				Tahun:  tt.tahun,
				Status: tt.status,
			}

			_, err := handler.Handle(context.Background(), cmd)
			assert.Error(t, err)

			commonErr, ok := err.(common.Error)
			assert.True(t, ok)

			assert.Equal(t, tt.errorCode, commonErr.Code)
		})
	}
}
