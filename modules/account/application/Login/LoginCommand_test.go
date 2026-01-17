package application_test

import (
	"context"
	"testing"

	app "UnpakSiamida/modules/account/application/Login"
	testutil "UnpakSiamida/modules/account/application/testutil"
	domain "UnpakSiamida/modules/account/domain"
	infra "UnpakSiamida/modules/account/infrastructure"
)

func TestLoginIntegration_Success(t *testing.T) {
	db, cleanup := testutil.SetupAccountMySQL(t)
	defer cleanup()

	repo := infra.NewAccountRepository(db)
	handler := app.LoginCommandHandler{Repo: repo}

	cmd := app.LoginCommand{
		Username: "admin",
		Password: "123",
	}

	// validation
	if err := app.LoginCommandValidation(cmd); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	res, err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Fatalf("token should not be empty")
	}

	if res.UserID == "" {
		t.Fatalf("user id should not be empty")
	}
}

func TestLoginIntegration_Failed(t *testing.T) {
	db, cleanup := testutil.SetupAccountMySQL(t)
	defer cleanup()

	repo := infra.NewAccountRepository(db)
	handler := app.LoginCommandHandler{Repo: repo}

	cmd := app.LoginCommand{
		Username: "admin",
		Password: "SALAH",
	}

	err := app.LoginCommandValidation(cmd)
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	res, err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Fatalf("expected error, got success %+v", res)
	}

	if err.Error() != domain.InvalidCredential().Error() {
		t.Fatalf("expected InvalidCredential error, got %v", err)
	}
}

func TestLoginIntegration_ValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		cmd  app.LoginCommand
	}{
		{"Empty username", app.LoginCommand{Username: "", Password: "123"}},
		{"Empty password", app.LoginCommand{Username: "admin", Password: ""}},
		{"XSS username", app.LoginCommand{Username: "<script>alert(1)</script>", Password: "123"}},
		{"XSS password", app.LoginCommand{Username: "admin", Password: "<img onerror=alert(1)>"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.LoginCommandValidation(tt.cmd)
			if err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

func TestLoginIntegration_ContextTimeout(t *testing.T) {
	db, cleanup := testutil.SetupAccountMySQL(t)
	defer cleanup()

	repo := infra.NewAccountRepository(db)
	handler := app.LoginCommandHandler{Repo: repo}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // force cancel

	_, err := handler.Handle(ctx, app.LoginCommand{
		Username: "admin",
		Password: "123",
	})

	if err == nil {
		t.Fatalf("expected context canceled error")
	}
}
