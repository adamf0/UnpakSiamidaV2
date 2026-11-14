package GetProkerBatalSetting

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"testing"
	"time"
)

type stubGetBatalRepo struct {
	domain.IProkerV3Repository
	statusCalled bool
	targetUnit   *int
}

func (s *stubGetBatalRepo) GetProkerBatalSettingStatus(_ context.Context, _ time.Time, targetUnitID *int) (*domain.ProkerBatalSettingStatus, error) {
	s.statusCalled = true
	s.targetUnit = targetUnitID
	return &domain.ProkerBatalSettingStatus{
		SettingKey: domain.ProkerBatalAktifTanggalKey,
		IsActive:   true,
		Status:     "active",
		StatusText: "Fitur Aktif",
	}, nil
}

func TestGetProkerBatalSettingQueryHandler(t *testing.T) {
	t.Run("returns status successfully", func(t *testing.T) {
		repo := &stubGetBatalRepo{}
		h := &GetProkerBatalSettingQueryHandler{Repo: repo}
		unitID := 5
		res, err := h.Handle(context.Background(), GetProkerBatalSettingQuery{
			TargetUnitID: &unitID,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || !res.IsActive {
			t.Fatalf("expected active status, got %+v", res)
		}
		if !repo.statusCalled || repo.targetUnit == nil || *repo.targetUnit != 5 {
			t.Fatalf("expected status called with unit 5, got unit=%v", repo.targetUnit)
		}
	})

	t.Run("returns error when repo is nil", func(t *testing.T) {
		h := &GetProkerBatalSettingQueryHandler{Repo: nil}
		_, err := h.Handle(context.Background(), GetProkerBatalSettingQuery{})
		if err == nil {
			t.Fatal("expected error when repo is nil")
		}
	})
}
