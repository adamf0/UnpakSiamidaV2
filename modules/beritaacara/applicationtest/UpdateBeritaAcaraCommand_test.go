package applicationtest

import (
	"context"
	"testing"

	app "UnpakSiamida/modules/beritaacara/application/UpdateBeritaAcara"
	infra "UnpakSiamida/modules/beritaacara/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBeritaAcaraCommandValidation_Success(t *testing.T) {
	Tahun := "2080"
	FakultasUnit := 1
	Tanggal := "2021-01-01"
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	validCmd := app.UpdateBeritaAcaraCommand{
		Uuid:         uuid.NewString(),
		Tahun:        Tahun,
		FakultasUnit: FakultasUnit,
		Tanggal:      Tanggal,
		Auditee:      &Auditee,
		Auditor1:     &Auditor1,
		Auditor2:     &Auditor2,
	}
	err := app.UpdateBeritaAcaraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestUpdateBeritaAcaraCommandValidation_Fail(t *testing.T) {
	Tahun := ""
	FakultasUnit := 0
	Tanggal := ""
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	invalidCmd := app.UpdateBeritaAcaraCommand{
		Uuid:         "",
		Tahun:        Tahun,
		FakultasUnit: FakultasUnit,
		Tanggal:      Tanggal,
		Auditee:      &Auditee,
		Auditor1:     &Auditor1,
		Auditor2:     &Auditor2,
	}
	err := app.UpdateBeritaAcaraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "FakultasUnit cannot be blank")
	assert.Contains(t, err.Error(), "Tanggal cannot be blank")
}

func TestUpdateBeritaAcaraCommand_Success(t *testing.T) {
	db, terminate := setupBeritaAcaraMySQL(t)
	defer terminate()

	repo := infra.NewBeritaAcaraRepository(db)
	handler := &app.UpdateBeritaAcaraCommandHandler{Repo: repo}

	Tahun := "2080"
	FakultasUnit := 1
	Tanggal := "2021-01-01"
	Auditee := 1
	Auditor1 := 2
	Auditor2 := 3

	// Update record
	cmd := app.UpdateBeritaAcaraCommand{
		Uuid:         "14212231-792f-4935-bb1c-9a38695a4b6b",
		Tahun:        Tahun,
		FakultasUnit: FakultasUnit,
		Tanggal:      Tanggal,
		Auditee:      &Auditee,
		Auditor1:     &Auditor1,
		Auditor2:     &Auditor2,
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "14212231-792f-4935-bb1c-9a38695a4b6b", updatedUUID)
}

// func TestUpdateBeritaAcaraCommand_Fail(t *testing.T) {
// 	db, terminate := setupBeritaAcaraMySQL(t)
// 	defer terminate()

// 	repo := infra.NewBeritaAcaraRepository(db)
// 	handler := &app.UpdateBeritaAcaraCommandHandler{Repo: repo}

// 	// Insert record awal
// 	original := domain.BeritaAcara{
// 		UUID: uuid.New(),
// 		Nama: "Dokumen Edge",
// 	}
// 	err := repo.Create(context.Background(), &original)
// 	assert.NoError(t, err)

// 	uuid := uuid.NewString()
// 	// Update dengan nama yang sama
// 	cmdSame := app.UpdateBeritaAcaraCommand{
// 		Uuid: uuid,
// 		Nama: "Dokumen Edge",
// 	}
// 	_, err = handler.Handle(context.Background(), cmdSame)
// 	assert.Error(t, err)

// 	commonErr, _ := err.(common.Error)

// 	assert.Equal(t, "BeritaAcara.NotFound", commonErr.Code)
// 	assert.Equal(t, fmt.Sprintf("BeritaAcara with identifier %s not found", uuid), commonErr.Description)
// }
