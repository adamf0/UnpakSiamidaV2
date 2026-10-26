package CreateUserAbac

import (
	"context"
	"errors"
	"strings"
	"time"

	domainuser "UnpakSiamida/modules/user/domain"

	"github.com/google/uuid"
)

type CreateUserAbacCommandHandler struct {
	Repo domainuser.IUserRepository
}

func (h *CreateUserAbacCommandHandler) Handle(ctx context.Context, cmd CreateUserAbacCommand) (*CreateUserAbacResult, error) {
	tahun := strings.TrimSpace(cmd.Tahun)
	if tahun == "" {
		return nil, errors.New("Tahun penugasan wajib diisi")
	}

	// Resolve User ID
	var userID uint
	if cmd.IdUser != nil && *cmd.IdUser > 0 {
		userID = *cmd.IdUser
	} else if strings.TrimSpace(cmd.UserUuid) != "" {
		parsedUuid, err := uuid.Parse(strings.TrimSpace(cmd.UserUuid))
		if err != nil {
			return nil, errors.New("UUID User tidak valid")
		}
		u, err := h.Repo.GetByUuid(ctx, parsedUuid)
		if err != nil || u == nil {
			return nil, errors.New("User yang dipilih tidak ditemukan dalam sistem")
		}
		userID = u.ID
	} else {
		return nil, errors.New("Pengguna wajib dipilih")
	}

	// Resolve Fakultas Unit & Role
	var fuID *uint
	role := strings.ToLower(strings.TrimSpace(cmd.Role))

	hasTarget := (cmd.IdFakultasUnit != nil && *cmd.IdFakultasUnit > 0) || strings.TrimSpace(cmd.UuidFakultasUnit) != ""
	if !hasTarget {
		if role == "upms" {
			// UPMS does not require a target unit (applies university-wide / school-wide)
			fuID = nil
		} else if role == "upmf" {
			return nil, errors.New("Target Fakultas wajib dipilih untuk role UPMF")
		} else {
			return nil, errors.New("Fakultas/Unit wajib dipilih untuk role ini")
		}
	} else {
		resolvedID, fuType, fuNama, err := h.Repo.GetFakultasUnitInfo(ctx, cmd.IdFakultasUnit, cmd.UuidFakultasUnit)
		if err != nil {
			return nil, errors.New("Fakultas/Unit yang dipilih tidak ditemukan")
		}
		fuID = &resolvedID

		// Rule: "kalau yg saya pilih fakultas maka hanya bisa pilih upml dan upmf, kalau bukan fakultas maka bisa pilih upml upms"
		isFakultas := fuType == "fakultas" || strings.Contains(strings.ToLower(fuNama), "(fakultas)") || strings.HasPrefix(strings.ToLower(fuNama), "fakultas")

		if role == "" || role == "auto" {
			if isFakultas {
				role = "upmf"
			} else {
				lowNama := strings.ToLower(fuNama)
				if strings.Contains(lowNama, "pascasarjana") || strings.Contains(lowNama, "vokasi") || strings.Contains(lowNama, "sekolah") {
					role = "upms"
				} else {
					role = "upml"
				}
			}
		} else {
			if isFakultas && role == "upms" {
				return nil, errors.New("Target Fakultas hanya dapat memilih role UPML atau UPMF")
			}
			if !isFakultas && role == "upmf" {
				return nil, errors.New("Target non-Fakultas hanya dapat memilih role UPML atau UPMS")
			}
		}
	}

	status := strings.ToLower(strings.TrimSpace(cmd.Status))
	if status == "" {
		status = "active"
	}

	// Check duplicate
	exists, err := h.Repo.CheckUserAbacExists(ctx, userID, tahun, fuID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Penugasan untuk user, tahun, dan target unit ini sudah ada")
	}

	newUUID := uuid.New()
	now := time.Now()
	abac := &domainuser.UserAbac{
		UUID:           newUUID,
		IdUser:         userID,
		Tahun:          tahun,
		IdFakultasUnit: fuID,
		Role:           role,
		Status:         status,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := h.Repo.CreateUserAbac(ctx, abac); err != nil {
		return nil, err
	}

	return &CreateUserAbacResult{
		Message: "Penugasan role berhasil disimpan",
		UUID:    newUUID.String(),
		ID:      abac.ID,
		Role:    role,
	}, nil
}
