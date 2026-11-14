package UpdateProkerBatalSetting

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"testing"
	"time"
)

type stubUpdateBatalRepo struct {
	domain.IProkerV3Repository
	rawVal     string
	savedVal   string
	saveCalled bool
}

func (s *stubUpdateBatalRepo) GetProkerBatalSettingRaw(_ context.Context) (string, error) {
	return s.rawVal, nil
}

func (s *stubUpdateBatalRepo) SaveProkerBatalSetting(_ context.Context, serialized string, _ *int) error {
	s.savedVal = serialized
	s.saveCalled = true
	return nil
}

func (s *stubUpdateBatalRepo) GetProkerBatalSettingStatus(_ context.Context, _ time.Time, _ *int) (*domain.ProkerBatalSettingStatus, error) {
	return &domain.ProkerBatalSettingStatus{
		SettingKey: domain.ProkerBatalAktifTanggalKey,
		IsActive:   true,
		Status:     "active",
		StatusText: "Fitur Aktif",
	}, nil
}

func (s *stubUpdateBatalRepo) GetUserIDByUUID(_ context.Context, _ string) (*int, error) {
	id := 1
	return &id, nil
}

func TestUpdateProkerBatalSettingCommandHandler(t *testing.T) {
	t.Run("rejects non-admin", func(t *testing.T) {
		repo := &stubUpdateBatalRepo{}
		h := &UpdateProkerBatalSettingCommandHandler{Repo: repo}
		_, err := h.Handle(context.Background(), UpdateProkerBatalSettingCommand{
			IsAdmin: false,
		})
		if err == nil {
			t.Fatal("expected error for non-admin")
		}
	})

	t.Run("updates global date successfully", func(t *testing.T) {
		repo := &stubUpdateBatalRepo{}
		h := &UpdateProkerBatalSettingCommandHandler{Repo: repo}
		targetDate := "2026-12-01"
		res, err := h.Handle(context.Background(), UpdateProkerBatalSettingCommand{
			IsAdmin:    true,
			GlobalDate: &targetDate,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !repo.saveCalled {
			t.Fatal("expected SaveProkerBatalSetting to be called")
		}
		if res == nil {
			t.Fatal("expected non-nil result")
		}
	})
}
