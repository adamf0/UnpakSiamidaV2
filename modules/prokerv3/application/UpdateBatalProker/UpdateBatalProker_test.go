package UpdateBatalProker

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"testing"

	"github.com/google/uuid"
)

type mockBatalRepo struct {
	domain.IProkerV3Repository
	proker *domain.ProkerV3
}

func (m *mockBatalRepo) GetByUuid(ctx context.Context, uid uuid.UUID) (*domain.ProkerV3, error) {
	if m.proker == nil || m.proker.UUID != uid {
		return nil, domain.ProkerNotFound()
	}
	return m.proker, nil
}

func (m *mockBatalRepo) UpdateBatalProker(ctx context.Context, uid uuid.UUID, isBatal bool, catatan string) (*domain.ProkerV3, error) {
	m.proker.IsBatal = isBatal
	m.proker.CatatanBatal = catatan
	return m.proker, nil
}

func TestUpdateBatalProkerCommandHandler(t *testing.T) {
	uid := uuid.New()
	p := &domain.ProkerV3{
		UUID:    uid,
		IsBatal: false,
	}
	repo := &mockBatalRepo{proker: p}
	h := &UpdateBatalProkerCommandHandler{Repo: repo}

	tr := true
	res, err := h.Handle(context.Background(), UpdateBatalProkerCommand{
		UUID:         uid,
		IsBatal:      &tr,
		CatatanBatal: "Program dibatalkan oleh yayasan",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsBatal {
		t.Fatalf("expected IsBatal to be true, got %v", res.IsBatal)
	}
	if res.CatatanBatal != "Program dibatalkan oleh yayasan" {
		t.Fatalf("expected CatatanBatal to match, got %v", res.CatatanBatal)
	}
}
