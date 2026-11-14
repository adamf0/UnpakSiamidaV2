package UpdateUserAbac

import (
	"context"
	"strings"
	"time"

	domainuser "UnpakSiamida/modules/user/domain"

	"github.com/google/uuid"
)

type UpdateUserAbacCommandHandler struct {
	Repo domainuser.IUserRepository
}

func (h *UpdateUserAbacCommandHandler) Handle(ctx context.Context, cmd UpdateUserAbacCommand) (string, error) {
	rawUUID := strings.TrimSpace(cmd.UUID)
	if rawUUID == "" {
		return "", domainuser.AbacUuidRequired()
	}

	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return "", domainuser.AbacInvalidUuid()
	}

	existing, err := h.Repo.GetUserAbacByUuid(ctx, parsedUUID)
	if err != nil || existing == nil {
		return "", domainuser.AbacNotFound()
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if cmd.Status != "" {
		updates["status"] = strings.ToLower(strings.TrimSpace(cmd.Status))
	}
	if cmd.Tahun != "" {
		updates["tahun"] = strings.TrimSpace(cmd.Tahun)
	}

	var targetFuID *uint = existing.IdFakultasUnit
	if cmd.IdFakultasUnit != nil {
		if *cmd.IdFakultasUnit > 0 {
			targetFuID = cmd.IdFakultasUnit
			updates["id_fakultas_unit"] = *cmd.IdFakultasUnit
		} else {
			targetFuID = nil
			updates["id_fakultas_unit"] = nil
		}
	} else if strings.TrimSpace(cmd.UuidFakultasUnit) != "" {
		fuID, _, _, err := h.Repo.GetFakultasUnitInfo(ctx, nil, cmd.UuidFakultasUnit)
		if err == nil {
			targetFuID = &fuID
			updates["id_fakultas_unit"] = fuID
		}
	}

	if cmd.Role != "" {
		reqRole := strings.ToLower(strings.TrimSpace(cmd.Role))
		if reqRole == "admin" {
			return "", domainuser.AbacAdminNotAllowed()
		}
		if reqRole == "auditee" {
			reqRole = "user"
		}
		if reqRole == "evaluator" || reqRole == "anggaran" {
			updates["id_fakultas_unit"] = nil
			targetFuID = nil
		}
		if reqRole != "auto" {
			if targetFuID != nil && *targetFuID > 0 {
				_, fuType, fuNama, err := h.Repo.GetFakultasUnitInfo(ctx, targetFuID, "")
				if err == nil {
					lowType := strings.ToLower(strings.TrimSpace(fuType))
					lowNama := strings.ToLower(strings.TrimSpace(fuNama))
					isFak := lowType == "fakultas" || strings.Contains(lowNama, "(fakultas)") || strings.HasPrefix(lowNama, "fakultas")
					isProdi := lowType == "prodi" || strings.Contains(lowNama, "(prodi)") || strings.HasPrefix(lowNama, "prodi") || strings.HasPrefix(lowNama, "program studi")

					if reqRole == "upms" && !isProdi {
						return "", domainuser.AbacTargetUpmsInvalid()
					}
					if (reqRole == "upmf" || reqRole == "fakultas") && !isFak {
						return "", domainuser.AbacTargetUpmfInvalid()
					}
					if reqRole == "upml" && isProdi {
						return "", domainuser.AbacTargetUpmlInvalid()
					}
				}
			} else {
				if reqRole != "evaluator" && reqRole != "anggaran" {
					return "", domainuser.AbacTargetRequired()
				}
			}
			updates["role"] = reqRole
		}
	}

	if err := h.Repo.UpdateUserAbac(ctx, existing.ID, updates); err != nil {
		return "", err
	}

	return "Penugasan role berhasil diperbarui", nil
}
