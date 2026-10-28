package UpdateUserAbac

import (
	"context"
	"errors"
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
		return "", errors.New("UUID penugasan wajib disertakan")
	}

	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return "", errors.New("Format UUID penugasan tidak valid")
	}

	existing, err := h.Repo.GetUserAbacByUuid(ctx, parsedUUID)
	if err != nil || existing == nil {
		return "", errors.New("Data penugasan tidak ditemukan")
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
	if cmd.IdFakultasUnit != nil && *cmd.IdFakultasUnit > 0 {
		targetFuID = cmd.IdFakultasUnit
		updates["id_fakultas_unit"] = *cmd.IdFakultasUnit
	} else if strings.TrimSpace(cmd.UuidFakultasUnit) != "" {
		fuID, _, _, err := h.Repo.GetFakultasUnitInfo(ctx, nil, cmd.UuidFakultasUnit)
		if err == nil {
			targetFuID = &fuID
			updates["id_fakultas_unit"] = fuID
		}
	}

	if cmd.Role != "" {
		reqRole := strings.ToLower(strings.TrimSpace(cmd.Role))
		if reqRole != "auto" {
			if reqRole != "upms" && reqRole != "upmf" && reqRole != "upml" {
				return "", errors.New("Role penugasan tidak valid (hanya diperbolehkan UPMS, UPMF, atau UPML)")
			}
			if targetFuID == nil || *targetFuID == 0 {
				return "", errors.New("Target Fakultas/Unit/Prodi wajib dipilih")
			}
			_, fuType, fuNama, err := h.Repo.GetFakultasUnitInfo(ctx, targetFuID, "")
			if err == nil {
				lowType := strings.ToLower(strings.TrimSpace(fuType))
				lowNama := strings.ToLower(strings.TrimSpace(fuNama))
				isFak := lowType == "fakultas" || strings.Contains(lowNama, "(fakultas)") || strings.HasPrefix(lowNama, "fakultas")
				isProdi := lowType == "prodi" || strings.Contains(lowNama, "(prodi)") || strings.HasPrefix(lowNama, "prodi") || strings.HasPrefix(lowNama, "program studi")

				if reqRole == "upms" && !isProdi {
					return "", errors.New("Target untuk role UPMS harus berupa Program Studi (Prodi)")
				}
				if reqRole == "upmf" && !isFak {
					return "", errors.New("Target untuk role UPMF harus berupa Fakultas")
				}
				if reqRole == "upml" && isProdi {
					return "", errors.New("Target untuk role UPML harus berupa Fakultas atau Unit (bukan Program Studi)")
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
