package CreateMataProgramV2

import (
	commondomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"errors"
	"testing"
)

type mockRepo struct {
	domain.IMataProgramV2Repository
	tahunID uint
	units   map[string]uint
	saved   []uint
}

func (m *mockRepo) ResolveOrCreateTahunID(ctx context.Context, tahunStr string) (uint, error) {
	return m.tahunID, nil
}

func (m *mockRepo) ResolveFakultasUnitID(ctx context.Context, rawTarget string) (uint, error) {
	return m.units[rawTarget], nil
}

func (m *mockRepo) UpsertOrRestoreAssignment(ctx context.Context, tahunID uint, unitID uint, mataProgram string) (string, bool, error) {
	m.saved = append(m.saved, unitID)
	return "uuid-first", true, nil
}

func TestCreateMataProgramV2CommandHandler(t *testing.T) {
	repo := &mockRepo{
		tahunID: 7,
		units:   map[string]uint{"10": 10, "25": 25},
	}
	h := &CreateMataProgramV2CommandHandler{Repo: repo}

	t.Run("requires mata_program", func(t *testing.T) {
		_, err := h.Handle(context.Background(), CreateMataProgramV2Command{Targets: []string{"10"}})
		var derr commondomain.Error
		if !errors.As(err, &derr) || derr.Code != "MataProgramV2.MataProgramRequired" {
			t.Fatalf("expected MataProgramRequired, got %v", err)
		}
	})

	t.Run("requires targets", func(t *testing.T) {
		_, err := h.Handle(context.Background(), CreateMataProgramV2Command{MataProgram: "Program A"})
		var derr commondomain.Error
		if !errors.As(err, &derr) || derr.Code != "MataProgramV2.TargetsRequired" {
			t.Fatalf("expected TargetsRequired, got %v", err)
		}
	})

	t.Run("saves valid targets and skips invalid ones", func(t *testing.T) {
		res, err := h.Handle(context.Background(), CreateMataProgramV2Command{
			MataProgram: "Program A",
			Tahun:       "2027",
			Targets:     []string{"10", "0", "invalid", "25"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Count != 2 || res.UUID != "uuid-first" || len(repo.saved) != 2 {
			t.Fatalf("unexpected result: %+v saved=%v", res, repo.saved)
		}
	})
}
