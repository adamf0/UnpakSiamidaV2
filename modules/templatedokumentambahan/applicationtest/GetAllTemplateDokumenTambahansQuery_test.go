package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	domain "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/templatedokumentambahan/application/GetAllTemplateDokumenTambahans"
	infra "UnpakSiamida/modules/templatedokumentambahan/infrastructure"
)

func TestGetAllTemplateDokumenTambahansIntegration(t *testing.T) {
	db, cleanup := setupTemplateDokumenTambahanMySQL(t)
	defer cleanup()

	// Use GORM wrapper repo constructor
	gormrepo := infra.NewTemplateDokumenTambahanRepository(db)
	handler := app.GetAllTemplateDokumenTambahansQueryHandler{Repo: gormrepo}

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
		{"No search, returns all", "", 10},
		{"Search matching 'fakultas#all'", "fakultas#all", 8},
		{"Search not matching anything", "TidakAda", 0},
	}

	for _, tt := range searchTests {
		t.Run("Search_"+tt.name, func(t *testing.T) {

			q := app.GetAllTemplateDokumenTambahansQuery{
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
		// jenisfile
		{"jenisfile eq 'Program Kerja Sesuai Dengan Template 2024 disertai Monev'", []domain.SearchFilter{
			{"jenisfile", "eq", str("Program Kerja Sesuai Dengan Template 2024 disertai Monev")},
		}, 3},
		{"jenisfile like 'Program Kerja Sesuai Dengan Template 2024 disertai Monev'", []domain.SearchFilter{
			{"jenisfile", "like", str("Program Kerja Sesuai Dengan Template 2024 disertai Monev")},
		}, 3},
		{"jenisfile neq 'Program Kerja Sesuai Dengan Template 2024 disertai Monev'", []domain.SearchFilter{
			{"jenisfile", "neq", str("Program Kerja Sesuai Dengan Template 20210 disertai Monev")},
		}, 10},
	}

	for _, tt := range filterTests {
		t.Run("Filter_"+tt.name, func(t *testing.T) {

			q := app.GetAllTemplateDokumenTambahansQuery{
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
