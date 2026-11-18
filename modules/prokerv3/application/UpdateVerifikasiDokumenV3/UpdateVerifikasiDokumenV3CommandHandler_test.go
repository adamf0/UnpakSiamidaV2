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
	lastUpdate        domain.DokumenVerifikasiUpdate
	updateCalledCount int
}

func (m *mockVerifikasiDokumenRepo) GetDokumenUnitContextByUuid(ctx context.Context, uid uuid.UUID) (*domain.DokumenUnitContext, error) {
	return m.docCtx, m.docCtxErr
}

func (m *mockVerifikasiDokumenRepo) IsUnitUnderAssignedFaculties(ctx context.Context, itemUnitID *int, itemUnitUUID string, assignedUnitIDs []uint) (bool, error) {
	return m.underFaculty, nil
}

func (m *mockVerifikasiDokumenRepo) UpdateVerifikasiDokumen(ctx context.Context, uid uuid.UUID, update domain.DokumenVerifikasiUpdate) (*domain.MatrixProkerDokumen, error) {
	m.updateCalledCount++
	m.lastUpdate = update
	doc := m.docCtx.Dokumen
	doc.StatusVerifikasi = update.StatusVerifikasi
	doc.CatatanPenolakan = update.CatatanPenolakan
	if update.StatusVerifikasiFakultas != nil {
		doc.StatusVerifikasiFakultas = *update.StatusVerifikasiFakultas
	}
	if update.CatatanTolakFakultas != nil {
		doc.CatatanTolakFakultas = *update.CatatanTolakFakultas
	}
	if update.StatusVerifikasiAdmin != nil {
		doc.StatusVerifikasiAdmin = *update.StatusVerifikasiAdmin
	}
	if update.CatatanTolakAdmin != nil {
		doc.CatatanTolakAdmin = *update.CatatanTolakAdmin
	}
	m.docCtx.Dokumen = doc
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

	t.Run("prodi: fakultas reject sets fakultas status & note without touching admin", func(t *testing.T) {
		repo := &mockVerifikasiDokumenRepo{
			docCtx: &domain.DokumenUnitContext{
				Dokumen: domain.MatrixProkerDokumen{
					UUID:                     docID,
					StatusVerifikasi:         "belum_terverifikasi",
					StatusVerifikasiFakultas: "belum_terverifikasi",
					StatusVerifikasiAdmin:    "belum_terverifikasi",
				},
				IdFakultasUnit:   &prodiUnitID,
				TypeFakultasUnit: "prodi",
				NamaFakultasUnit: "ILMU KOMPUTER",
				ParentFakultas:   "MIPA",
			},
		}
		h := &UpdateVerifikasiDokumenV3CommandHandler{Repo: repo}

		res, err := h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:           docID.String(),
			StatusVerifikasi:  "ditolak",
			CatatanPenolakan:  "SK belum ditandatangani Dekan",
			ActorLevel:        "fakultas",
			AssignedUnitNames: []string{"MIPA"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StatusVerifikasiFakultas != "ditolak" || res.CatatanTolakFakultas != "SK belum ditandatangani Dekan" {
			t.Fatalf("unexpected faculty state: %+v", res)
		}
		if res.StatusVerifikasiAdmin != "belum_terverifikasi" || res.CatatanTolakAdmin != "" {
			t.Fatalf("admin fields should remain empty: %+v", res)
		}
		// Admin still cannot reject or verify while rejected by faculty
		_, err = h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "ditolak",
			CatatanPenolakan: "coba tolak admin",
			ActorLevel:       "admin",
		})
		assertDomainErrorCode(t, err, "ProkerV3.DokumenProdiRejectRequiresFacultyFirst")
	})

	t.Run("prodi: fakultas verifies then admin rejects keeping fakultas verified", func(t *testing.T) {
		repo := &mockVerifikasiDokumenRepo{
			docCtx: &domain.DokumenUnitContext{
				Dokumen: domain.MatrixProkerDokumen{
					UUID:                     docID,
					StatusVerifikasi:         "belum_terverifikasi",
					StatusVerifikasiFakultas: "belum_terverifikasi",
					StatusVerifikasiAdmin:    "belum_terverifikasi",
				},
				IdFakultasUnit:   &prodiUnitID,
				TypeFakultasUnit: "prodi",
				NamaFakultasUnit: "ILMU KOMPUTER",
				ParentFakultas:   "MIPA",
			},
		}
		h := &UpdateVerifikasiDokumenV3CommandHandler{Repo: repo}

		// Step 1: Fakultas verifies
		res, err := h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:           docID.String(),
			StatusVerifikasi:  "terverifikasi_fakultas",
			ActorLevel:        "fakultas",
			AssignedUnitNames: []string{"MIPA"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StatusVerifikasi != "terverifikasi_fakultas" || res.StatusVerifikasiFakultas != "terverifikasi" || res.StatusVerifikasiAdmin != "belum_terverifikasi" {
			t.Fatalf("unexpected state after faculty verify: %+v", res)
		}

		// Step 2: Admin rejects with admin-specific note
		res, err = h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "ditolak",
			CatatanPenolakan: "Lampiran RAB kurang detail",
			ActorLevel:       "admin",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StatusVerifikasiFakultas != "terverifikasi" || res.CatatanTolakFakultas != "" {
			t.Fatalf("faculty verification must stay intact: %+v", res)
		}
		if res.StatusVerifikasiAdmin != "ditolak" || res.CatatanTolakAdmin != "Lampiran RAB kurang detail" || res.StatusVerifikasi != "ditolak" {
			t.Fatalf("admin rejection fields not set properly: %+v", res)
		}

		// Step 3: Admin can later approve directly because faculty is still verified
		res, err = h.Handle(context.Background(), UpdateVerifikasiDokumenV3Command{
			DocUUID:          docID.String(),
			StatusVerifikasi: "terverifikasi",
			ActorLevel:       "admin",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.StatusVerifikasi != "terverifikasi" || res.StatusVerifikasiAdmin != "terverifikasi" || res.CatatanTolakAdmin != "" {
			t.Fatalf("expected full verified state: %+v", res)
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
	})

	t.Run("non-prodi: fakultas blocked, admin verifies/rejects on admin columns", func(t *testing.T) {
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
			StatusVerifikasi: "ditolak",
			CatatanPenolakan: "Format file salah",
			ActorLevel:       "admin",
		})
		if err != nil || res.StatusVerifikasiAdmin != "ditolak" || res.CatatanTolakAdmin != "Format file salah" {
			t.Fatalf("unexpected non-prodi admin reject: %+v err=%v", res, err)
		}
	})
}
