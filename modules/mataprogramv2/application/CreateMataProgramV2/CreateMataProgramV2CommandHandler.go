package CreateMataProgramV2

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"fmt"
	"strings"
)

type CreateMataProgramV2CommandHandler struct {
	Repo domain.IMataProgramV2Repository
}

func (h *CreateMataProgramV2CommandHandler) Handle(ctx context.Context, cmd CreateMataProgramV2Command) (*CreateMataProgramV2Result, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	mataProgram := strings.TrimSpace(cmd.MataProgram)
	if mataProgram == "" {
		return nil, domain.MataProgramRequired()
	}
	if len(cmd.Targets) == 0 && len(cmd.TemplateRenstraIDs) == 0 {
		return nil, domain.TargetsRequired()
	}

	tahunStr := strings.TrimSpace(cmd.Tahun)
	if tahunStr == "" {
		tahunStr = "2027"
	}

	tahunID, err := h.Repo.ResolveOrCreateTahunID(ctx, tahunStr)
	if err != nil {
		return nil, err
	}
	if tahunID == 0 {
		tahunID = 1
	}

	// Mode V3: Simpan relasi ke mata_programv2 dan mata_programv2_template_renstra
	if len(cmd.TemplateRenstraIDs) > 0 {
		savedUUID, savedCount, err := h.Repo.SaveMataProgramV2WithTemplates(ctx, tahunID, mataProgram, cmd.TemplateRenstraIDs)
		if err != nil {
			return nil, err
		}
		return &CreateMataProgramV2Result{
			Status:  "success",
			Message: fmt.Sprintf("Berhasil menyimpan mata program dengan %d template renstra", savedCount),
			UUID:    savedUUID,
			Count:   int(savedCount),
		}, nil
	}

	savedCount := 0
	var firstUuid string

	for _, rawTarget := range cmd.Targets {
		targetStr := strings.TrimSpace(rawTarget)
		if targetStr == "" || targetStr == "0" || targetStr == "null" {
			continue
		}

		fuID, err := h.Repo.ResolveFakultasUnitID(ctx, targetStr)
		if err != nil || fuID == 0 {
			continue
		}

		savedUUID, ok, err := h.Repo.UpsertOrRestoreAssignment(ctx, tahunID, fuID, mataProgram)
		if err != nil || !ok {
			continue
		}
		if firstUuid == "" {
			firstUuid = savedUUID
		}
		savedCount++
	}

	if savedCount == 0 {
		return nil, domain.InvalidTargets()
	}

	return &CreateMataProgramV2Result{
		Status:  "success",
		Message: fmt.Sprintf("Berhasil menyimpan %d target mata program", savedCount),
		UUID:    firstUuid,
		Count:   savedCount,
	}, nil
}
