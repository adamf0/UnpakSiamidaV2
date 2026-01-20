package applicationtest

import (
	"context"
	"strings"
	"testing"

	app "UnpakSiamida/modules/standarrenstra/application/GetStandarRenstra"
	infra "UnpakSiamida/modules/standarrenstra/infrastructure"
)

func TestGetStandarRenstraByUuid_Success(t *testing.T) {
	db, cleanup := setupStandarRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewStandarRenstraRepository(db)
	handler := app.GetStandarRenstraByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu tentukan
	fixedUUID := "14212231-792f-4935-bb1c-9a38695a4b6b"

	q := app.GetStandarRenstraByUuidQuery{Uuid: fixedUUID}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UUID.String() != fixedUUID {
		t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
	}
}

func TestGetStandarRenstraByUuid_Errors(t *testing.T) {
	db, cleanup := setupStandarRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewStandarRenstraRepository(db)
	handler := app.GetStandarRenstraByUuidQueryHandler{Repo: repo}

	tests := []struct {
		name   string
		uuid   string
		expect string
	}{
		{
			name:   "Invalid UUID format",
			uuid:   "abc-invalid-uuid",
			expect: "invalid", // parse UUID gagal
		},
		{
			name:   "UUID valid tapi tidak ada di DB",
			uuid:   "11111111-1111-1111-1111-111111111111",
			expect: "not found", // GORM akan return ErrRecordNotFound
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			q := app.GetStandarRenstraByUuidQuery{Uuid: tt.uuid}

			_, err := handler.Handle(context.Background(), q)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !strings.Contains(strings.ToLower(err.Error()), tt.expect) {
				t.Fatalf("expected error containing %q, got %v", tt.expect, err)
			}
		})
	}
}
