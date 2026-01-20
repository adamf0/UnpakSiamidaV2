package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	domain "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/templaterenstra/application/GetAllTemplateRenstras"
	infra "UnpakSiamida/modules/templaterenstra/infrastructure"
)

func TestGetAllTemplateRenstrasIntegration(t *testing.T) {
	db, cleanup := setupTemplateRenstraMySQL(t)
	defer cleanup()

	// Use GORM wrapper repo constructor
	gormrepo := infra.NewTemplateRenstraRepository(db)
	handler := app.GetAllTemplateRenstrasQueryHandler{Repo: gormrepo}

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
		{"Search matching '2024'", "2024", 10},
		{"Search not matching anything", "TidakAda", 0},
	}

	for _, tt := range searchTests {
		t.Run("Search_"+tt.name, func(t *testing.T) {

			q := app.GetAllTemplateRenstrasQuery{
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
		// tahun
		{"tahun eq '2024'", []domain.SearchFilter{
			{"tahun", "eq", str("2024")},
		}, 10},
		{"tahun like '2024'", []domain.SearchFilter{
			{"tahun", "like", str("2024")},
		}, 10},
		{"tahun neq '2024'", []domain.SearchFilter{
			{"tahun", "neq", str("2024")},
		}, 10},
	}

	for _, tt := range filterTests {
		t.Run("Filter_"+tt.name, func(t *testing.T) {

			q := app.GetAllTemplateRenstrasQuery{
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
