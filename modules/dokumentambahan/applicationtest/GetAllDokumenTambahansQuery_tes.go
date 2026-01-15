package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	domain "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/dokumentambahan/application/GetAllDokumenTambahans"
	infra "UnpakSiamida/modules/dokumentambahan/infrastructure"
)

func TestGetAllDokumenTambahansIntegration(t *testing.T) {
	db, cleanup := setupDokumenTambahanMySQL(t)
	defer cleanup()

	// Use GORM wrapper repo constructor
	gormrepo := infra.NewDokumenTambahanRepository(db)
	handler := app.GetAllDokumenTambahansQueryHandler{Repo: gormrepo}

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
		{"Search matching 2025", "2025", 10},
		{"Search not matching anything", "TidakAda", 0},
	}

	for _, tt := range searchTests {
		t.Run("Search_"+tt.name, func(t *testing.T) {

			q := app.GetAllDokumenTambahansQuery{
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
		{"dokumen eq Program Kerja Sesuai Dengan Template 2024 disertai Monev", []domain.SearchFilter{
			{"dokumen", "eq", str("Program Kerja Sesuai Dengan Template 2024 disertai Monev")},
		}, 10},
		{"dokumen like Program Kerja Sesuai Dengan Template 2024 disertai Monev", []domain.SearchFilter{
			{"dokumen", "like", str("Program Kerja Sesuai Dengan Template 2024 disertai Monev")},
		}, 10},
		{"dokumen neq Program Kerja Sesuai Dengan Template 2024 disertai Monev", []domain.SearchFilter{
			{"dokumen", "neq", str("Program Kerja Sesuai Dengan Template 2024 disertai Monev")},
		}, 10},

		{"pertanyaan eq Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?", []domain.SearchFilter{
			{"pertanyaan", "eq", str("Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?")},
		}, 10},
		{"pertanyaan like Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?", []domain.SearchFilter{
			{"pertanyaan", "like", str("Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?")},
		}, 10},
		{"pertanyaan neq Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?", []domain.SearchFilter{
			{"pertanyaan", "neq", str("Apakah Sudah Lengkap Sesuai Dengan Template Proker 2024 Beserta Monevnya?")},
		}, 10},

		{"targetaudit eq PUTIK (PUSAT TEKNOLOGI INFORMASI DAN KOMUNIKASI)", []domain.SearchFilter{
			{"targetaudit", "eq", str("PUTIK (PUSAT TEKNOLOGI INFORMASI DAN KOMUNIKASI)")},
		}, 1},
		{"targetaudit like PUTIK (PUSAT TEKNOLOGI INFORMASI DAN KOMUNIKASI)", []domain.SearchFilter{
			{"targetaudit", "like", str("PUTIK (PUSAT TEKNOLOGI INFORMASI DAN KOMUNIKASI)")},
		}, 1},
		{"targetaudit neq PUTIK (PUSAT TEKNOLOGI INFORMASI DAN KOMUNIKASI)", []domain.SearchFilter{
			{"targetaudit", "neq", str("PUTIK (PUSAT TEKNOLOGI INFORMASI DAN KOMUNIKASI)")},
		}, 10},

		{"tipe eq fakultas", []domain.SearchFilter{
			{"tipe", "eq", str("fakultas")},
		}, 10},
		{"tipe like fakultas", []domain.SearchFilter{
			{"tipe", "like", str("fakultas")},
		}, 10},
		{"tipe neq fakultas", []domain.SearchFilter{
			{"tipe", "neq", str("fakultas")},
		}, 10},

		{"tahun eq fakultas", []domain.SearchFilter{
			{"tahun", "eq", str("2025")},
		}, 10},
		{"tahun like 2025", []domain.SearchFilter{
			{"tahun", "like", str("2025")},
		}, 10},
		{"tahun neq 2025", []domain.SearchFilter{
			{"tahun", "neq", str("2025")},
		}, 10},
		{"tahun in 2025", []domain.SearchFilter{
			{"tahun", "in", str("2025,2024")},
		}, 10},
		{"tahun gt 2024", []domain.SearchFilter{
			{"tahun", "gt", str("2024")},
		}, 10},
		{"tahun gte 2024", []domain.SearchFilter{
			{"tahun", "gte", str("2024")},
		}, 10},
		{"tahun lt 2025", []domain.SearchFilter{
			{"tahun", "lt", str("2025")},
		}, 10},
		{"tahun lte 2025", []domain.SearchFilter{
			{"tahun", "lte", str("2025")},
		}, 10},
	}

	for _, tt := range filterTests {
		t.Run("Filter_"+tt.name, func(t *testing.T) {

			q := app.GetAllDokumenTambahansQuery{
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
