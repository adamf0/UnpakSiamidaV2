package applicationtest

import (
    "context"
    "strings"
    "testing"

    app "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstraDefault"
    infra "UnpakSiamida/modules/indikatorrenstra/infrastructure"
)

// ------------------------------
// SUCCESS
// ------------------------------
func TestGetIndikatorRenstraDefaultByUuid_Success(t *testing.T) {
    db, cleanup := setupMySQL(t)
    defer cleanup()

    repo := infra.NewIndikatorRenstraRepository(db)
    handler := app.GetIndikatorRenstraDefaultByUuidQueryHandler{Repo: repo}

    // UUID fix yang kamu seed di setupMySQL
    fixedUUID := "186f2427-8bdd-42d9-a757-65808f364eeb"

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
    db, cleanup := setupMySQL(t)
    defer cleanup()

    repo := infra.NewIndikatorRenstraRepository(db)
    handler := app.GetIndikatorRenstraDefaultByUuidQueryHandler{Repo: repo}

    tests := []struct {
        name   string
        uuid   string
        expect string
    }{
        {
            name:   "Invalid UUID format",
            uuid:   "not-a-valid-uuid",
            expect: "invalid", // dari domain.NotFound()
        },
        {
            name:   "UUID valid tapi tidak ada di database",
            uuid:   "11111111-1111-1111-1111-111111111111",
            expect: "not found",
        },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            q := app.GetIndikatorRenstraDefaultByUuidQuery{Uuid: tt.uuid}

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
