package UpdateAktivitasV3

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"
	"testing"

	"github.com/google/uuid"
)

type stubProkerV3Repo struct {
	domain.IAktivitasV2Repository
	entity      *domain.AktivitasV2
	savedRAB    string
	savedItems  int
	updateCalls int
}

func (s *stubProkerV3Repo) GetAktivitasEntityByUuid(_ context.Context, _ uuid.UUID) (*domain.AktivitasV2, error) {
	return s.entity, nil
}

func (s *stubProkerV3Repo) UpdateAktivitas(_ context.Context, data *domain.AktivitasV2) error {
	s.updateCalls++
	s.savedRAB = data.RAB
	s.savedItems = len(data.RabItems)
	return nil
}

func strPtr(s string) *string    { return &s }
func f64Ptr(v float64) *float64 { return &v }

func TestUpdateAktivitasV3CatatanRevisiAnggaran(t *testing.T) {
	uid := uuid.New()
	repo := &stubProkerV3Repo{
		entity: &domain.AktivitasV2{
			UUID:     uid,
			Anggaran: 1500000,
			RAB:      `[{"rincian":"Item 1","total":1500000}]`,
			RabItems: []domain.AktivitasV2Rab{{Rincian: "Item 1", Total: 1500000}},
		},
	}
	h := &UpdateAktivitasV3CommandHandler{Repo: repo}

	// 1. Update realisasi + note preserves Anggaran and avoids rewriting RAB rows.
	res, err := h.Handle(context.Background(), UpdateAktivitasV3Command{
		UUID:                  uid.String(),
		Realisasi:             f64Ptr(1000000),
		CatatanRevisiAnggaran: strPtr("Penyesuaian jumlah peserta"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.CatatanRevisiAnggaran != "Penyesuaian jumlah peserta" {
		t.Fatalf("got note %q, want Penyesuaian jumlah peserta", res.CatatanRevisiAnggaran)
	}
	if repo.savedRAB != "" || repo.savedItems != 0 {
		t.Fatalf("expected RAB not passed to UpdateAktivitas, got rab=%q items=%d", repo.savedRAB, repo.savedItems)
	}
	if res.Anggaran != 1500000 || len(res.RabItems) != 1 {
		t.Fatalf("expected response Anggaran/RabItems preserved, got %.0f / %d", res.Anggaran, len(res.RabItems))
	}

	// 2. Omitting CatatanRevisiAnggaran (nil) preserves existing note.
	res, err = h.Handle(context.Background(), UpdateAktivitasV3Command{
		UUID:      uid.String(),
		Aktivitas: "Aktivitas Diperbarui",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.CatatanRevisiAnggaran != "Penyesuaian jumlah peserta" {
		t.Fatalf("note should be preserved when omitted, got %q", res.CatatanRevisiAnggaran)
	}

	// 3. Passing empty string clears existing note.
	res, err = h.Handle(context.Background(), UpdateAktivitasV3Command{
		UUID:                  uid.String(),
		Realisasi:             f64Ptr(1500000),
		CatatanRevisiAnggaran: strPtr("   "),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.CatatanRevisiAnggaran != "" {
		t.Fatalf("expected note cleared, got %q", res.CatatanRevisiAnggaran)
	}
}
