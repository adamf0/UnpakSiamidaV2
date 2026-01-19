package applicationtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/templatedokumentambahan/application/DeleteTemplateDokumenTambahan"
	infra "UnpakSiamida/modules/templatedokumentambahan/infrastructure"
)

// -----------------------
// SUCCESS VALIDATION
// -----------------------
func TestDeleteTemplateDokumenTambahanCommandValidation_Success(t *testing.T) {
	cmd := app.DeleteTemplateDokumenTambahanCommand{
		Uuid: "9b354f31-be71-4173-9e26-c319d163660d",
	}

	err := app.DeleteTemplateDokumenTambahanCommandValidation(cmd)
	assert.NoError(t, err)
}

// -----------------------
// FAIL VALIDATION
// -----------------------
func TestDeleteTemplateDokumenTambahanCommandValidation_Fail(t *testing.T) {
	cmd := app.DeleteTemplateDokumenTambahanCommand{
		Uuid: "",
	}

	err := app.DeleteTemplateDokumenTambahanCommandValidation(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

// -----------------------
// SUCCESS HANDLER
// -----------------------
func TestDeleteTemplateDokumenTambahanCommandHandler_Success(t *testing.T) {
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	repo := infra.NewTemplateDokumenTambahanRepository(db)
	handler := &app.DeleteTemplateDokumenTambahanCommandHandler{
		Repo: repo,
	}

	uuid := "9b354f31-be71-4173-9e26-c319d163660d"
	cmd := app.DeleteTemplateDokumenTambahanCommand{
		Uuid: uuid,
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uuid, res)
}

// -----------------------
// FAIL HANDLER (Invalid UUID & NotFound)
// -----------------------
func TestDeleteTemplateDokumenTambahanCommandHandler_Fail(t *testing.T) {
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	repo := infra.NewTemplateDokumenTambahanRepository(db)
	handler := &app.DeleteTemplateDokumenTambahanCommandHandler{
		Repo: repo,
	}

	uuid := uuid.NewString()
	cmd := app.DeleteTemplateDokumenTambahanCommand{
		Uuid: uuid,
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "TemplateDokumenTambahan.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("TemplateDokumenTambahan with identifier %s not found", uuid), commonErr.Description)
}

// -----------------------
// EDGE CASE
// -----------------------
func TestDeleteTemplateDokumenTambahanCommandHandler_ContextTimeout(t *testing.T) { //tolong perbaiki
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	repo := infra.NewTemplateDokumenTambahanRepository(db)
	handler := app.DeleteTemplateDokumenTambahanCommandHandler{Repo: repo}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // force cancel

	cmd := app.DeleteTemplateDokumenTambahanCommand{
		Uuid: "9b354f31-be71-4173-9e26-c319d163660d",
	}
	_, err := handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.True(t, err == context.Canceled || err == context.DeadlineExceeded, "expected context canceled or timeout error")
}
