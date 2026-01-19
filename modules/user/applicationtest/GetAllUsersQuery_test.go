package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	domain "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/user/application/GetAllUsers"
	infra "UnpakSiamida/modules/user/infrastructure"
)

func GetAllUsersIntegration(t *testing.T) {
	db, cleanup := setupUserMySQL(t)
	defer cleanup()

	// Use GORM wrapper repo constructor
	gormrepo := infra.NewUserRepository(db)
	handler := app.GetAllUsersQueryHandler{Repo: gormrepo}

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
		{"No search, returns all", "", 10}, //aslinya 178, 10 kena paging
		{"Search matching 'admin", "admin", 2},
		{"Search not matching anything", "TidakAda", 0},
	}

	for _, tt := range searchTests {
		t.Run("Search_"+tt.name, func(t *testing.T) {

			q := app.GetAllUsersQuery{
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
		{"nidn_username eq 'admin'", []domain.SearchFilter{
			{"nidn_username", "eq", str("admin")},
		}, 10},
		{"nidn_username like 'admin'", []domain.SearchFilter{
			{"nidn_username", "like", str("admin")},
		}, 10},
		{"nidn_username neq 'admin'", []domain.SearchFilter{
			{"nidn_username", "neq", str("admin")},
		}, 10}, //aslinya 177, 10 kena paging

		// level
		{"level eq admin", []domain.SearchFilter{
			{"level", "eq", str("admin")},
		}, 2}, //aslinya 1
		// {"level in", []domain.SearchFilter{
		//     {"level", "in", str("auditor1,auditor2,auditee")}, //tunggu hingga modul renstra selesai
		// }, x}, //aslinya x, 10 kena paging

		// MULTI FILTERS (AND)
		// {"indikator eq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing' AND tahun eq '2025'",
		//     []domain.SearchFilter{
		//         {"indikator", "eq", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
		//         {"tahun", "eq", str("2025")},
		//     },
		//     1,
		// },
		// {"indikator eq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing' AND tahun eq '2025'",
		//     []domain.SearchFilter{
		//         {"indikator", "eq", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
		//         {"tahun", "eq", str("2023")},
		//     },
		//     0,
		// },
	}

	for _, tt := range filterTests {
		t.Run("Filter_"+tt.name, func(t *testing.T) {

			q := app.GetAllUsersQuery{
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
