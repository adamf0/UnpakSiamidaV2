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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	validFakultasUUID = "dea9a83f-70b3-4295-85ed-459eb1a9f6a0"
	validAuditeeUUID  = "c7fd1d83-2d34-42a7-9cfe-38fa5f813188"
	validAuditor1UUID = "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e"
	validAuditor2UUID = "63b1c4b2-5e13-407f-a9fc-a8c775d9ecaa"

	validTahun = "2032"
	testLayout = "2006-01-02"
)

// ============================
// Helpers
// ============================

func datePlus(days int) string {
	return time.Now().Add(time.Duration(days) * 24 * time.Hour).Format(testLayout)
}

func buildValidCommand() app.CreateRenstraCommand {
	return app.CreateRenstraCommand{
		FakultasUnit:                  validFakultasUUID,
		Auditee:                       validAuditeeUUID,
		Auditor1:                      validAuditor1UUID,
		Auditor2:                      validAuditor2UUID,
		Tahun:                         validTahun,
		PeriodeUploadMulai:            datePlus(0),
		PeriodeUploadAkhir:            datePlus(1),
		PeriodeAssesmentDokumenMulai:  datePlus(2),
		PeriodeAssesmentDokumenAkhir:  datePlus(3),
		PeriodeAssesmentLapanganMulai: datePlus(4),
		PeriodeAssesmentLapanganAkhir: datePlus(5),
	}
}

// ============================
// Validation Tests
// ============================

func TestCreateRenstraCommandValidation_Success(t *testing.T) {
	cmd := buildValidCommand()
	err := app.CreateRenstraCommandValidation(cmd)
	assert.NoError(t, err)
}

func TestCreateRenstraCommandValidation_Fail(t *testing.T) {
	cmd := app.CreateRenstraCommand{}

	err := app.CreateRenstraCommandValidation(cmd)
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "Fakultas Unit cannot be blank")
	assert.Contains(t, err.Error(), "Periode Upload Mulai cannot be blank")
	assert.Contains(t, err.Error(), "Periode Upload Akhir cannot be blank")
	assert.Contains(t, err.Error(), "Auditee cannot be blank")
	assert.Contains(t, err.Error(), "Auditor1 cannot be blank")
	assert.Contains(t, err.Error(), "Auditor2 cannot be blank")
}

// ============================
// Handler Success
// ============================

func TestCreateRenstraCommandHandler_Success(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	handler := &app.CreateRenstraCommandHandler{
		Repo:             infra.NewRenstraRepository(db),
		FakultasUnitRepo: infraFakultas.NewFakultasUnitRepository(db),
		UserRepo:         infraUser.NewUserRepository(db),
	}

	cmd := buildValidCommand()

	_, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

// ============================
// Handler Fail
// ============================

func TestCreateRenstraCommandHandler_Fail(t *testing.T) {
	db, terminate := setupRenstraMySQL(t)
	defer terminate()

	handler := &app.CreateRenstraCommandHandler{
		Repo:             infra.NewRenstraRepository(db),
		FakultasUnitRepo: infraFakultas.NewFakultasUnitRepository(db),
		UserRepo:         infraUser.NewUserRepository(db),
	}

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
