package applicationtest

import (
	"context"
	"testing"

	"UnpakSiamida/common/domain"
	appgetall "UnpakSiamida/modules/tahunrenstra/application/GetAllTahunRenstras"
	infra "UnpakSiamida/modules/tahunrenstra/infrastructure"
)

func TestGetAllTahunRenstras_Basic(t *testing.T) {
	db, cleanup := setupTahunRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewTahunRenstraRepository(db)
	handler := appgetall.GetAllTahunRenstrasQueryHandler{Repo: repo}

	page := 1
	limit := 10

	q := appgetall.GetAllTahunRenstrasQuery{
		Search:        "",
		SearchFilters: []domain.SearchFilter{},
		Page:          &page,
		Limit:         &limit,
	}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Data) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(res.Data))
	}

	if res.Total != 4 {
		t.Fatalf("expected total=4, got %d", res.Total)
	}
}

func TestGetAllTahunRenstras_Filter(t *testing.T) {
	db, cleanup := setupTahunRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewTahunRenstraRepository(db)
	handler := appgetall.GetAllTahunRenstrasQueryHandler{Repo: repo}

	page := 1
	limit := 10

	tests := []struct {
		name         string
		filter       []domain.SearchFilter
		expectedRows int
	}{
		{"tahun eq 2025", []domain.SearchFilter{
			{"tahun", "eq", str("2025")},
		}, 1},

		{"status eq active", []domain.SearchFilter{
			{"status", "eq", str("active")},
		}, 1},

		{"status neq active", []domain.SearchFilter{
			{"status", "neq", str("active")},
		}, 3},

		{"multi filter tahun=2024 AND status active",
			[]domain.SearchFilter{
				{"tahun", "eq", str("2024")},
				{"status", "eq", str("active")},
			},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			q := appgetall.GetAllTahunRenstrasQuery{
				Search:        "",
				SearchFilters: tt.filter,
				Page:          &page,
				Limit:         &limit,
			}

			res, err := handler.Handle(context.Background(), q)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(res.Data) != tt.expectedRows {
				t.Fatalf("[%s] expected %d rows, got %d",
					tt.name, tt.expectedRows, len(res.Data))
			}
		})
	}
}

func TestGetAllTahunRenstras_GlobalSearch(t *testing.T) {
	db, cleanup := setupTahunRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewTahunRenstraRepository(db)
	handler := appgetall.GetAllTahunRenstrasQueryHandler{Repo: repo}

	page := 1
	limit := 10

	tests := []struct {
		name         string
		search       string
		expectedRows int
	}{
		{"search empty", "", 4},
		{"search '2025'", "2025", 1},
		{"search 'active'", "active", 1},
		{"search 'non'", "non", 3},
		{"search nothing", "xxxx", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			q := appgetall.GetAllTahunRenstrasQuery{
				Search:        tt.search,
				SearchFilters: []domain.SearchFilter{},
				Page:          &page,
				Limit:         &limit,
			}

			res, err := handler.Handle(context.Background(), q)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
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
