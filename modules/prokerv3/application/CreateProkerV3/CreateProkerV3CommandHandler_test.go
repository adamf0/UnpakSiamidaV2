package CreateProkerV3

import (
	commondomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
	"testing"
)

type stubCreateProkerRepo struct {
	domain.IProkerV3Repository
	existsResult bool
	existsErr    error

	checkedFuID    int
	checkedTahun   string
	checkedMpID    *int
	checkedNama    string
	createCalled   bool
	resolveMpCalls int
}

func (s *stubCreateProkerRepo) ExistsActiveProker(_ context.Context, fuID int, tahunStr string, mpID *int, namaProg string, _ uint) (bool, error) {
	s.checkedFuID = fuID
	s.checkedTahun = tahunStr
	s.checkedMpID = mpID
	s.checkedNama = namaProg
	return s.existsResult, s.existsErr
}

func (s *stubCreateProkerRepo) ResolveMasterTahun(_ context.Context, _ string) (int, error) {
	return 1, nil
}

func (s *stubCreateProkerRepo) ResolveOrCreateMataProgram(_ context.Context, _ string, _ int, _ *int) (*int, error) {
	s.resolveMpCalls++
	id := 10
	return &id, nil
}

func (s *stubCreateProkerRepo) UpdateMataProgramTahun(_ context.Context, _ int, _ int) error {
	return nil
}

func (s *stubCreateProkerRepo) Create(_ context.Context, _ *domain.ProkerV3) error {
	s.createCalled = true
	return nil
}

func (s *stubCreateProkerRepo) SaveMatrixProkerScopes(_ context.Context, _ []domain.MatrixProkerScope) error {
	return nil
}

func TestCreateProkerV3CommandHandler_Uniqueness(t *testing.T) {
	fuID := 25
	mpID := 99

	t.Run("blocks duplicate (mata_program, fakultas_unit, tahun) before mutating DB", func(t *testing.T) {
		repo := &stubCreateProkerRepo{existsResult: true}
		h := &CreateProkerV3CommandHandler{Repo: repo}

		_, err := h.Handle(context.Background(), CreateProkerV3Command{
			IdFakultasUnit: &fuID,
			IdMataProgram:  &mpID,
			NamaProgram:    "Program Uji Coba",
			Tahun:          "2027",
		})
		if err == nil {
			t.Fatalf("expected conflict error, got nil")
		}
		var domErr commondomain.Error
		if !errors.As(err, &domErr) || domErr.Code != "ProkerV3.DuplicateProkerUnitTahun" {
			t.Fatalf("expected ProkerV3.DuplicateProkerUnitTahun, got %v", err)
		}
		if repo.createCalled || repo.resolveMpCalls > 0 {
			t.Fatalf("expected no DB writes on duplicate, create=%v resolveMpCalls=%d", repo.createCalled, repo.resolveMpCalls)
		}
		if repo.checkedFuID != fuID || repo.checkedTahun != "2027" || repo.checkedMpID == nil || *repo.checkedMpID != mpID {
			t.Fatalf("unexpected uniqueness check args: fu=%d tahun=%q mp=%v", repo.checkedFuID, repo.checkedTahun, repo.checkedMpID)
		}
	})

	t.Run("allows unique (mata_program, fakultas_unit, tahun)", func(t *testing.T) {
		repo := &stubCreateProkerRepo{existsResult: false}
		h := &CreateProkerV3CommandHandler{Repo: repo}

		uid, err := h.Handle(context.Background(), CreateProkerV3Command{
			IdFakultasUnit: &fuID,
			IdMataProgram:  &mpID,
			NamaProgram:    "Program Baru",
			Tahun:          "2027",
		})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if uid == "" || !repo.createCalled {
			t.Fatalf("expected created uuid and repo.Create called, uid=%q create=%v", uid, repo.createCalled)
		}
	})
}
