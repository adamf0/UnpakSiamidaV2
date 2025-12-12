package applicationtest

import (
        "context"
        "testing"

				commonDomain "UnpakSiamida/common/domain"
        appactive "UnpakSiamida/modules/tahunrenstra/application/GetActiveTahunRenstra"
        infra "UnpakSiamida/modules/tahunrenstra/infrastructure"
    )

func TestGetActiveTahunRenstra_Success(t *testing.T) {
        db, cleanup := setupMySQL(t)
        defer cleanup()

        repo := infra.NewTahunRenstraRepository(db)
        handler := appactive.GetActiveTahunRenstraQueryHandler{Repo: repo}

        q := appactive.GetActiveTahunRenstraQuery{}
        res, err := handler.Handle(context.Background(), q)

        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }

        if res == nil {
            t.Fatalf("expected data, got nil")
        }

        if res.Tahun != "2025" {
            t.Fatalf("expected tahun 2025 (active), got %s", res.Tahun)
        }
}

func TestGetActiveTahunRenstra_NotFound(t *testing.T) {
    db, cleanup := setupMySQL(t)
    defer cleanup()

    _ = db.Exec("DELETE FROM v_tahun_renstra")

    repo := infra.NewTahunRenstraRepository(db)
    handler := appactive.GetActiveTahunRenstraQueryHandler{Repo: repo}

    _, err := handler.Handle(context.Background(), appactive.GetActiveTahunRenstraQuery{})

    if err == nil {
        t.Fatalf("expected error, got nil")
    }

    derr, ok := err.(commonDomain.Error)
    if !ok {
        t.Fatalf("expected commonDomain.Error, got %T (%v)", err, err)
    }

    if derr.Code != "TahunRenstra.EmptyData" {
        t.Fatalf("expected TahunRenstra.EmptyData, got %s", derr.Code)
    }
}