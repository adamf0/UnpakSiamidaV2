package applicationtest

import (
	"context"
	"strings"
	"testing"

	app "UnpakSiamida/modules/jenisfile/application/GetJenisFile"
	infra "UnpakSiamida/modules/jenisfile/infrastructure"
)

func GetJenisFileByUuid_Success(t *testing.T) {
	db, cleanup := setupJenisFileMySQL(t)
	defer cleanup()

	repo := infra.NewJenisFileRepository(db)
	handler := app.GetJenisFileByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu tentukan
	fixedUUID := "14212231-792f-4935-bb1c-9a38695a4b6b"

	q := app.GetJenisFileByUuidQuery{Uuid: fixedUUID}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UUID.String() != fixedUUID {
		t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
	}
}

func GetJenisFileByUuid_Errors(t *testing.T) {
	db, cleanup := setupJenisFileMySQL(t)
	defer cleanup()

	repo := infra.NewJenisFileRepository(db)
	handler := app.GetJenisFileByUuidQueryHandler{Repo: repo}

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
			q := app.GetJenisFileByUuidQuery{Uuid: tt.uuid}

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
