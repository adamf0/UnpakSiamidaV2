package applicationtest

import (
	"context"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
	domain "UnpakSiamida/modules/indikatorrenstra/domain"
	infra "UnpakSiamida/modules/indikatorrenstra/infrastructure"
	infraStandar "UnpakSiamida/modules/standarrenstra/infrastructure"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndikatorRenstraCommandValidation_Success(t *testing.T) {
	validCmd := app.CreateIndikatorRenstraCommand{
		StandarRenstra: "5fd713d0-adfe-4086-a000-21c948faf84d",
		Indikator:      "Lulusan memiliki sertifikat kompetensi atau Bahasa asing",
		Parent:         nil,
		Tahun:          "2080",
		TipeTarget:     "numerik",
		Operator:       nil,
	}
	err := app.CreateIndikatorRenstraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestCreateIndikatorRenstraCommandValidation_Fail(t *testing.T) {
	invalidCmd := app.CreateIndikatorRenstraCommand{
		StandarRenstra: "",
		Indikator:      "",
		Parent:         nil,
		Tahun:          "",
		TipeTarget:     "",
		Operator:       nil,
	}
	err := app.CreateIndikatorRenstraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Standar Renstra cannot be blank")
	assert.Contains(t, err.Error(), "Indikator cannot be blank")
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "Tipe Target cannot be blank")
}

func TestCreateIndikatorRenstraCommand_Success(t *testing.T) {
	db, terminate := setupIndikatorRenstraMySQL(t)
	defer terminate()

	handler := &app.CreateIndikatorRenstraCommandHandler{
		Repo:               infra.NewIndikatorRenstraRepository(db),
		RepoStandarRenstra: infraStandar.NewStandarRenstraRepository(db),
	}

	cmd := app.CreateIndikatorRenstraCommand{
		StandarRenstra: "5fd713d0-adfe-4086-a000-21c948faf84d",
		Indikator:      "Lulusan memiliki sertifikat kompetensi atau Bahasa asing",
		Parent:         nil,
		Tahun:          "2080",
		TipeTarget:     "numerik",
		Operator:       nil,
	}
	uuidStr, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr)

	// Pastikan record tersimpan di DB
	var saved domain.IndikatorRenstra
	err = db.Where("uuid = ?", uuidStr).First(&saved).Error
	assert.NoError(t, err)
}

func TestCreateIndikatorRenstraCommand_Fail(t *testing.T) {
	tests := []struct {
		name         string
		cmd          app.CreateIndikatorRenstraCommand
		expectedCode string
	}{
		{
			name: "InvalidStandar",
			cmd: app.CreateIndikatorRenstraCommand{
				StandarRenstra: "5fd713d0-adfe-4086-a000-000000000000",
				Indikator:      "Lulusan memiliki sertifikat kompetensi atau Bahasa asing",
				Tahun:          "2080",
				TipeTarget:     "numerik",
			},
			expectedCode: "IndikatorRenstra.InvalidStandar",
		},
		{
			name: "NotUniqueIndikator",
			cmd: app.CreateIndikatorRenstraCommand{
				StandarRenstra: "5fd713d0-adfe-4086-a000-21c948faf84d",
				Indikator:      "Lulusan memiliki sertifikat kompetensi atau Bahasa asing",
				Tahun:          "2024",
				TipeTarget:     "numerik",
			},
			expectedCode: "IndikatorRenstra.NotUniqueIndikator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, terminate := setupIndikatorRenstraMySQL(t)
			defer terminate()

			handler := &app.CreateIndikatorRenstraCommandHandler{
				Repo:               infra.NewIndikatorRenstraRepository(db),
				RepoStandarRenstra: infraStandar.NewStandarRenstraRepository(db),
			}

			_, err := handler.Handle(context.Background(), tt.cmd)
			assert.Error(t, err)

			commonErr, ok := err.(common.Error)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, commonErr.Code)
		})
	}
}

func TestCreateTemplateRenstra_ContextTimeout(t *testing.T) {
	db, terminate := setupIndikatorRenstraMySQL(t)
	defer terminate()

	handler := &app.CreateIndikatorRenstraCommandHandler{
		Repo:               infra.NewIndikatorRenstraRepository(db),
		RepoStandarRenstra: infraStandar.NewStandarRenstraRepository(db),
	}

	cmd := app.CreateIndikatorRenstraCommand{
		StandarRenstra: "5fd713d0-adfe-4086-a000-21c948faf84d",
		Indikator:      "Lulusan memiliki sertifikat kompetensi atau Bahasa asing",
		Parent:         nil,
		Tahun:          "2080",
		TipeTarget:     "numerik",
		Operator:       nil,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.True(t, err == context.Canceled || err == context.DeadlineExceeded, "expected context canceled or timeout error")
}
