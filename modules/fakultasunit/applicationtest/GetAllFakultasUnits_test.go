package applicationtest

import (
    "context"
    "testing"
    
    _ "github.com/go-sql-driver/mysql"
    
    app "UnpakSiamida/modules/fakultasunit/application/GetAllFakultasUnits"
    infra "UnpakSiamida/modules/fakultasunit/infrastructure"
    domain "UnpakSiamida/common/domain"
)

func TestGetAllFakultasUnitsIntegration(t *testing.T) {
    db, cleanup := setupMySQL(t)
    defer cleanup()

    // Use GORM wrapper repo constructor
    gormrepo := infra.NewFakultasUnitRepository(db)
    handler := app.GetAllFakultasUnitsQueryHandler{Repo: gormrepo}

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
        {"No search, returns all", "", 2},
        {"Search matching Teknik", "Teknik", 1},
        {"Search not matching anything", "TidakAda", 0},
    }

    for _, tt := range searchTests {
        t.Run("Search_"+tt.name, func(t *testing.T) {

            q := app.GetAllFakultasUnitsQuery{
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
        // nama_fak_prod_unit
        {"nama eq Teknik", []domain.SearchFilter{
            {"nama_fak_prod_unit", "eq", str("Teknik")},
        }, 1},
        {"nama like Tek", []domain.SearchFilter{
            {"nama_fak_prod_unit", "like", str("Tek")},
        }, 1},
        {"nama neq Ekonomi", []domain.SearchFilter{
            {"nama_fak_prod_unit", "neq", str("Ekonomi")},
        }, 1},

        // fakultas
        {"fakultas eq VOKASI", []domain.SearchFilter{
            {"fakultas", "eq", str("VOKASI")},
        }, 1},
        {"fakultas in", []domain.SearchFilter{
            {"fakultas", "in", str("VOKASI,EKONOMI DAN BISNIS")},
        }, 2},

        // jenjang
        {"jenjang eq S1", []domain.SearchFilter{
            {"jenjang", "eq", str("S1")},
        }, 2},
        {"jenjang like S", []domain.SearchFilter{
            {"jenjang", "like", str("S")},
        }, 2},

        // MULTI FILTERS (AND)
        {"fakultas FT AND jenjang S1",
            []domain.SearchFilter{
                {"fakultas", "eq", str("FT")},
                {"jenjang", "eq", str("S1")},
            },
            2,
        },
        {"fakultas FE AND jenjang D3",
            []domain.SearchFilter{
                {"fakultas", "eq", str("FE")},
                {"jenjang", "eq", str("D3")},
            },
            1,
        },
        {"fakultas VOKASI AND type Unit",
            []domain.SearchFilter{
                {"fakultas", "eq", str("VOKASI")},
                {"type", "eq", str("Unit")},
            },
            1,
        },

        // LIKE COMBINATIONS
        {"nama like 'Tek%'",
            []domain.SearchFilter{
                {"nama_fak_prod_unit", "like", str("Tek")},
            },
            1,
        },

        // MIX LIKE + EQ
        {"nama like 'Sistem%' AND fakultas FT",
            []domain.SearchFilter{
                {"nama_fak_prod_unit", "like", str("Sistem")},
                {"fakultas", "eq", str("FT")},
            },
            1,
        },

        // IN operator
        {"fakultas in (FT,FE)",
            []domain.SearchFilter{
                {"fakultas", "in", str("FT,FE")},
            },
            4,
        },

        // COMPLEX 3 FILTERS
        {"fakultas FE AND jenjang S1 AND type Prodi",
            []domain.SearchFilter{
                {"fakultas", "eq", str("FE")},
                {"jenjang", "eq", str("S1")},
                {"type", "eq", str("Prodi")},
            },
            1,
        },
    }

    for _, tt := range filterTests {
        t.Run("Filter_"+tt.name, func(t *testing.T) {

            q := app.GetAllFakultasUnitsQuery{
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