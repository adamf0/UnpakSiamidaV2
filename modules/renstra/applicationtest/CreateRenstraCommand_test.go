package applicationtest

import (
	"context"
	"testing"
	"time"

	common "UnpakSiamida/common/domain"
	infraFakultas "UnpakSiamida/modules/fakultasunit/infrastructure"
	app "UnpakSiamida/modules/renstra/application/CreateRenstra"
	infra "UnpakSiamida/modules/renstra/infrastructure"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/goforj/godump"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Contoh validasi sederhana command
func TestCreateRenstraCommandValidation_Success(t *testing.T) {
	cmd := app.CreateRenstraCommand{
		FakultasUnit:                  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Auditee:                       "c7fd1d83-2d34-42a7-9cfe-38fa5f813188",
		Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
		Auditor2:                      "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa",
		Tahun:                         "2032",
		PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
		PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).Format("2006-01-02"),
	}

	err := app.CreateRenstraCommandValidation(cmd)
	assert.NoError(t, err)
}

func TestCreateRenstraCommandValidation_Fail(t *testing.T) {
	cmd := app.CreateRenstraCommand{
		FakultasUnit:                  "",
		Auditee:                       "",
		Auditor1:                      "",
		Auditor2:                      "",
		Tahun:                         "",
		PeriodeUploadMulai:            "",
		PeriodeUploadAkhir:            "",
		PeriodeAssesmentDokumenMulai:  "",
		PeriodeAssesmentDokumenAkhir:  "",
		PeriodeAssesmentLapanganMulai: "",
		PeriodeAssesmentLapanganAkhir: "",
	}
	err := app.CreateRenstraCommandValidation(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "Fakultas Unit cannot be blank")
	assert.Contains(t, err.Error(), "Periode Upload Mulai cannot be blank")
	assert.Contains(t, err.Error(), "Periode Upload Akhir cannot be blank")
	assert.Contains(t, err.Error(), "Periode Upload Mulai cannot be blank")
	assert.Contains(t, err.Error(), "Periode Upload Akhir cannot be blank")
	assert.Contains(t, err.Error(), "Periode Upload Mulai cannot be blank")
	assert.Contains(t, err.Error(), "Periode Upload Akhir cannot be blank")
	assert.Contains(t, err.Error(), "Auditee cannot be blank")
	assert.Contains(t, err.Error(), "Auditor1 cannot be blank")
	assert.Contains(t, err.Error(), "Auditor2 cannot be blank")
}

// Test handler sukses
func TestCreateRenstraCommandHandler_Success(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	repo := infra.NewRenstraRepository(db)
	repoFakultas := infraFakultas.NewFakultasUnitRepository(db)
	repoUser := infraUser.NewUserRepository(db)

	handler := &app.CreateRenstraCommandHandler{
		Repo:             repo,
		FakultasUnitRepo: repoFakultas,
		UserRepo:         repoUser,
	}

	cmd := app.CreateRenstraCommand{
		FakultasUnit:                  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Auditee:                       "c7fd1d83-2d34-42a7-9cfe-38fa5f813188",
		Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
		Auditor2:                      "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa",
		Tahun:                         "2032",
		PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
		PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).Format("2006-01-02"),
	}

	x, err := handler.Handle(context.Background(), cmd)
	godump.Dump(cmd, x)
	assert.NoError(t, err)
}

// Test handler gagal karena UUID invalid / tidak ada data
func TestCreateRenstraCommandHandler_Fail(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	repo := infra.NewRenstraRepository(db)
	repoFakultas := infraFakultas.NewFakultasUnitRepository(db)
	repoUser := infraUser.NewUserRepository(db)

	handler := &app.CreateRenstraCommandHandler{
		Repo:             repo,
		FakultasUnitRepo: repoFakultas,
		UserRepo:         repoUser,
	}

	// FakultasUnit UUID invalid
	cmd := app.CreateRenstraCommand{
		FakultasUnit: uuid.NewString(),
		Auditee:      uuid.NewString(),
		Auditor1:     uuid.NewString(),
		Auditor2:     uuid.NewString(),
		Tahun:        "2025",
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Contains(t, commonErr.Code, "Renstra.InvalidFakultasUnit")
	assert.Contains(t, commonErr.Description, "fakultas unit is invalid")
}
