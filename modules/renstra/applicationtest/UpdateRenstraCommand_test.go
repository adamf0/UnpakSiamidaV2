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
	infra "UnpakSiamida/modules/renstra/infrastructure"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Contoh validasi sederhana command
func TestUpdateRenstraCommandValidation_Success(t *testing.T) {
	cmd := app.UpdateRenstraCommand{
		Uuid:                          "c67a37c3-7f25-43de-835d-e4bece0eb308",
		FakultasUnit:                  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Auditee:                       "c7fd1d83-2d34-42a7-9cfe-38fa5f813188",
		Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
		Auditor2:                      "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa",
		Tahun:                         "2031",
		PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
		PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).Format("2006-01-02"),
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
		FakultasUnit:                  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
		Auditee:                       "c7fd1d83-2d34-42a7-9cfe-38fa5f813188",
		Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
		Auditor2:                      "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa",
		Tahun:                         "2031",
		PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
		PeriodeUploadAkhir:            time.Now().Add(1 * 24 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentDokumenMulai:  time.Now().Add(2 * 24 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentDokumenAkhir:  time.Now().Add(3 * 24 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentLapanganMulai: time.Now().Add(4 * 24 * time.Hour).Format("2006-01-02"),
		PeriodeAssesmentLapanganAkhir: time.Now().Add(5 * 24 * time.Hour).Format("2006-01-02"),
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
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

	invalidUuid := uuid.NewString()

	tests := []struct {
		name       string
		cmd        app.UpdateRenstraCommand
		expectCode string
		expectDesc string
	}{
		{
			"NotFound",
			app.UpdateRenstraCommand{
				Uuid:                          invalidUuid,
				FakultasUnit:                  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
				Auditee:                       "c7fd1d83-2d34-42a7-9cfe-38fa5f813188",
				Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
				Auditor2:                      "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa",
				Tahun:                         "2031",
				PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
				PeriodeUploadAkhir:            time.Now().Add(1 * 24 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenMulai:  time.Now().Add(2 * 24 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenAkhir:  time.Now().Add(3 * 24 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganMulai: time.Now().Add(4 * 24 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganAkhir: time.Now().Add(5 * 24 * time.Hour).Format("2006-01-02"),
			},
			"Renstra.NotFound",
			fmt.Sprintf("Renstra with identifier %s not found", invalidUuid),
		},
		{
			"InvalidFakultasUnit",
			app.UpdateRenstraCommand{
				Uuid:                          "c67a37c3-7f25-43de-835d-e4bece0eb308",
				FakultasUnit:                  "dea9a83f-70b3-4295-85ed-000000000000",
				Auditee:                       "c7fd1d83-2d34-42a7-9cfe-38fa5f813188",
				Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
				Auditor2:                      "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa",
				Tahun:                         "2031",
				PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
				PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).Format("2006-01-02"),
			},
			"Renstra.InvalidFakultasUnit",
			"fakultas unit is invalid",
		},
		{
			"MissingAuditee",
			app.UpdateRenstraCommand{
				Uuid:                          "c67a37c3-7f25-43de-835d-e4bece0eb308",
				FakultasUnit:                  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
				Auditee:                       "c7fd1d83-2d34-42a7-9cfe-000000000000",
				Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
				Auditor2:                      "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa",
				Tahun:                         "2031",
				PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
				PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).Format("2006-01-02"),
			},
			"Renstra.MissingAuditee",
			"auditee have not been assigned",
		},
		{
			"MissingAuditor1",
			app.UpdateRenstraCommand{
				Uuid:                          "c67a37c3-7f25-43de-835d-e4bece0eb308",
				FakultasUnit:                  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
				Auditee:                       "c7fd1d83-2d34-42a7-9cfe-38fa5f813188",
				Auditor1:                      "56ce6c95-e23f-463b-bcf6-000000000000",
				Auditor2:                      "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa",
				Tahun:                         "2031",
				PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
				PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).Format("2006-01-02"),
			},
			"Renstra.MissingAuditor1",
			"auditor1 have not been assigned",
		},
		{
			"MissingAuditor2",
			app.UpdateRenstraCommand{
				Uuid:                          "c67a37c3-7f25-43de-835d-e4bece0eb308",
				FakultasUnit:                  "dea9a83f-70b3-4295-85ed-459eb1a9f6a0",
				Auditee:                       "c7fd1d83-2d34-42a7-9cfe-38fa5f813188",
				Auditor1:                      "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
				Auditor2:                      "63b1c4b2-5e13-407f-a9fc-000000000000",
				Tahun:                         "2031",
				PeriodeUploadMulai:            time.Now().Format("2006-01-02"),
				PeriodeUploadAkhir:            time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenMulai:  time.Now().Add(25 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentDokumenAkhir:  time.Now().Add(27 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganMulai: time.Now().Add(28 * time.Hour).Format("2006-01-02"),
				PeriodeAssesmentLapanganAkhir: time.Now().Add(30 * time.Hour).Format("2006-01-02"),
			},
			"Renstra.MissingAuditor2",
			"auditor2 have not been assigned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Handle(context.Background(), tt.cmd)
			assert.Error(t, err)
			commonErr, ok := err.(common.Error)
			assert.True(t, ok)
			assert.Contains(t, commonErr.Code, tt.expectCode)
			assert.Contains(t, commonErr.Description, tt.expectDesc)
		})
	}
}
