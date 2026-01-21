package applicationtest

import (
	"context"
	"testing"

	app "UnpakSiamida/modules/beritaacara/application/CreateBeritaAcara"
	infra "UnpakSiamida/modules/beritaacara/infrastructure"

	"github.com/stretchr/testify/assert"
)

func TestCreateBeritaAcaraCommandValidation_Success(t *testing.T) {
	Tahun := "2080"
	FakultasUnit := 1
	Tanggal := "2021-01-01"
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	validCmd := app.CreateBeritaAcaraCommand{
		Tahun:        Tahun,
		FakultasUnit: FakultasUnit,
		Tanggal:      Tanggal,
		Auditee:      &Auditee,
		Auditor1:     &Auditor1,
		Auditor2:     &Auditor2,
	}
	err := app.CreateBeritaAcaraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestCreateBeritaAcaraCommandValidation_Fail(t *testing.T) {
	Tahun := ""
	FakultasUnit := 0
	Tanggal := ""
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	invalidCmd := app.CreateBeritaAcaraCommand{
		Tahun:        Tahun,
		FakultasUnit: FakultasUnit,
		Tanggal:      Tanggal,
		Auditee:      &Auditee,
		Auditor1:     &Auditor1,
		Auditor2:     &Auditor2,
	}
	err := app.CreateBeritaAcaraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "FakultasUnit cannot be blank")
	assert.Contains(t, err.Error(), "Tanggal cannot be blank")
}

func TestCreateBeritaAcaraCommand_Success(t *testing.T) {
	db, terminate := setupBeritaAcaraMySQL(t)
	defer terminate()

	repo := infra.NewBeritaAcaraRepository(db)
	handler := &app.CreateBeritaAcaraCommandHandler{Repo: repo}

	Tahun := "2080"
	FakultasUnit := 1
	Tanggal := "2021-01-01"
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	cmd := app.CreateBeritaAcaraCommand{
		Tahun:        Tahun,
		FakultasUnit: FakultasUnit,
		Tanggal:      Tanggal,
		Auditee:      &Auditee,
		Auditor1:     &Auditor1,
		Auditor2:     &Auditor2,
	}
	uuidStr, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuidStr)
}
