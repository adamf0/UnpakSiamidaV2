package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	domain "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/beritaacara/application/GetAllBeritaAcaras"
	infra "UnpakSiamida/modules/beritaacara/infrastructure"
)

func TestGetAllBeritaAcarasIntegration(t *testing.T) {
	db, cleanup := setupBeritaAcaraMySQL(t)
	defer cleanup()

	// Use GORM wrapper repo constructor
	gormrepo := infra.NewBeritaAcaraRepository(db)
	handler := app.GetAllBeritaAcarasQueryHandler{Repo: gormrepo}

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
		{"Search matching 'HUKUM'", "Program", 7},
		{"Search not matching anything", "TidakAda", 0},
	}

	for _, tt := range searchTests {
		t.Run("Search_"+tt.name, func(t *testing.T) {

			q := app.GetAllBeritaAcarasQuery{
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
	// filterTests := []struct {
	// 	name         string
	// 	filter       []domain.SearchFilter
	// 	expectedRows int
	// }{
	// 	// nama
	// 	{"nama eq HUKUM", []domain.SearchFilter{
	// 		{"nama_fak_prod_unit", "eq", str("HUKUM")},
	// 	}, 5},
	// 	{"nama like HUKUM", []domain.SearchFilter{
	// 		{"nama_fak_prod_unit", "like", str("HUKUM")},
	// 	}, 5},
	// 	{"nama neq HUKUM", []domain.SearchFilter{
	// 		{"nama_fak_prod_unit", "neq", str("HUKUM")},
	// 	}, 10},
	// }

	// for _, tt := range filterTests {
	// 	t.Run("Filter_"+tt.name, func(t *testing.T) {

	// 		q := app.GetAllBeritaAcarasQuery{
	// 			Search:        "",
	// 			SearchFilters: tt.filter,
	// 			Page:          &page,
	// 			Limit:         &limit,
	// 		}

	// 		res, err := handler.Handle(context.Background(), q)
	// 		if err != nil {
	// 			t.Fatalf("handler returned error: %v", err)
	// 		}

	// 		if len(res.Data) != tt.expectedRows {
	// 			t.Fatalf("[%s] expected %d rows, got %d",
	// 				tt.name, tt.expectedRows, len(res.Data))
	// 		}
	// 	})
	// }
}

func str(v string) *string {
	return &v
}
