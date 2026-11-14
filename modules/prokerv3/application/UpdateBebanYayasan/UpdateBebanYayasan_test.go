package UpdateBebanYayasan

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"testing"

	"github.com/google/uuid"
)

type mockBebanRepo struct {
	domain.IProkerV3Repository
	proker *domain.ProkerV3
}

func (m *mockBebanRepo) GetByUuid(ctx context.Context, uid uuid.UUID) (*domain.ProkerV3, error) {
	if m.proker == nil || m.proker.UUID != uid {
		return nil, domain.ProkerNotFound()
	}
	return m.proker, nil
}

func (m *mockBebanRepo) UpdateBebanYayasan(ctx context.Context, uid uuid.UUID, beban bool) (*domain.ProkerV3, error) {
	m.proker.BebanYayasan = beban
	return m.proker, nil
}

func TestUpdateBebanYayasanCommandHandler(t *testing.T) {
	uid := uuid.New()
	p := &domain.ProkerV3{
		UUID:         uid,
		BebanYayasan: false,
	}
	repo := &mockBebanRepo{proker: p}
	h := &UpdateBebanYayasanCommandHandler{Repo: repo}

	// Test toggle
	res, err := h.Handle(context.Background(), UpdateBebanYayasanCommand{
		UUID: uid,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.BebanYayasan {
		t.Fatalf("expected BebanYayasan to be true, got %v", res.BebanYayasan)
	}

	// Test explicit set
	f := false
	res2, err := h.Handle(context.Background(), UpdateBebanYayasanCommand{
		UUID:         uid,
		BebanYayasan: &f,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res2.BebanYayasan {
		t.Fatalf("expected BebanYayasan to be false, got %v", res2.BebanYayasan)
	}
}
