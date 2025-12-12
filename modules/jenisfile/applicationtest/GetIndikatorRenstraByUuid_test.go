package applicationtest

import (
    "context"
    "testing"
    "strings"

    app "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstra"
    infra "UnpakSiamida/modules/indikatorrenstra/infrastructure"
)

func TestGetIndikatorRenstraByUuid_Success(t *testing.T) {
    db, cleanup := setupJenisFileMySQL(t)
    defer cleanup()

    repo := infra.NewIndikatorRenstraRepository(db)
    handler := app.GetIndikatorRenstraByUuidQueryHandler{Repo: repo}

    // UUID fix yang kamu tentukan
    fixedUUID := "186f2427-8bdd-42d9-a757-65808f364eeb"

    q := app.GetIndikatorRenstraByUuidQuery{Uuid: fixedUUID}

    res, err := handler.Handle(context.Background(), q)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if res.UUID.String() != fixedUUID {
        t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
    }
}

func TestGetIndikatorRenstraByUuid_Errors(t *testing.T) {
    db, cleanup := setupJenisFileMySQL(t)
    defer cleanup()

    repo := infra.NewIndikatorRenstraRepository(db)
    handler := app.GetIndikatorRenstraByUuidQueryHandler{Repo: repo}

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
            q := app.GetIndikatorRenstraByUuidQuery{Uuid: tt.uuid}

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