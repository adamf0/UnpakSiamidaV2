package UpdateVerifikasiDokumenV3

import (
	commondomain "UnpakSiamida/common/domain"
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type mockVerifikasiDokumenRepo struct {
	domain.IProkerV3Repository
	docCtx            *domain.DokumenUnitContext
	docCtxErr         error
	underFaculty      bool
	updatedStatus     string
	updatedCatatan    string
	updatedPenolakan  string
	updateCalledCount int
}

func (m *mockVerifikasiDokumenRepo) GetDokumenUnitContextByUuid(ctx context.Context, uid uuid.UUID) (*domain.DokumenUnitContext, error) {
	return m.docCtx, m.docCtxErr
}

func (m *mockVerifikasiDokumenRepo) IsUnitUnderAssignedFaculties(ctx context.Context, itemUnitID *int, itemUnitUUID string, assignedUnitIDs []uint) (bool, error) {
	return m.underFaculty, nil
}

func (m *mockVerifikasiDokumenRepo) UpdateVerifikasiDokumen(ctx context.Context, uid uuid.UUID, status, catatan, catatanPenolakan string) (*domain.MatrixProkerDokumen, error) {
	m.updateCalledCount++
	m.updatedStatus = status
	m.updatedCatatan = catatan
	m.updatedPenolakan = catatanPenolakan
	doc := m.docCtx.Dokumen
	doc.StatusVerifikasi = status
	doc.CatatanPenolakan = catatanPenolakan
	return &doc, nil
}

func assertDomainErrorCode(t *testing.T, err error, wantCode string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error %s, got nil", wantCode)
	}
	var derr commondomain.Error
	if !errors.As(err, &derr) {
		t.Fatalf("expected commondomain.Error, got %T (%v)", err, err)
	}
	if derr.Code != wantCode {
		t.Fatalf("expected error code %s, got %s (%s)", wantCode, derr.Code, derr.Description)
	}
}

func TestUpdateVerifikasiDokumenV3_Workflow(t *testing.T) {
	docID := uuid.New()
	prodiUnitID := 25
	unitUnitID := 10

	t.Run("prodi: fakultas matching parent verifies to terverifikasi_fakultas", func(t *testing.T) {
		repo := &mockVerifikasiDokumenRepo{
			docCtx: &domain.DokumenUnitContext{
				Dokumen:          domain.MatrixProkerDokumen{UUID: docID, StatusVerifikasi: "belum_terverifikasi"},
				IdFakultasUnit:   &prodiUnitID,
				TypeFakultasUnit: "prodi",
				NamaFakultasUnit: "ILMU KOMPUTER",
				ParentFakultas:   "MIPA",
			},
		}
		h := &UpdateVerifikasiDokumenV3CommandHandler{Repo: repo}
		res, err := h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:           docID.String(),
			StatusVerifikasi:  "terverifikasi",
			ActorLevel:        "fakultas",
			AssignedUnitNames: []string{"MIPA"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StatusVerifikasi != "terverifikasi_fakultas" || repo.updatedStatus != "terverifikasi_fakultas" {
			t.Fatalf("expected terverifikasi_fakultas, got %s", res.StatusVerifikasi)
		}
	})

	t.Run("prodi: fakultas outside parent faculty is rejected", func(t *testing.T) {
		repo := &mockVerifikasiDokumenRepo{
			docCtx: &domain.DokumenUnitContext{
				Dokumen:          domain.MatrixProkerDokumen{UUID: docID, StatusVerifikasi: "belum_terverifikasi"},
				IdFakultasUnit:   &prodiUnitID,
				TypeFakultasUnit: "prodi",
				NamaFakultasUnit: "ILMU KOMPUTER",
				ParentFakultas:   "MIPA",
			},
			underFaculty: false,
		}
		h := &UpdateVerifikasiDokumenV3CommandHandler{Repo: repo}
		_, err := h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:           docID.String(),
			StatusVerifikasi:  "terverifikasi_fakultas",
			ActorLevel:        "fakultas",
			AssignedUnitIDs:   []uint{99},
			AssignedUnitNames: []string{"HUKUM"},
		})
		assertDomainErrorCode(t, err, "ProkerV3.VerifikasiDokumenOutsideFaculty")
		if repo.updateCalledCount != 0 {
			t.Fatalf("expected no repo update")
		}
	})

	t.Run("prodi: admin cannot verify or reject before fakultas verifies", func(t *testing.T) {
		repo := &mockVerifikasiDokumenRepo{
			docCtx: &domain.DokumenUnitContext{
				Dokumen:          domain.MatrixProkerDokumen{UUID: docID, StatusVerifikasi: "belum_terverifikasi"},
				IdFakultasUnit:   &prodiUnitID,
				TypeFakultasUnit: "prodi",
				NamaFakultasUnit: "ILMU KOMPUTER",
				ParentFakultas:   "MIPA",
			},
		}
		h := &UpdateVerifikasiDokumenV3CommandHandler{Repo: repo}

		_, err := h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "terverifikasi",
			ActorLevel:       "admin",
		})
		assertDomainErrorCode(t, err, "ProkerV3.DokumenProdiRequiresFacultyFirst")

		_, err = h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "ditolak",
			CatatanPenolakan: "kurang lengkap",
			ActorLevel:       "admin",
		})
		assertDomainErrorCode(t, err, "ProkerV3.DokumenProdiRejectRequiresFacultyFirst")
	})

	t.Run("prodi: admin can verify and reject after terverifikasi_fakultas", func(t *testing.T) {
		repo := &mockVerifikasiDokumenRepo{
			docCtx: &domain.DokumenUnitContext{
				Dokumen:          domain.MatrixProkerDokumen{UUID: docID, StatusVerifikasi: "terverifikasi_fakultas"},
				IdFakultasUnit:   &prodiUnitID,
				TypeFakultasUnit: "prodi",
				NamaFakultasUnit: "ILMU KOMPUTER",
				ParentFakultas:   "MIPA",
			},
		}
		h := &UpdateVerifikasiDokumenV3CommandHandler{Repo: repo}

		res, err := h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "terverifikasi",
			ActorLevel:       "admin",
		})
		if err != nil || res.StatusVerifikasi != "terverifikasi" {
			t.Fatalf("expected terverifikasi, got %v err=%v", res, err)
		}

		res, err = h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "ditolak",
			CatatanPenolakan: "perbaiki lampiran",
			ActorLevel:       "admin",
		})
		if err != nil || res.StatusVerifikasi != "ditolak" {
			t.Fatalf("expected ditolak, got %v err=%v", res, err)
		}
	})

	t.Run("non-prodi: fakultas blocked, admin verifies directly", func(t *testing.T) {
		repo := &mockVerifikasiDokumenRepo{
			docCtx: &domain.DokumenUnitContext{
				Dokumen:          domain.MatrixProkerDokumen{UUID: docID, StatusVerifikasi: "belum_terverifikasi"},
				IdFakultasUnit:   &unitUnitID,
				TypeFakultasUnit: "unit",
				NamaFakultasUnit: "BPSI",
			},
		}
		h := &UpdateVerifikasiDokumenV3CommandHandler{Repo: repo}

		_, err := h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "terverifikasi_fakultas",
			ActorLevel:       "fakultas",
		})
		assertDomainErrorCode(t, err, "ProkerV3.DokumenNonProdiAdminOnly")

		res, err := h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "terverifikasi",
			ActorLevel:       "admin",
		})
		if err != nil || res.StatusVerifikasi != "terverifikasi" {
			t.Fatalf("expected terverifikasi, got %v err=%v", res, err)
		}
	})
}
