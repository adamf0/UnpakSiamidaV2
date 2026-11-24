package UpdateMataProgramV2

import (
	commondomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"errors"
	"testing"
)

type mockUpdateRepo struct {
	domain.IMataProgramV2Repository
	record           *domain.MataProgramV2Record
	resolvedTahunID  uint
	updatedNama      string
	updatedTahunID   uint
	updatedTemplates []uint
	replaced         bool
	updateErr        error
}

func (m *mockUpdateRepo) GetRecordByUuid(ctx context.Context, rawUuid string, includeDeleted bool) (*domain.MataProgramV2Record, error) {
	if m.record == nil {
		return nil, domain.NotFoundOrDeleted()
	}
	return m.record, nil
}

func (m *mockUpdateRepo) ResolveOrCreateTahunID(ctx context.Context, tahunStr string) (uint, error) {
	return m.resolvedTahunID, nil
}

func (m *mockUpdateRepo) UpdateMataProgramV2WithTemplates(ctx context.Context, rawUuid string, tahunID uint, mataProgram string, templateRenstraIDs []uint, replaceTemplates bool) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updatedNama = mataProgram
	m.updatedTahunID = tahunID
	m.updatedTemplates = append([]uint(nil), templateRenstraIDs...)
	m.replaced = replaceTemplates
	return nil
}

func TestUpdateMataProgramV2CommandHandler(t *testing.T) {
	tahun15 := uint(15)
	baseRec := &domain.MataProgramV2Record{
		ID:            11,
		UUID:          "04325df3-5050-47b2-94e6-2ca5d762413b",
		MataProgram:   "Penerimaan Mahasiswa Baru",
		IdMasterTahun: &tahun15,
	}

	t.Run("rejects empty templates when template field explicitly provided", func(t *testing.T) {
		repo := &mockUpdateRepo{record: baseRec, resolvedTahunID: 15}
		h := &UpdateMataProgramV2CommandHandler{Repo: repo}
		_, err := h.Handle(context.Background(), UpdateMataProgramV2Command{
			UUID:         baseRec.UUID,
			MataProgram:  "PMB Baru",
			Tahun:        "2027",
			HasTemplates: true,
		})
		var derr commondomain.Error
		if !errors.As(err, &derr) || derr.Code != "MataProgramV2.TargetsRequired" {
			t.Fatalf("expected TargetsRequired, got %v", err)
		}
	})

	t.Run("updates name, year, and templates without id_fakultas_unit", func(t *testing.T) {
		repo := &mockUpdateRepo{record: baseRec, resolvedTahunID: 16}
		h := &UpdateMataProgramV2CommandHandler{Repo: repo}
		res, err := h.Handle(context.Background(), UpdateMataProgramV2Command{
			UUID:               baseRec.UUID,
			MataProgram:        "PMB 2028",
			Tahun:              "2028",
			TemplateRenstraIDs: []uint{501, 502},
			HasTemplates:       true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.UUID != baseRec.UUID || repo.updatedNama != "PMB 2028" || repo.updatedTahunID != 16 || !repo.replaced || len(repo.updatedTemplates) != 2 {
			t.Fatalf("unexpected state: res=%+v repo=%+v", res, repo)
		}
	})

	t.Run("propagates duplicate conflict", func(t *testing.T) {
		repo := &mockUpdateRepo{record: baseRec, resolvedTahunID: 15, updateErr: domain.DuplicateMataProgram()}
		h := &UpdateMataProgramV2CommandHandler{Repo: repo}
		_, err := h.Handle(context.Background(), UpdateMataProgramV2Command{
			UUID:               baseRec.UUID,
			MataProgram:        "Internasionalisasi",
			Tahun:              "2027",
			TemplateRenstraIDs: []uint{10},
			HasTemplates:       true,
		})
		var derr commondomain.Error
		if !errors.As(err, &derr) || derr.Code != "MataProgramV2.DuplicateMataProgram" {
			t.Fatalf("expected DuplicateMataProgram, got %v", err)
		}
	})
}
