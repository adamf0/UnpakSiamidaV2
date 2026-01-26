package applicationtest

import (
	"context"
	"strings"
	"testing"

	app "UnpakSiamida/modules/tahunproker/application/GetTahunProker"
	infra "UnpakSiamida/modules/tahunproker/infrastructure"
)

func TestGetTahunProkerByUuid_Success(t *testing.T) {
	db, cleanup := setupTahunProkerMySQL(t)
	defer cleanup()

	repo := infra.NewTahunProkerRepository(db)
	handler := app.GetTahunProkerByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu tentukan
	fixedUUID := "666a6b72-d2b4-481f-adb8-298d807e9e20"

	q := app.GetTahunProkerByUuidQuery{Uuid: fixedUUID}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UUID.String() != fixedUUID {
		t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
	}
}

func TestGetTahunProkerByUuid_Errors(t *testing.T) {
	db, cleanup := setupTahunProkerMySQL(t)
	defer cleanup()

	repo := infra.NewTahunProkerRepository(db)
	handler := app.GetTahunProkerByUuidQueryHandler{Repo: repo}

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
			q := app.GetTahunProkerByUuidQuery{Uuid: tt.uuid}

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
