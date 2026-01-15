package applicationtest

import (
	"context"
	"testing"

	app "UnpakSiamida/modules/account/application/Whoami"
	domain "UnpakSiamida/modules/account/domain"
	infra "UnpakSiamida/modules/account/infrastructure"
)

func TestWhoamiIntegration_Success(t *testing.T) {
	db, cleanup := setupAccountMySQL(t)
	defer cleanup()

	repo := infra.NewAccountRepository(db)
	handler := app.WhoamiCommandHandler{Repo: repo}

	cmd := app.WhoamiCommand{
		SID: "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
	}

	// validation
	if err := app.WhoamiCommandValidation(cmd); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	user, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if user.UUID == "" {
		t.Fatalf("uuid should not be empty")
	}

	if user.ExtraRole == nil {
		t.Fatalf("ExtraRole should be initialized (empty slice)")
	}
}

func TestWhoamiIntegration_Failed_InvalidUUID(t *testing.T) {
	db, cleanup := setupAccountMySQL(t)
	defer cleanup()

	repo := infra.NewAccountRepository(db)
	handler := app.WhoamiCommandHandler{Repo: repo}

	cmd := app.WhoamiCommand{
		SID: "not-a-uuid",
	}

	err := app.WhoamiCommandValidation(cmd)
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	user, err := handler.Handle(context.Background(), cmd)
	if err == nil {
		t.Fatalf("expected error, got success %+v", user)
	}

	if err.Error() != domain.NotFound(cmd.SID).Error() {
		t.Fatalf("expected NotFound error, got %v", err)
	}
}

func TestWhoamiIntegration_Failed_UserNotFound(t *testing.T) {
	db, cleanup := setupAccountMySQL(t)
	defer cleanup()

	_ = db.Exec("DELETE FROM users")

	repo := infra.NewAccountRepository(db)
	handler := app.WhoamiCommandHandler{Repo: repo}

	cmd := app.WhoamiCommand{
		SID: "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e",
	}

	err := app.WhoamiCommandValidation(cmd)
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	user, err := handler.Handle(context.Background(), cmd)
	if err == nil {
		t.Fatalf("expected error, got success %+v", user)
	}

	if err.Error() != domain.InvalidCredential().Error() {
		t.Fatalf("expected InvalidCredential error, got %v", err)
	}
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
