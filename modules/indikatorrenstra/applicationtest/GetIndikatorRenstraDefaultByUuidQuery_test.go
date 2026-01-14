package applicationtest

import (
	"context"
	"strings"
	"testing"

	app "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstraDefault"
	domain "UnpakSiamida/modules/indikatorrenstra/domain"
	infra "UnpakSiamida/modules/indikatorrenstra/infrastructure"
)

// ------------------------------
// SUCCESS
// ------------------------------
func TestGetIndikatorRenstraDefaultByUuid_Success(t *testing.T) {
	db, cleanup := setupIndikatorRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewIndikatorRenstraRepository(db)
	handler := app.GetIndikatorRenstraDefaultByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu seed di setupIndikatorRenstraMySQL
	fixedUUID := "d0754056-4d55-4091-a13d-0fae624a7616"

	q := app.GetIndikatorRenstraDefaultByUuidQuery{Uuid: fixedUUID}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Uuid.String() != fixedUUID {
		t.Fatalf("expected UUID %s, got %s", fixedUUID, res.Uuid)
	}
}

// ------------------------------
// ERROR CASES
// ------------------------------
func TestGetIndikatorRenstraDefaultByUuid_Errors(t *testing.T) {
	db, cleanup := setupIndikatorRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewIndikatorRenstraRepository(db)
	handler := app.GetIndikatorRenstraDefaultByUuidQueryHandler{Repo: repo}

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
			q := app.GetIndikatorRenstraDefaultByUuidQuery{Uuid: tt.uuid}

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
