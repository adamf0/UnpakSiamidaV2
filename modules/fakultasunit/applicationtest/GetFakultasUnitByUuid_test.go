package applicationtest

import (
    "context"
    "testing"
    "strings"

    app "UnpakSiamida/modules/fakultasunit/application/GetFakultasUnit"
    infra "UnpakSiamida/modules/fakultasunit/infrastructure"
)

func TestGetFakultasUnitByUuid_Success(t *testing.T) {
    db, cleanup := setupMySQL(t)
    defer cleanup()

    repo := infra.NewFakultasUnitRepository(db)
    handler := app.GetFakultasUnitByUuidQueryHandler{Repo: repo}

    // UUID fix yang kamu tentukan
    fixedUUID := "e76447c9-097a-4a1f-8c85-066058e0c299"

    q := app.GetFakultasUnitByUuidQuery{Uuid: fixedUUID}

    res, err := handler.Handle(context.Background(), q)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if res.UUID.String() != fixedUUID {
        t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
    }
}

func TestGetFakultasUnitByUuid_Errors(t *testing.T) {
    db, cleanup := setupMySQL(t)
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
            t.Parallel()

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