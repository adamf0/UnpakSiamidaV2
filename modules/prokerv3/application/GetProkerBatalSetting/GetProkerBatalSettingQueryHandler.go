package GetProkerBatalSetting

import (
	"context"
	"strconv"
	"strings"
	"time"

	domain "UnpakSiamida/modules/prokerv3/domain"
)

type GetProkerBatalSettingQueryHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *GetProkerBatalSettingQueryHandler) Handle(ctx context.Context, query GetProkerBatalSettingQuery) (*GetProkerBatalSettingResult, error) {
	if h.Repo == nil {
		return nil, domain.RepoNotInitialized()
	}

	targetUnitID := query.TargetUnitID

	if query.TargetUnitUUID != "" {
		uid, err := h.Repo.GetFakultasUnitIDByUUIDOrName(ctx, query.TargetUnitUUID)
		if err == nil && uid != nil && *uid > 0 {
			targetUnitID = uid
		}
	}

	if targetUnitID == nil && !strings.EqualFold(query.UserLevel, "admin") && query.UserUnit != "" {
		if id, err := strconv.Atoi(query.UserUnit); err == nil && id > 0 {
			targetUnitID = &id
		} else {
			uid, err := h.Repo.GetFakultasUnitIDByUUIDOrName(ctx, query.UserUnit)
			if err == nil && uid != nil && *uid > 0 {
				targetUnitID = uid
			}
		}
	}

	status, err := h.Repo.GetProkerBatalSettingStatus(ctx, time.Now(), targetUnitID)
	if err != nil {
		return nil, err
	}

	return status, nil
}
