package UpdateVerifikasiDokumenV3

import (
	"UnpakSiamida/common/helper"
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strings"

	"github.com/google/uuid"
)

type UpdateVerifikasiDokumenV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateVerifikasiDokumenV3CommandHandler) Handle(ctx context.Context, cmd UpdateVerifikasiDokumenV3Command) (*domain.MatrixProkerDokumen, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	docUUID, err := uuid.Parse(cmd.DocUUID)
	if err != nil {
		return nil, domain.InvalidDocUuid()
	}

	level := strings.ToLower(strings.TrimSpace(cmd.ActorLevel))
	if level != "admin" && level != "fakultas" && level != "upmf" {
		return nil, domain.VerifikasiDokumenForbidden()
	}

	status := strings.TrimSpace(cmd.StatusVerifikasi)
	if status == "" {
		status = strings.TrimSpace(cmd.Status)
	}
	isReject := status == "ditolak" || status == "gagal_terverifikasi"

	docCtx, err := h.Repo.GetDokumenUnitContextByUuid(ctx, docUUID)
	if err != nil {
		return nil, err
	}
	if docCtx == nil {
		return nil, domain.DokumenNotFound()
	}

	finalStatus := status
	if docCtx.IsProdi() {
		// Unit pengusul bertipe PRODI: verifikasi fakultas terlebih dahulu, lalu verifikasi admin
		if level == "fakultas" || level == "upmf" {
			allowed, err := h.canFacultyVerify(ctx, docCtx, cmd)
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, domain.VerifikasiDokumenOutsideFaculty()
			}
			if docCtx.Dokumen.StatusVerifikasi == "terverifikasi" {
				return nil, domain.DokumenAlreadyVerifiedByAdmin()
			}
			if isReject {
				finalStatus = "ditolak"
			} else {
				finalStatus = "terverifikasi_fakultas"
			}
		} else {
			// level == "admin"
			if isReject {
				if docCtx.Dokumen.StatusVerifikasi != "terverifikasi_fakultas" && docCtx.Dokumen.StatusVerifikasi != "terverifikasi" {
					return nil, domain.DokumenProdiRejectRequiresFacultyFirst()
				}
				finalStatus = "ditolak"
			} else {
				if docCtx.Dokumen.StatusVerifikasi != "terverifikasi_fakultas" {
					return nil, domain.DokumenProdiRequiresFacultyFirst()
				}
				finalStatus = "terverifikasi"
			}
		}
	} else {
		// Unit pengusul SELAIN PRODI: cukup verifikasi admin
		if level == "fakultas" || level == "upmf" {
			return nil, domain.DokumenNonProdiAdminOnly()
		}
		if isReject {
			finalStatus = "ditolak"
		} else {
			finalStatus = "terverifikasi"
		}
	}

	return h.Repo.UpdateVerifikasiDokumen(ctx, docUUID, finalStatus, cmd.Catatan, cmd.CatatanPenolakan)
}

func (h *UpdateVerifikasiDokumenV3CommandHandler) canFacultyVerify(
	ctx context.Context,
	docCtx *domain.DokumenUnitContext,
	cmd UpdateVerifikasiDokumenV3Command,
) (bool, error) {
	// 1. Direct match berdasarkan ID / nama unit
	if docCtx.IdFakultasUnit != nil && *docCtx.IdFakultasUnit > 0 {
		for _, aid := range cmd.AssignedUnitIDs {
			if uint(*docCtx.IdFakultasUnit) == aid {
				return true, nil
			}
		}
	}
	itemUnitName := strings.TrimSpace(docCtx.NamaFakultasUnit)
	if itemUnitName != "" {
		for _, aname := range cmd.AssignedUnitNames {
			if strings.EqualFold(itemUnitName, strings.TrimSpace(aname)) {
				return true, nil
			}
		}
	}
	if len(cmd.AssignedUnitIDs) == 0 && helper.MatchUnit(docCtx.IdFakultasUnit, docCtx.UuidFakultasUnit, cmd.ActorUserFU) {
		return true, nil
	}

	// 2. Match parent fakultas berdasarkan nama penugasan
	parentFaculty := strings.TrimSpace(docCtx.ParentFakultas)
	if parentFaculty != "" {
		for _, aname := range cmd.AssignedUnitNames {
			if strings.EqualFold(parentFaculty, strings.TrimSpace(aname)) {
				return true, nil
			}
		}
	}

	// 3. Fallback cek relasi fakultas-prodi di repository
	return h.Repo.IsUnitUnderAssignedFaculties(ctx, docCtx.IdFakultasUnit, docCtx.UuidFakultasUnit, cmd.AssignedUnitIDs)
}
