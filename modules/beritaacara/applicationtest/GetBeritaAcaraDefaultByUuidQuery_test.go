package applicationtest

import (
	"context"
	"strings"
	"testing"

	app "UnpakSiamida/modules/beritaacara/application/GetBeritaAcaraDefault"
	domain "UnpakSiamida/modules/beritaacara/domain"
	infra "UnpakSiamida/modules/beritaacara/infrastructure"
)

// ------------------------------
// SUCCESS
// ------------------------------
func TestGetBeritaAcaraDefaultByUuid_Success(t *testing.T) {
	db, cleanup := setupBeritaAcaraMySQL(t)
	defer cleanup()

	repo := infra.NewBeritaAcaraRepository(db)
	handler := app.GetBeritaAcaraDefaultByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu seed di setupBeritaAcaraMySQL
	fixedUUID := "14212231-792f-4935-bb1c-9a38695a4b6b"

	q := app.GetBeritaAcaraDefaultByUuidQuery{Uuid: fixedUUID}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UUID.String() != fixedUUID {
		t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
	}
}

// ------------------------------
// ERROR CASES
// ------------------------------
func TestGetBeritaAcaraDefaultByUuid_Errors(t *testing.T) {
	db, cleanup := setupBeritaAcaraMySQL(t)
	defer cleanup()

	repo := infra.NewBeritaAcaraRepository(db)
	handler := app.GetBeritaAcaraDefaultByUuidQueryHandler{Repo: repo}

	tests := []struct {
		name   string
		uuid   string
		expect error
	}{
		{
			name:   "Invalid UUID format",
			uuid:   "not-a-valid-uuid",
			expect: domain.NotFound("not-a-valid-uuid"),
		},
		{
			name:   "UUID valid but not in DB",
			uuid:   "11111111-1111-1111-1111-111111111111",
			expect: domain.NotFound("11111111-1111-1111-1111-111111111111"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := app.GetBeritaAcaraDefaultByUuidQuery{Uuid: tt.uuid}

			_, err := handler.Handle(context.Background(), q)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.expect.Error())) {
				t.Fatalf("expected error containing %q, got %v", tt.expect.Error(), err)
			}
		})
	}
}
