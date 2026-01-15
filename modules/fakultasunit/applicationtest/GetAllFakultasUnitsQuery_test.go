package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	commonDomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/fakultasunit/application/GetAllFakultasUnits"
	infra "UnpakSiamida/modules/fakultasunit/infrastructure"
)

func TestGetAllFakultasUnitsIntegration(t *testing.T) {
	db, cleanup := setupFakultasUnitMySQL(t)
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
		{"No search, returns all", "", 10},
		{"Search matching Teknik", "Teknik", 7},
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
		}, 9},
		{"nama neq Ekonomi", []domain.SearchFilter{
			{"nama_fak_prod_unit", "neq", str("Ekonomi")},
		}, 10},

		// fakultas
		{"fakultas eq VOKASI", []domain.SearchFilter{
			{"fakultas", "eq", str("VOKASI")},
		}, 7},
		{"fakultas in", []domain.SearchFilter{
			{"fakultas", "in", str("VOKASI,EKONOMI DAN BISNIS")},
		}, 10},

		// jenjang
		{"jenjang eq S1", []domain.SearchFilter{
			{"jenjang", "eq", str("s1")},
		}, 10},
		{"jenjang neq S1", []domain.SearchFilter{
			{"jenjang", "neq", str("s1")},
		}, 10},
		{"jenjang like S", []domain.SearchFilter{
			{"jenjang", "like", str("S")},
		}, 10},

		// MULTI FILTERS
		{"fakultas FT AND jenjang S1",
			[]domain.SearchFilter{
				{"fakultas", "eq", str("TEKNIK")},
				{"jenjang", "eq", str("s1")},
			},
			5,
		},
		{"fakultas EKONOMI DAN BISNIS AND type prodi",
			[]domain.SearchFilter{
				{"fakultas", "eq", str("TEKNIK")},
				{"jenjang", "like", str("s1")},
			},
			5,
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

func TestGetFakultasUnit_NotFound(t *testing.T) {
	db, cleanup := setupFakultasUnitMySQL(t)
	defer cleanup()

	_ = db.Exec("DELETE FROM sijamu_fakultas_unit")

	repo := infra.NewFakultasUnitRepository(db)
	handler := app.GetAllFakultasUnitsQueryHandler{Repo: repo}

	_, err := handler.Handle(context.Background(), app.GetAllFakultasUnitsQuery{})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	derr, ok := err.(commonDomain.Error)
	if !ok {
		t.Fatalf("expected commonDomain.Error, got %T (%v)", err, err)
	}

	if derr.Code != "FakultasUnit.EmptyData" {
		t.Fatalf("expected FakultasUnit.EmptyData, got %s", derr.Code)
	}
}

func str(v string) *string {
	return &v
}
