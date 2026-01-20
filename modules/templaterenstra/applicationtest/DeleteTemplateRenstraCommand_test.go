package applicationtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/templaterenstra/application/DeleteTemplateRenstra"
	infra "UnpakSiamida/modules/templaterenstra/infrastructure"
)

// -----------------------
// SUCCESS VALIDATION
// -----------------------
func TestDeleteTemplateRenstraCommandValidation_Success(t *testing.T) {
	cmd := app.DeleteTemplateRenstraCommand{
		Uuid: "c6df396d-b15e-4129-b1c8-4f312b2830ca",
	}

	err := app.DeleteTemplateRenstraCommandValidation(cmd)
	assert.NoError(t, err)
}

// -----------------------
// FAIL VALIDATION
// -----------------------
func TestDeleteTemplateRenstraCommandValidation_Fail(t *testing.T) {
	cmd := app.DeleteTemplateRenstraCommand{
		Uuid: "",
	}

	err := app.DeleteTemplateRenstraCommandValidation(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID cannot be blank")
}

// -----------------------
// SUCCESS HANDLER
// -----------------------
func TestDeleteTemplateRenstraCommandHandler_Success(t *testing.T) {
	db, cleanup := setupTemplateRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewTemplateRenstraRepository(db)
	handler := &app.DeleteTemplateRenstraCommandHandler{
		Repo: repo,
	}

	uuid := "c6df396d-b15e-4129-b1c8-4f312b2830ca"
	cmd := app.DeleteTemplateRenstraCommand{
		Uuid: uuid,
	}

	res, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, uuid, res)
}

// -----------------------
// FAIL HANDLER (Invalid UUID & NotFound)
// -----------------------
func TestDeleteTemplateRenstraCommandHandler_Fail(t *testing.T) {
	db, cleanup := setupTemplateRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewTemplateRenstraRepository(db)
	handler := &app.DeleteTemplateRenstraCommandHandler{
		Repo: repo,
	}

	uuid := uuid.NewString()
	cmd := app.DeleteTemplateRenstraCommand{
		Uuid: uuid,
	}

	_, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "TemplateRenstra.NotFound", commonErr.Code)
	assert.Equal(t, fmt.Sprintf("TemplateRenstra with identifier %s not found", uuid), commonErr.Description)
}

// -----------------------
// EDGE CASE
// -----------------------
func TestDeleteTemplateRenstraCommandHandler_ContextTimeout(t *testing.T) { //tolong perbaiki
	db, cleanup := setupTemplateRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewTemplateRenstraRepository(db)
	handler := app.DeleteTemplateRenstraCommandHandler{Repo: repo}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // force cancel

	cmd := app.DeleteTemplateRenstraCommand{
		Uuid: "c6df396d-b15e-4129-b1c8-4f312b2830ca",
	}
	_, err := handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.True(t, err == context.Canceled || err == context.DeadlineExceeded, "expected context canceled or timeout error")
}
