package applicationtest

import (
	"context"
	"strings"
	"testing"

	app "UnpakSiamida/modules/fakultasunit/application/GetFakultasUnit"
	infra "UnpakSiamida/modules/fakultasunit/infrastructure"
)

func GetFakultasUnitByUuid_Success(t *testing.T) {
	db, cleanup := setupFakultasUnitMySQL(t)
	defer cleanup()

	repo := infra.NewFakultasUnitRepository(db)
	handler := app.GetFakultasUnitByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu tentukan
	fixedUUID := "0d2fa3f8-6df3-45b8-8985-654cb49d5d03"

	q := app.GetFakultasUnitByUuidQuery{Uuid: fixedUUID}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UUID.String() != fixedUUID {
		t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
	}
}

func GetFakultasUnitByUuid_Errors(t *testing.T) {
	db, cleanup := setupFakultasUnitMySQL(t)
	defer cleanup()

	repo := infra.NewFakultasUnitRepository(db)
	handler := app.GetFakultasUnitByUuidQueryHandler{Repo: repo}

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
			q := app.GetFakultasUnitByUuidQuery{Uuid: tt.uuid}

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
