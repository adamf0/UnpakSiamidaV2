package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	domain "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/renstra/application/GetAllRenstras"
	infra "UnpakSiamida/modules/renstra/infrastructure"
)

// TestGetAllRenstrasIntegration menjalankan test integrasi untuk query GetAllRenstrasQueryHandler
func TestGetAllRenstrasIntegration(t *testing.T) {
	db, cleanup := setupRenstraMySQL(t)
	defer cleanup()

	// Use GORM wrapper repo constructor
	repo := infra.NewRenstraRepository(db)
	handler := app.GetAllRenstrasQueryHandler{Repo: repo}

	page := 1
	limit := 10

	// -------------------------------------------------------
	// GLOBAL SEARCH THEORY TEST
	// -------------------------------------------------------
	searchTests := []struct {
		name         string
		search       string
		expectedRows int
	}{
		{"No search, returns all", "", 10}, // sesuaikan expectedRows dengan data test
		{"Search matching '2024'", "2024", 10},
		{"Search not matching anything", "TidakAda", 0},
	}

	for _, tt := range searchTests {
		t.Run("Search_"+tt.name, func(t *testing.T) {
			q := app.GetAllRenstrasQuery{
				Search:        tt.search,
				SearchFilters: []domain.SearchFilter{},
				Page:          &page,
				Limit:         &limit,
				Scope:         "", // bisa "audit" kalau perlu
			}

			res, err := handler.Handle(context.Background(), q)
			assert.NoError(t, err)
			assert.Len(t, res.Data, tt.expectedRows)
		})
	}

	// -------------------------------------------------------
	// SEARCH FILTER THEORY (All Columns × Operators)
	// -------------------------------------------------------
	filterTests := []struct {
		name         string
		filter       []domain.SearchFilter
		expectedRows int
	}{
		{"tahun eq '2024'", []domain.SearchFilter{
			{"tahun", "eq", str("2024")},
		}, 10},
		{"fakultas_unit like 'ILMU HUKUM'", []domain.SearchFilter{
			{"fakultas_unit", "like", str("ILMU HUKUM")},
		}, 4},
		{"tahun neq '2023'", []domain.SearchFilter{
			{"tahun", "neq", str("2023")},
		}, 10},
	}

	for _, tt := range filterTests {
		t.Run("Filter_"+tt.name, func(t *testing.T) {
			q := app.GetAllRenstrasQuery{
				Search:        "",
				SearchFilters: tt.filter,
				Page:          &page,
				Limit:         &limit,
				Scope:         "", // bisa "audit" kalau perlu
			}

			res, err := handler.Handle(context.Background(), q)
			assert.NoError(t, err)
			assert.Len(t, res.Data, tt.expectedRows, "[%s] unexpected row count", tt.name)
		})
	}
}

func str(v string) *string {
	return &v
}
