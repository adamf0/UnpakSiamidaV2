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
			if targetFuID != nil && *targetFuID > 0 {
				_, fuType, fuNama, err := h.Repo.GetFakultasUnitInfo(ctx, targetFuID, "")
				if err == nil {
					isFak := fuType == "fakultas" || strings.Contains(strings.ToLower(fuNama), "(fakultas)") || strings.HasPrefix(strings.ToLower(fuNama), "fakultas")
					if isFak && reqRole == "upms" {
						return "", errors.New("Target Fakultas hanya dapat memilih role UPML atau UPMF")
					}
					if !isFak && reqRole == "upmf" {
						return "", errors.New("Target non-Fakultas hanya dapat memilih role UPML atau UPMS")
					}
				}
			} else if reqRole == "upmf" {
				return "", errors.New("Target Fakultas wajib dipilih untuk role UPMF")
			}
			updates["role"] = reqRole
		}
	}

	if err := h.Repo.UpdateUserAbac(ctx, existing.ID, updates); err != nil {
		return "", err
	}

	return "Penugasan role berhasil diperbarui", nil
}
