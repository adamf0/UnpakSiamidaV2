package applicationtest

import (
	"context"
	"strings"
	"testing"

	app "UnpakSiamida/modules/tahunproker/application/GetTahunProkerDefault"
	domain "UnpakSiamida/modules/tahunproker/domain"
	infra "UnpakSiamida/modules/tahunproker/infrastructure"
)

// ------------------------------
// SUCCESS
// ------------------------------
func TestGetTahunProkerDefaultByUuid_Success(t *testing.T) {
	db, cleanup := setupTahunProkerMySQL(t)
	defer cleanup()

	repo := infra.NewTahunProkerRepository(db)
	handler := app.GetTahunProkerDefaultByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu seed di setupTahunProkerMySQL
	fixedUUID := "666a6b72-d2b4-481f-adb8-298d807e9e20"

	q := app.GetTahunProkerDefaultByUuidQuery{Uuid: fixedUUID}

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
func TestGetTahunProkerDefaultByUuid_Errors(t *testing.T) {
	db, cleanup := setupTahunProkerMySQL(t)
	defer cleanup()

	repo := infra.NewTahunProkerRepository(db)
	handler := app.GetTahunProkerDefaultByUuidQueryHandler{Repo: repo}

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
			q := app.GetTahunProkerDefaultByUuidQuery{Uuid: tt.uuid}

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
