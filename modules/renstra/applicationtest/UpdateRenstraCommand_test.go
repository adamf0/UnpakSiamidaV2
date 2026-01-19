package applicationtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	// common "UnpakSiamida/common/domain"
	common "UnpakSiamida/common/domain"
	infraFakultas "UnpakSiamida/modules/fakultasunit/infrastructure"
	app "UnpakSiamida/modules/renstra/application/UpdateRenstra"
	domain "UnpakSiamida/modules/renstra/domain"
	infra "UnpakSiamida/modules/renstra/infrastructure"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Contoh validasi sederhana command
func TestUpdateRenstraCommandValidation_Success(t *testing.T) {
	cmd := app.UpdateRenstraCommand{
		Uuid:                          "c67a37c3-7f25-43de-835d-e4bece0eb308",
		FakultasUnit:                  "0d2fa3f8-6df3-45b8-8985-654cb49d5d03",
		Auditee:                       "495fe283-3e42-4323-a172-c110036b0c60",
		Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
		Auditor2:                      "b15f4b66-c696-40a4-a047-c51d2be63d4b",
		Tahun:                         "2025",
		PeriodeUploadMulai:            time.Now().String(),
		PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).String(),
		PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).String(),
		PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).String(),
		PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).String(),
		PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).String(),
	}

	err := app.UpdateRenstraCommandValidation(cmd)
	assert.NoError(t, err)
}

func TestUpdateRenstraCommandValidation_Fail(t *testing.T) {
	cmd := app.UpdateRenstraCommand{
		Uuid:                          "",
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
	err := app.UpdateRenstraCommandValidation(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
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
func TestUpdateRenstraCommandHandler_Success(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	repo := infra.NewRenstraRepository(db)
	repoFakultas := infraFakultas.NewFakultasUnitRepository(db)
	repoUser := infraUser.NewUserRepository(db)

	handler := &app.UpdateRenstraCommandHandler{
		Repo:             repo,
		FakultasUnitRepo: repoFakultas,
		UserRepo:         repoUser,
	}

	cmd := app.UpdateRenstraCommand{
		Uuid:                          "c67a37c3-7f25-43de-835d-e4bece0eb308",
		FakultasUnit:                  "0d2fa3f8-6df3-45b8-8985-654cb49d5d03",
		Auditee:                       "495fe283-3e42-4323-a172-c110036b0c60",
		Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
		Auditor2:                      "b15f4b66-c696-40a4-a047-c51d2be63d4b",
		Tahun:                         "2025",
		PeriodeUploadMulai:            time.Now().String(),
		PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).String(),
		PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).String(),
		PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).String(),
		PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).String(),
		PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).String(),
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)

	var saved domain.Renstra
	err = db.Where("uuid = ?", updatedUUID).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, cmd.Tahun, saved.Tahun)
}

// Test handler gagal karena UUID invalid / tidak ada data
func TestUpdateRenstraCommandHandler_Fail(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	repo := infra.NewRenstraRepository(db)
	repoFakultas := infraFakultas.NewFakultasUnitRepository(db)
	repoUser := infraUser.NewUserRepository(db)

	handler := &app.UpdateRenstraCommandHandler{
		Repo:             repo,
		FakultasUnitRepo: repoFakultas,
		UserRepo:         repoUser,
	}

	uuid := uuid.NewString()

	// UUID renstra tidak ditemukan
	cmd := app.UpdateRenstraCommand{
		Uuid:                          uuid,
		FakultasUnit:                  "0d2fa3f8-6df3-45b8-8985-654cb49d5d03",
		Auditee:                       "495fe283-3e42-4323-a172-c110036b0c60",
		Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
		Auditor2:                      "b15f4b66-c696-40a4-a047-c51d2be63d4b",
		Tahun:                         "2025",
		PeriodeUploadMulai:            time.Now().String(),
		PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).String(),
		PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).String(),
		PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).String(),
		PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).String(),
		PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).String(),
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Contains(t, commonErr.Code, "Renstra.NotFound")
	assert.Contains(t, commonErr.Description, fmt.Sprintf("Renstra with identifier %s not found", uuid))
}
