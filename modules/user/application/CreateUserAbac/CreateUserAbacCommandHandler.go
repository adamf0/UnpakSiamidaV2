package CreateUserAbac

import (
	"context"
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
		return nil, domainuser.AbacTahunRequired()
	}

	// Resolve User ID
	var userID uint
	if cmd.IdUser != nil && *cmd.IdUser > 0 {
		userID = *cmd.IdUser
	} else if strings.TrimSpace(cmd.UserUuid) != "" {
		parsedUuid, err := uuid.Parse(strings.TrimSpace(cmd.UserUuid))
		if err != nil {
			return nil, domainuser.AbacUserUuidInvalid()
		}
		u, err := h.Repo.GetByUuid(ctx, parsedUuid)
		if err != nil || u == nil {
			return nil, domainuser.AbacUserNotFound()
		}
		userID = u.ID
	} else {
		return nil, domainuser.AbacUserRequired()
	}

	// Resolve Fakultas Unit & Role
	var fuID *uint
	role := strings.ToLower(strings.TrimSpace(cmd.Role))
	if role == "admin" {
		return nil, domainuser.AbacAdminNotAllowed()
	}
	if role == "auditee" {
		role = "user"
	}

	hasTarget := (cmd.IdFakultasUnit != nil && *cmd.IdFakultasUnit > 0) || strings.TrimSpace(cmd.UuidFakultasUnit) != ""
	if hasTarget {
		resolvedID, fuType, fuNama, err := h.Repo.GetFakultasUnitInfo(ctx, cmd.IdFakultasUnit, cmd.UuidFakultasUnit)
		if err != nil {
			return nil, domainuser.AbacFakultasUnitNotFound()
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
			if role == "upms" && !isProdi {
				return nil, domainuser.AbacTargetUpmsInvalid()
			}
			if (role == "upmf" || role == "fakultas") && !isFakultas {
				return nil, domainuser.AbacTargetUpmfInvalid()
			}
			if role == "upml" && isProdi {
				return nil, domainuser.AbacTargetUpmlInvalid()
			}
		}
	} else {
		if role != "evaluator" && role != "anggaran" {
			return nil, domainuser.AbacTargetRequired()
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
		return nil, domainuser.AbacDuplicateAssignment()
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
