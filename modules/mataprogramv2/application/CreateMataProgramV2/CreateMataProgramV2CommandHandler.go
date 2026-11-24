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
	if len(cmd.TemplateRenstraIDs) == 0 {
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
