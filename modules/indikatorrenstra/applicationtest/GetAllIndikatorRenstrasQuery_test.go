package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	commonDomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
	infra "UnpakSiamida/modules/indikatorrenstra/infrastructure"
)

func TestGetAllIndikatorRenstrasIntegration(t *testing.T) {
	db, cleanup := setupIndikatorRenstraMySQL(t)
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
		{"No search, returns all", "", 10}, //aslinya 158, 10 kena paging
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
		{"indikator neq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing'", []domain.SearchFilter{
			{"indikator", "neq", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
		}, 10}, //aslinya 156, 10 kena paging

		// tahun
		{"tahun eq 2024", []domain.SearchFilter{
			{"tahun", "eq", str("2024")},
		}, 10}, //aslinya 80, 10 kena paging
		{"tahun in", []domain.SearchFilter{
			{"tahun", "in", str("2025,2024")},
		}, 10}, //aslinya 158, 10 kena paging

		// MULTI FILTERS (AND)
		{"indikator eq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing' AND tahun eq '2025'",
			[]domain.SearchFilter{
				{"indikator", "eq", str("Lulusan memiliki sertifikat kompetensi atau Bahasa asing")},
				{"tahun", "eq", str("2025")},
			},
			1,
		},
		{"indikator eq 'Lulusan memiliki sertifikat kompetensi atau Bahasa asing' AND tahun eq '2025'",
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

func TestGetIndikatorRenstra_NotFound(t *testing.T) {
	db, cleanup := setupIndikatorRenstraMySQL(t)
	defer cleanup()

	_ = db.Exec("TRUNCATE TABLE master_indikator_renstra")

	repo := infra.NewIndikatorRenstraRepository(db)
	handler := app.GetAllIndikatorRenstrasQueryHandler{Repo: repo}

	_, err := handler.Handle(context.Background(), app.GetAllIndikatorRenstrasQuery{})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	derr, ok := err.(commonDomain.Error)
	if !ok {
		t.Fatalf("expected commonDomain.Error, got %T (%v)", err, err)
	}

	if derr.Code != "IndikatorRenstra.EmptyData" {
		t.Fatalf("expected IndikatorRenstra.EmptyData, got %s", derr.Code)
	}
}

func str(v string) *string {
	return &v
}
