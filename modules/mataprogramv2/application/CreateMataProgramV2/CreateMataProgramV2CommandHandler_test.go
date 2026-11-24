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
	tahunID          uint
	savedTahunID     uint
	savedMataProgram string
	savedTemplates   []uint
	saveErr          error
}

func (m *mockRepo) ResolveOrCreateTahunID(ctx context.Context, tahunStr string) (uint, error) {
	return m.tahunID, nil
}

func (m *mockRepo) SaveMataProgramV2WithTemplates(ctx context.Context, tahunID uint, mataProgram string, templateRenstraIDs []uint) (string, uint, error) {
	if m.saveErr != nil {
		return "", 0, m.saveErr
	}
	m.savedTahunID = tahunID
	m.savedMataProgram = mataProgram
	m.savedTemplates = append([]uint(nil), templateRenstraIDs...)
	return "uuid-created", uint(len(templateRenstraIDs)), nil
}

func TestCreateMataProgramV2CommandHandler(t *testing.T) {
	repo := &mockRepo{tahunID: 15}
	h := &CreateMataProgramV2CommandHandler{Repo: repo}

	t.Run("requires mata_program", func(t *testing.T) {
		_, err := h.Handle(context.Background(), CreateMataProgramV2Command{TemplateRenstraIDs: []uint{101}})
		var derr commondomain.Error
		if !errors.As(err, &derr) || derr.Code != "MataProgramV2.MataProgramRequired" {
			t.Fatalf("expected MataProgramRequired, got %v", err)
		}
	})

	t.Run("requires template_renstras", func(t *testing.T) {
		_, err := h.Handle(context.Background(), CreateMataProgramV2Command{MataProgram: "Program A"})
		var derr commondomain.Error
		if !errors.As(err, &derr) || derr.Code != "MataProgramV2.TargetsRequired" {
			t.Fatalf("expected TargetsRequired, got %v", err)
		}
	})

	t.Run("saves mata program with template renstras", func(t *testing.T) {
		res, err := h.Handle(context.Background(), CreateMataProgramV2Command{
			MataProgram:        "Program A",
			Tahun:              "2027",
			TemplateRenstraIDs: []uint{101, 102},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Count != 2 || res.UUID != "uuid-created" || repo.savedTahunID != 15 || len(repo.savedTemplates) != 2 {
			t.Fatalf("unexpected result: %+v repo=%+v", res, repo)
		}
	})

	t.Run("propagates duplicate error", func(t *testing.T) {
		dupRepo := &mockRepo{tahunID: 15, saveErr: domain.DuplicateMataProgram()}
		dupHandler := &CreateMataProgramV2CommandHandler{Repo: dupRepo}
		_, err := dupHandler.Handle(context.Background(), CreateMataProgramV2Command{
			MataProgram:        "Program A",
			Tahun:              "2027",
			TemplateRenstraIDs: []uint{101},
		})
		var derr commondomain.Error
		if !errors.As(err, &derr) || derr.Code != "MataProgramV2.DuplicateMataProgram" {
			t.Fatalf("expected DuplicateMataProgram, got %v", err)
		}
	})
}
