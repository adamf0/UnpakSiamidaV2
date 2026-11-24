package UpdateMataProgramV2

import (
	domain "UnpakSiamida/modules/mataprogramv2/domain"
	"context"
	"strings"
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

	original, err := h.Repo.GetRecordByUuid(ctx, rawUuid, false)
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, domain.NotFound()
	}

	var finalTahunID uint
	if original.IdMasterTahun != nil {
		finalTahunID = *original.IdMasterTahun
	}
	tahunStr := strings.TrimSpace(cmd.Tahun)
	if tahunStr != "" {
		resolvedID, err := h.Repo.ResolveOrCreateTahunID(ctx, tahunStr)
		if err != nil {
			return nil, err
		}
		if resolvedID > 0 {
			finalTahunID = resolvedID
		}
	}

	mataProgram := strings.TrimSpace(cmd.MataProgram)
	newMataProg := original.MataProgram
	if mataProgram != "" {
		newMataProg = mataProgram
	}

	replaceTemplates := cmd.HasTemplates || len(cmd.TemplateRenstraIDs) > 0
	if replaceTemplates && len(cmd.TemplateRenstraIDs) == 0 {
		return nil, domain.TargetsRequired()
	}

	if err := h.Repo.UpdateMataProgramV2WithTemplates(ctx, rawUuid, finalTahunID, newMataProg, cmd.TemplateRenstraIDs, replaceTemplates); err != nil {
		return nil, err
	}

	return &UpdateMataProgramV2Result{
		Status:  "success",
		Message: "Berhasil memperbarui mata program dan template renstra",
		UUID:    rawUuid,
	}, nil
}
