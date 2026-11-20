package UpdateMataProgramV2

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"strconv"
	"strings"
	"time"
)

type UpdateMataProgramV2CommandHandler struct {
	Repo domain.IMataProgramV2Repository
}

func (h *UpdateMataProgramV2CommandHandler) Handle(ctx context.Context, cmd UpdateMataProgramV2Command) (*UpdateMataProgramV2Result, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	rawUuid := strings.TrimSpace(cmd.UUID)
	if rawUuid == "" {
		return nil, domain.UuidRequired()
	}

	scope := strings.TrimSpace(cmd.Scope)
	if scope == "" {
		scope = "group"
	}

	original, err := h.Repo.GetRecordByUuid(ctx, rawUuid, false)
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, domain.NotFound()
	}

	newTahunID := original.IdMasterTahun
	tahunStr := strings.TrimSpace(cmd.Tahun)
	if tahunStr != "" {
		if resolvedID, err := h.Repo.ResolveTahunID(ctx, tahunStr); err == nil && resolvedID > 0 {
			newTahunID = &resolvedID
		}
	}

	mataProgram := strings.TrimSpace(cmd.MataProgram)
	newMataProg := original.MataProgram
	if mataProgram != "" {
		newMataProg = mataProgram
	}

	if scope == "group" {
		var finalTahunID uint = 1
		if newTahunID != nil {
			finalTahunID = *newTahunID
		}

		if len(cmd.TemplateRenstraIDs) > 0 {
			if err := h.Repo.UpdateMataProgramV2WithTemplates(ctx, rawUuid, finalTahunID, newMataProg, cmd.TemplateRenstraIDs); err != nil {
				return nil, err
			}
			return &UpdateMataProgramV2Result{
				Status:  "success",
				Message: "Berhasil memperbarui mata program dan template renstra",
				UUID:    rawUuid,
			}, nil
		}

		updateMap := map[string]interface{}{}
		if newMataProg != original.MataProgram {
			updateMap["mata_program"] = newMataProg
		}
		if newTahunID != original.IdMasterTahun && newTahunID != nil {
			updateMap["id_master_tahun"] = *newTahunID
		}

		if len(updateMap) > 0 {
			if err := h.Repo.UpdateGroupFields(ctx, original.MataProgram, original.IdMasterTahun, updateMap); err != nil {
				return nil, err
			}
		}

		if len(cmd.Targets) > 0 {
			var targetUnitIDs []uint
			for _, tStr := range cmd.Targets {
				fuID, err := h.Repo.ResolveFakultasUnitID(ctx, strings.TrimSpace(tStr))
				if err == nil && fuID > 0 {
					targetUnitIDs = append(targetUnitIDs, fuID)
				}
			}

			if len(targetUnitIDs) > 0 {
				var finalTahunID uint = 1
				if newTahunID != nil {
					finalTahunID = *newTahunID
				}

				if err := h.Repo.SoftDeleteGroupUnitsNotIn(ctx, newMataProg, finalTahunID, targetUnitIDs, time.Now()); err != nil {
					return nil, err
				}

				for _, fuID := range targetUnitIDs {
					_, _, _ = h.Repo.UpsertOrRestoreAssignment(ctx, finalTahunID, fuID, newMataProg)
				}
			}
		}
	} else {
		updateMap := map[string]interface{}{}
		if mataProgram != "" {
			updateMap["mata_program"] = mataProgram
		}
		if newTahunID != nil {
			updateMap["id_master_tahun"] = *newTahunID
		}
		targetStr := strings.TrimSpace(cmd.IdFakultasUnit)
		if targetStr != "" {
			if id, err := strconv.Atoi(targetStr); err == nil && id > 0 {
				updateMap["id_fakultas_unit"] = id
			}
		}
		if len(updateMap) > 0 {
			if err := h.Repo.UpdateSingleByUuid(ctx, rawUuid, updateMap); err != nil {
				return nil, err
			}
		}
	}

	return &UpdateMataProgramV2Result{
		Status:  "success",
		Message: "Berhasil memperbarui data Mata Program V2",
		UUID:    rawUuid,
	}, nil
}
