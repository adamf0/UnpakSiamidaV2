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
		return nil, errors.New("Target Fakultas/Unit/Prodi wajib dipilih")
	}

	resolvedID, fuType, fuNama, err := h.Repo.GetFakultasUnitInfo(ctx, cmd.IdFakultasUnit, cmd.UuidFakultasUnit)
	if err != nil {
		return nil, errors.New("Fakultas/Unit/Prodi yang dipilih tidak ditemukan")
	}
	fuID = &resolvedID

	lowType := strings.ToLower(strings.TrimSpace(fuType))
	lowNama := strings.ToLower(strings.TrimSpace(fuNama))

	isFakultas := lowType == "fakultas" || strings.Contains(lowNama, "(fakultas)") || strings.HasPrefix(lowNama, "fakultas")
	isProdi := lowType == "prodi" || strings.Contains(lowNama, "(prodi)") || strings.HasPrefix(lowNama, "prodi") || strings.HasPrefix(lowNama, "program studi")

	if role == "" || role == "auto" {
		if isProdi {
			role = "upms"
		} else if isFakultas {
			role = "upmf"
		} else {
			role = "upml"
		}
	} else {
		if role != "upms" && role != "upmf" && role != "upml" {
			return nil, errors.New("Role penugasan tidak valid (hanya diperbolehkan UPMS, UPMF, atau UPML)")
		}
		if role == "upms" && !isProdi {
			return nil, errors.New("Target untuk role UPMS harus berupa Program Studi (Prodi)")
		}
		if role == "upmf" && !isFakultas {
			return nil, errors.New("Target untuk role UPMF harus berupa Fakultas")
		}
		if role == "upml" && isProdi {
			return nil, errors.New("Target untuk role UPML harus berupa Fakultas atau Unit (bukan Program Studi)")
		}
	}

	if role != "upms" && role != "upmf" && role != "upml" {
		return nil, errors.New("Role penugasan tidak valid (hanya diperbolehkan UPMS, UPMF, atau UPML)")
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
