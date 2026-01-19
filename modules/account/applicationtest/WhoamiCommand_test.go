package applicationtest

import (
	"context"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/account/application/Whoami"
	infra "UnpakSiamida/modules/account/infrastructure"

	"github.com/stretchr/testify/assert"
)

func TestWhoami_Success(t *testing.T) {
	db, cleanup := setupAccountMySQL(t)
	defer cleanup()

	repo := infra.NewAccountRepository(db)
	handler := app.WhoamiCommandHandler{Repo: repo}

	cmd := app.WhoamiCommand{
		SID: "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
	}

	// validation
	err := app.WhoamiCommandValidation(cmd)
	assert.NoError(t, err)

	user, err := handler.Handle(context.Background(), cmd)
	assert.NotNil(t, user)
	// if err != nil {
	// 	t.Fatalf("expected success, got error: %v", err)
	// }

	// if user.UUID == "" {
	// 	t.Fatalf("uuid should not be empty")
	// }

	// if user.ExtraRole == nil {
	// 	t.Fatalf("ExtraRole should be initialized (empty slice)")
	// }
}

func TestWhoamiIntegration_Failed_InvalidUUID(t *testing.T) {
	_, cleanup := setupAccountMySQL(t)
	defer cleanup()

	cmd := app.WhoamiCommand{
		SID: "not-a-uuid",
	}

	err := app.WhoamiCommandValidation(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UUID must be a valid UUIDv4 format")
}

func TestWhoamiIntegration_Failed_UserNotFound(t *testing.T) {
	db, cleanup := setupAccountMySQL(t)
	defer cleanup()

	_ = db.Exec("TRUNCATE TABLE users")

	repo := infra.NewAccountRepository(db)
	handler := app.WhoamiCommandHandler{Repo: repo}

	cmd := app.WhoamiCommand{
		SID: "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
	}

	err := app.WhoamiCommandValidation(cmd)
	assert.NoError(t, err)

	_, err = handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, commonErr.Code, "Account.InvalidCredential")
	assert.Equal(t, commonErr.Description, "invalid credentials")
}

func TestWhoamiIntegration_ValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		cmd  app.WhoamiCommand
	}{
		{"Empty SID", app.WhoamiCommand{SID: ""}},
		{"XSS SID", app.WhoamiCommand{SID: "<script>alert(1)</script>"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.WhoamiCommandValidation(tt.cmd)
			if err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

func TestWhoamiIntegration_ContextTimeout(t *testing.T) {
	db, cleanup := setupAccountMySQL(t)
	defer cleanup()

	repo := infra.NewAccountRepository(db)
	handler := app.WhoamiCommandHandler{Repo: repo}

	cmd := app.WhoamiCommand{
		SID: "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
	}

	// validation
	err := app.WhoamiCommandValidation(cmd)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // force cancel

	_, err = handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.True(t, err == context.Canceled || err == context.DeadlineExceeded, "expected context canceled or timeout error")
}
