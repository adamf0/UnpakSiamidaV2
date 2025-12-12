package applicationtest

import (
    "context"
    "testing"
    
    _ "github.com/go-sql-driver/mysql"
    
    app "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
    infra "UnpakSiamida/modules/indikatorrenstra/infrastructure"
    domain "UnpakSiamida/common/domain"
)

func TestGetAllIndikatorRenstrasIntegration(t *testing.T) {
    db, cleanup := setupMySQL(t)
    defer cleanup()

    // Use GORM wrapper repo constructor
    gormrepo := infra.NewIndikatorRenstraRepository(db)
    handler := app.GetAllIndikatorRenstrasQueryHandler{Repo: gormrepo}

    page := 1
    limit := 10

    // -------------------------------------------------------
    //  GLOBAL SEARCH THEORY TEST
    // -------------------------------------------------------
    searchTests := []struct {
        name         string
        search       string
        expectedRows int
    }{
        {"No search, returns all", "", 158},
        {"Search matching 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing'", "Lulusan memiliki sertifikat kompetensi atau Bahasa asing", 2},
        {"Search not matching anything", "TidakAda", 0},
    }

    for _, tt := range searchTests {
        t.Run("Search_"+tt.name, func(t *testing.T) {

            q := app.GetAllIndikatorRenstrasQuery{
                Search:        tt.search,
                SearchFilters: []domain.SearchFilter{},
                Page:          &page,
                Limit:         &limit,
            }

            res, err := handler.Handle(context.Background(), q)
            if err != nil {
                t.Fatalf("handler returned error: %v", err)
            }

            if len(res.Data) != tt.expectedRows {
                t.Fatalf("expected %d rows, got %d", tt.expectedRows, len(res.Data))
            }
        })
    }

    // -------------------------------------------------------
    //  SEARCH FILTER THEORY  (All Columns × Operators)
    // -------------------------------------------------------
    filterTests := []struct {
        name         string
        filter       []domain.SearchFilter
        expectedRows int
    }{
        // indikator
        {"indikator eq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing'", []domain.SearchFilter{
            {"indikator", "eq", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
        }, 2},
       {"indikator like 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing'", []domain.SearchFilter{
            {"indikator", "like", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
        }, 2},
        {"indikator neq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing'", []domain.SearchFilter{ //fail
            {"indikator", "neq", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
        }, 156},

        // tahun
        {"tahun eq 2024", []domain.SearchFilter{
            {"tahun", "eq", str("2024")},
        }, 80},
        {"tahun in", []domain.SearchFilter{
            {"tahun", "in", str("2025,2024")},
        }, 158},

        // MULTI FILTERS (AND)
        {"indikator eq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing' AND tahun eq '2025'", //fail
            []domain.SearchFilter{
                {"indikator", "eq", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
                {"tahun", "eq", str("2025")},
            },
            1,
        },
        {"indikator eq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing' AND tahun eq '2025'", //fail
            []domain.SearchFilter{
                {"indikator", "eq", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
                {"tahun", "eq", str("2023")},
            },
            0,
        },
    }

    for _, tt := range filterTests {
        t.Run("Filter_"+tt.name, func(t *testing.T) {

            q := app.GetAllIndikatorRenstrasQuery{
                Search:        "",
                SearchFilters: tt.filter,
                Page:          &page,
                Limit:         &limit,
            }

            res, err := handler.Handle(context.Background(), q)
            if err != nil {
                t.Fatalf("handler returned error: %v", err)
            }

            if len(res.Data) != tt.expectedRows {
                t.Fatalf("[%s] expected %d rows, got %d",
                    tt.name, tt.expectedRows, len(res.Data))
            }
        })
    }
}

func str(v string) *string {
    return &v
}