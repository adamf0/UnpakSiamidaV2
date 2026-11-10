package DeleteProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"testing"

	"github.com/google/uuid"
)

type stubDeleteRepo struct {
	domain.IProkerV3Repository
	hardCalls int
	softCalls int
	lastSoft  domain.SoftDeleteProkerParams
}

func (s *stubDeleteRepo) Delete(_ context.Context, _ uuid.UUID) error {
	s.hardCalls++
	return nil
}

func (s *stubDeleteRepo) SoftDelete(_ context.Context, params domain.SoftDeleteProkerParams) error {
	s.softCalls++
	s.lastSoft = params
	return nil
}

func TestDeleteProkerV3CommandHandler(t *testing.T) {
	uid := uuid.New().String()

	t.Run("soft delete rejects empty reason", func(t *testing.T) {
		repo := &stubDeleteRepo{}
		h := &DeleteProkerV3CommandHandler{Repo: repo}
		_, err := h.Handle(context.Background(), DeleteProkerV3Command{
			UUID:       uid,
			SoftDelete: true,
			Catatan:    "   ",
		})
		if err == nil {
			t.Fatal("expected error for empty soft delete reason")
		}
		if repo.softCalls != 0 || repo.hardCalls != 0 {
			t.Fatalf("expected zero repo calls, got soft=%d hard=%d", repo.softCalls, repo.hardCalls)
		}
	})

	t.Run("soft delete delegates to SoftDelete with trimmed metadata", func(t *testing.T) {
		repo := &stubDeleteRepo{}
		h := &DeleteProkerV3CommandHandler{Repo: repo}
		actorID := 12
		gotUUID, err := h.Handle(context.Background(), DeleteProkerV3Command{
			UUID:          uid,
			SoftDelete:    true,
			Catatan:       "  Program kerja digabung ke kegiatan lain  ",
			Role:          "user",
			ActorID:       &actorID,
			ActorName:     "Aries Maesya, M.Kom",
			ActorUsername: "aries",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotUUID != uid || repo.softCalls != 1 || repo.hardCalls != 0 {
			t.Fatalf("unexpected result: uuid=%s soft=%d hard=%d", gotUUID, repo.softCalls, repo.hardCalls)
		}
		if repo.lastSoft.Catatan != "Program kerja digabung ke kegiatan lain" {
			t.Fatalf("got catatan %q", repo.lastSoft.Catatan)
		}
	})

	t.Run("hard delete path preserved for admin", func(t *testing.T) {
		repo := &stubDeleteRepo{}
		h := &DeleteProkerV3CommandHandler{Repo: repo}
		if _, err := h.Handle(context.Background(), DeleteProkerV3Command{
			UUID:       uid,
			SoftDelete: false,
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.hardCalls != 1 || repo.softCalls != 0 {
			t.Fatalf("expected hard=1 soft=0, got hard=%d soft=%d", repo.hardCalls, repo.softCalls)
		}
	})
}
