package applicationtest

import (
	"context"
	"testing"

	app "UnpakSiamida/modules/beritaacara/application/UpdateBeritaAcara"
	infra "UnpakSiamida/modules/beritaacara/infrastructure"
	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBeritaAcaraCommandValidation_Success(t *testing.T) {
	Tahun := "2080"
	FakultasUnitUuid := "0d2fa3f8-6df3-45b8-8985-654cb49d5d03"
	Tanggal := "2021-01-01"
	AuditeeUuid := "0a853a5f-0475-4b95-aa55-f9009b165771"
	Auditor1Uuid := "495fe283-3e42-4323-a172-c110036b0c60"
	Auditor2Uuid := "d3d2b976-49c5-4fc8-8a78-a92484a97189"

	validCmd := app.UpdateBeritaAcaraCommand{
		Uuid:             uuid.NewString(),
		Tahun:            Tahun,
		FakultasUnitUuid: FakultasUnitUuid,
		Tanggal:          Tanggal,
		AuditeeUuid:      AuditeeUuid,
		Auditor1Uuid:     &Auditor1Uuid,
		Auditor2Uuid:     &Auditor2Uuid,
	}
	err := app.UpdateBeritaAcaraCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestUpdateBeritaAcaraCommandValidation_Fail(t *testing.T) {
	Uuid := ""
	Tahun := ""
	FakultasUnitUuid := ""
	Tanggal := ""
	AuditeeUuid := ""
	Auditor1Uuid := ""
	Auditor2Uuid := ""

	invalidCmd := app.UpdateBeritaAcaraCommand{
		Uuid:             Uuid,
		Tahun:            Tahun,
		FakultasUnitUuid: FakultasUnitUuid,
		Tanggal:          Tanggal,
		AuditeeUuid:      AuditeeUuid,
		Auditor1Uuid:     &Auditor1Uuid,
		Auditor2Uuid:     &Auditor2Uuid,
	}
	err := app.UpdateBeritaAcaraCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
	assert.Contains(t, err.Error(), "Tahun cannot be blank")
	assert.Contains(t, err.Error(), "FakultasUnit cannot be blank")
	assert.Contains(t, err.Error(), "Auditee cannot be blank")
	assert.Contains(t, err.Error(), "Tanggal cannot be blank")
}

func TestUpdateBeritaAcaraCommand_Success(t *testing.T) {
	db, terminate := setupBeritaAcaraMySQL(t)
	defer terminate()

	repo := infra.NewBeritaAcaraRepository(db)
	handler := &app.UpdateBeritaAcaraCommandHandler{Repo: repo}

	Uuid := "14212231-792f-4935-bb1c-9a38695a4b6b"
	Tahun := "2080"
	FakultasUnitUuid := "0d2fa3f8-6df3-45b8-8985-654cb49d5d03"
	Tanggal := "2021-01-01"
	AuditeeUuid := "0a853a5f-0475-4b95-aa55-f9009b165771"
	Auditor1Uuid := "495fe283-3e42-4323-a172-c110036b0c60"
	Auditor2Uuid := "d3d2b976-49c5-4fc8-8a78-a92484a97189"

	// Update record
	cmd := app.UpdateBeritaAcaraCommand{
		Uuid:             Uuid,
		Tahun:            Tahun,
		FakultasUnitUuid: FakultasUnitUuid,
		Tanggal:          Tanggal,
		AuditeeUuid:      AuditeeUuid,
		Auditor1Uuid:     &Auditor1Uuid,
		Auditor2Uuid:     &Auditor2Uuid,
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, Uuid, updatedUUID)
}

func TestUpdateBeritaAcaraCommand_Fail(t *testing.T) {
	db, terminate := setupBeritaAcaraMySQL(t)
	defer terminate()

	handler := &app.UpdateBeritaAcaraCommandHandler{
		Repo:             infra.NewBeritaAcaraRepository(db),
		RepoFakultasUnit: infraFakultasUnit.NewFakultasUnitRepository(db),
		RepoUser:         infraUser.NewUserRepository(db),
	}

	validUuid := "14212231-792f-4935-bb1c-9a38695a4b6b"
	validTahun := "2080"
	validFakultasUnitUuid := "0d2fa3f8-6df3-45b8-8985-654cb49d5d03"
	validTanggal := "2021-01-01"
	validAuditeeUuid := "0a853a5f-0475-4b95-aa55-f9009b165771"
	validAuditor1Uuid := "495fe283-3e42-4323-a172-c110036b0c60"
	validAuditor2Uuid := "d3d2b976-49c5-4fc8-8a78-a92484a97189"

	tests := []struct {
		name          string
		cmd           app.UpdateBeritaAcaraCommand
		expectedError string
	}{
		{
			name: "Fail - EmptyData",
			cmd: app.UpdateBeritaAcaraCommand{
				Uuid:             uuid.NewString(), // ❌ tidak ada
				Tahun:            validTahun,
				FakultasUnitUuid: validFakultasUnitUuid,
				Tanggal:          validTanggal,
				AuditeeUuid:      validAuditeeUuid,
				Auditor1Uuid:     &validAuditor1Uuid,
				Auditor2Uuid:     &validAuditor2Uuid,
			},
			expectedError: "EmptyData",
		},
		{
			name: "Fail - InvalidTanggal",
			cmd: app.UpdateBeritaAcaraCommand{
				Uuid:             validUuid,
				Tahun:            validTahun,
				FakultasUnitUuid: validFakultasUnitUuid,
				Tanggal:          "2021-02-32", // ❌
				AuditeeUuid:      validAuditeeUuid,
				Auditor1Uuid:     &validAuditor1Uuid,
				Auditor2Uuid:     &validAuditor2Uuid,
			},
			expectedError: "InvalidTanggal",
		},
		{
			name: "Fail - NotFoundFakultas",
			cmd: app.UpdateBeritaAcaraCommand{
				Uuid:             validUuid,
				Tahun:            validTahun,
				FakultasUnitUuid: "0d2fa3f8-6df3-45b8-8985-654cb49d5d00", // ❌
				Tanggal:          validTanggal,
				AuditeeUuid:      validAuditeeUuid,
				Auditor1Uuid:     &validAuditor1Uuid,
				Auditor2Uuid:     &validAuditor2Uuid,
			},
			expectedError: "NotFoundFakultas",
		},
		{
			name: "Fail - NotFoundAuditee",
			cmd: app.UpdateBeritaAcaraCommand{
				Uuid:             validUuid,
				Tahun:            validTahun,
				FakultasUnitUuid: validFakultasUnitUuid,
				Tanggal:          validTanggal,
				AuditeeUuid:      uuid.NewString(), // ❌
				Auditor1Uuid:     &validAuditor1Uuid,
				Auditor2Uuid:     &validAuditor2Uuid,
			},
			expectedError: "NotFoundAuditee",
		},
		{
			name: "Fail - DuplicateAssignment (Auditee == Auditor1)",
			cmd: app.UpdateBeritaAcaraCommand{
				Uuid:             validUuid,
				Tahun:            validTahun,
				FakultasUnitUuid: validFakultasUnitUuid,
				Tanggal:          validTanggal,
				AuditeeUuid:      validAuditeeUuid,
				Auditor1Uuid:     &validAuditeeUuid, // ❌ duplicate
				Auditor2Uuid:     &validAuditor2Uuid,
			},
			expectedError: "DuplicateAssignment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedUUID, err := handler.Handle(context.Background(), tt.cmd)

			assert.Error(t, err)
			assert.Empty(t, updatedUUID)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}
