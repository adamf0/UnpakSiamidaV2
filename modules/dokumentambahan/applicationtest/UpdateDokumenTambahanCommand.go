package applicationtest

import (
	"context"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/dokumentambahan/application/UpdateDokumenTambahan"
	domain "UnpakSiamida/modules/dokumentambahan/domain"
	infra "UnpakSiamida/modules/dokumentambahan/infrastructure"
	infraRenstra "UnpakSiamida/modules/renstra/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDokumenTambahanCommandValidation_Success(t *testing.T) {
	link := "https://drive.google.com/example"
	capaianAuditor := ""

	validCmd := app.UpdateDokumenTambahanCommand{
		Uuid:           uuid.NewString(),
		UuidRenstra:    uuid.NewString(),
		Tahun:          "2025",
		Mode:           "auditee",
		Granted:        "2025#auditor2",
		Link:           &link,
		CapaianAuditor: &capaianAuditor,
	}
	err := app.UpdateDokumenTambahanCommandValidation(validCmd)
	assert.NoError(t, err)
}

func TestUpdateDokumenTambahanCommandValidation_Fail(t *testing.T) {
	link := "https://drive.google.com/example"
	capaianAuditor := ""

	invalidCmd := app.UpdateDokumenTambahanCommand{
		Uuid:           "invalid-uuid",
		UuidRenstra:    "",
		Tahun:          "",
		Mode:           "",
		Granted:        "2025#auditor2",
		Link:           &link,
		CapaianAuditor: &capaianAuditor,
	}
	err := app.UpdateDokumenTambahanCommandValidation(invalidCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID wajib diisi")
	assert.Contains(t, err.Error(), "UUID Renstra wajib diisi")
	assert.Contains(t, err.Error(), "Tahun wajib diisi")
	assert.Contains(t, err.Error(), "Mode harus auditee, auditor1 atau auditor2")
}

func TestUpdateDokumenTambahanCommandHandler_Success(t *testing.T) {
	db, terminate := setupDokumenTambahanMySQL(t)
	defer terminate()

	repo := infra.NewDokumenTambahanRepository(db)
	repoRenstra := infraRenstra.NewRenstraRepository(db)
	handler := &app.UpdateDokumenTambahanCommandHandler{
		Repo:        repo,
		RepoRenstra: repoRenstra,
	}

	link := "https://drive.google.com/example"
	capaianAuditor := ""
	catatanAuditor := "entah lah"

	cmd := app.UpdateDokumenTambahanCommand{
		Uuid:           "c836800f-8c09-4e04-ba16-e0ca027ca571",
		UuidRenstra:    "c67a37c3-7f25-43de-835d-e4bece0eb308",
		Tahun:          "2024",
		Mode:           "auditee",
		Granted:        "2024#auditee",
		Link:           &link,
		CapaianAuditor: &capaianAuditor,
		CatatanAuditor: &catatanAuditor,
	}

	updatedUUID, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "", updatedUUID)

	var saved domain.DokumenTambahan
	err = db.Where("uuid = ?", updatedUUID).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, link, saved.Link)
	assert.Equal(t, capaianAuditor, saved.CatatanAuditor)
	assert.Equal(t, catatanAuditor, saved.CatatanAuditor)

	link = "https://drive.google.com/example"
	capaianAuditor = "0"
	catatanAuditor = "entah lah"

	////
	cmd = app.UpdateDokumenTambahanCommand{
		Uuid:           "c836800f-8c09-4e04-ba16-e0ca027ca571",
		UuidRenstra:    "c67a37c3-7f25-43de-835d-e4bece0eb308",
		Tahun:          "2024",
		Mode:           "auditor",
		Granted:        "2024#auditor",
		Link:           &link,
		CapaianAuditor: &capaianAuditor,
		CatatanAuditor: &catatanAuditor,
	}

	updatedUUID, err = handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, "", updatedUUID)

	err = db.Where("uuid = ?", updatedUUID).First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, link, saved.Link)
	assert.Equal(t, capaianAuditor, saved.CatatanAuditor)
	assert.Equal(t, catatanAuditor, saved.CatatanAuditor)
}

func TestUpdateDokumenTambahanCommandHandler_Fail(t *testing.T) {
	db, terminate := setupDokumenTambahanMySQL(t)
	defer terminate()

	repo := infra.NewDokumenTambahanRepository(db)
	repoRenstra := infraRenstra.NewRenstraRepository(db)
	handler := &app.UpdateDokumenTambahanCommandHandler{
		Repo:        repo,
		RepoRenstra: repoRenstra,
	}

	link := "https://drive.google.com/example"
	capaianAuditor := ""

	cmd := app.UpdateDokumenTambahanCommand{
		Uuid:           "c836800f-8c09-4e04-ba16-e0ca027ca571",
		UuidRenstra:    "c67a37c3-7f25-43de-835d-e4bece0eb308",
		Tahun:          "2024",
		Mode:           "auditee",
		Granted:        "2024#auditor2",
		Link:           &link,
		CapaianAuditor: &capaianAuditor,
	}
	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "DokumenTambahan.NotGranted", commonErr.Code)
	assert.Contains(t, "you are not granted permission in this action", commonErr.Description)
}
