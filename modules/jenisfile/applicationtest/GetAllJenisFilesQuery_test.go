package applicationtest

import (
    "context"
    "testing"
    
    _ "github.com/go-sql-driver/mysql"
    
    app "UnpakSiamida/modules/jenisfile/application/GetAllJenisFiles"
    infra "UnpakSiamida/modules/jenisfile/infrastructure"
    domain "UnpakSiamida/common/domain"
)

func TestGetAllJenisFilesIntegration(t *testing.T) {
    db, cleanup := setupJenisFileMySQL(t)
    defer cleanup()

    // Use GORM wrapper repo constructor
    gormrepo := infra.NewJenisFileRepository(db)
    handler := app.GetAllJenisFilesQueryHandler{Repo: gormrepo}

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
        {"No search, returns all", "", 4}, 
        {"Search matching 'LED Prodi'", "LED Prodi", 1},
        {"Search not matching anything", "TidakAda", 0},
    }

    for _, tt := range searchTests {
        t.Run("Search_"+tt.name, func(t *testing.T) {

            q := app.GetAllJenisFilesQuery{
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
        // nama
        {"nama eq 'LED Prodi'", []domain.SearchFilter{
            {"nama", "eq", str("LED Prodi")},
        }, 1},
       {"nama like 'LED Prodi'", []domain.SearchFilter{
            {"nama", "like", str("LED Prodi")},
        }, 1},
        {"nama neq 'LED Prodi'", []domain.SearchFilter{
            {"nama", "neq", str("LED Prodi")},
        }, 3},
    }

    for _, tt := range filterTests {
        t.Run("Filter_"+tt.name, func(t *testing.T) {

            q := app.GetAllJenisFilesQuery{
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