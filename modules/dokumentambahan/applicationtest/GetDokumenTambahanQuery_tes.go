package applicationtest

import (
	"context"
	"strings"
	"testing"

	app "UnpakSiamida/modules/dokumentambahan/application/GetDokumenTambahan"
	infra "UnpakSiamida/modules/dokumentambahan/infrastructure"
)

func TestGetDokumenTambahanByUuid_Success(t *testing.T) {
	db, cleanup := setupDokumenTambahanMySQL(t)
	defer cleanup()

	repo := infra.NewDokumenTambahanRepository(db)
	handler := app.GetDokumenTambahanByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu tentukan
	fixedUUID := "c836800f-8c09-4e04-ba16-e0ca027ca571"

	q := app.GetDokumenTambahanByUuidQuery{Uuid: fixedUUID}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UUID.String() != fixedUUID {
		t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
	}
}

func TestGetDokumenTambahanByUuid_Errors(t *testing.T) {
	db, cleanup := setupDokumenTambahanMySQL(t)
	defer cleanup()

	repo := infra.NewDokumenTambahanRepository(db)
	handler := app.GetDokumenTambahanByUuidQueryHandler{Repo: repo}

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
			q := app.GetDokumenTambahanByUuidQuery{Uuid: tt.uuid}

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
