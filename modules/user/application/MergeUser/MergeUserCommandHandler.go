package application

import (
	commondomain "UnpakSiamida/common/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MergeUserCommandHandler struct {
	Repo domainuser.IUserRepository
}

func (h *MergeUserCommandHandler) Handle(
	ctx context.Context,
	cmd MergeUserCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var primaryID uint
	if cmd.PrimaryID != nil && *cmd.PrimaryID > 0 {
		primaryID = *cmd.PrimaryID
	} else if strings.TrimSpace(cmd.PrimaryUUID) != "" {
		pUUID, err := uuid.Parse(strings.TrimSpace(cmd.PrimaryUUID))
		if err != nil {
			return "", domainuser.InvalidParsing("Primary UUID")
		}
		pUser, err := h.Repo.GetByUuid(ctx, pUUID)
		if err != nil || pUser == nil {
			return "", domainuser.NotFound("Primary User")
		}
		primaryID = pUser.ID
	} else {
		return "", commondomain.FailureError("User.Merge", "Akun primary wajib dipilih")
	}

	secondaryMap := make(map[uint]bool)
	for _, id := range cmd.SecondaryIDs {
		if id > 0 && id != primaryID {
			secondaryMap[id] = true
		}
	}

	for _, rawUUID := range cmd.SecondaryUUIDs {
		trimmed := strings.TrimSpace(rawUUID)
		if trimmed == "" {
			continue
		}
		sUUID, err := uuid.Parse(trimmed)
		if err != nil {
			continue
		}
		sUser, err := h.Repo.GetByUuid(ctx, sUUID)
		if err == nil && sUser != nil && sUser.ID != primaryID {
			secondaryMap[sUser.ID] = true
		}
	}

	if len(secondaryMap) == 0 {
		// Jika tidak ada secondary, reset id_user_merge milik primary menjadi null (unmerge)
		if err := h.Repo.MergeUser(ctx, primaryID, []uint{}); err != nil {
			return "", err
		}
		return "Status merge akun berhasil direset (unmerge)", nil
	}

	secondaryIDs := make([]uint, 0, len(secondaryMap))
	for id := range secondaryMap {
		secondaryIDs = append(secondaryIDs, id)
	}

	if err := h.Repo.MergeUser(ctx, primaryID, secondaryIDs); err != nil {
		return "", err
	}

	return "Data user berhasil digabung", nil
}
