package application

import (
	"context"

	domaingeneraterenstra "UnpakSiamida/modules/generaterenstra/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	helper "UnpakSiamida/common/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"errors"
)

type DeleteGenerateRenstraCommandHandler struct {
	Repo domaingeneraterenstra.IGenerateRenstraRepository
	RepoRenstra domainrenstra.IRenstraRepository
}

func (h *DeleteGenerateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd DeleteGenerateRenstraCommand,
) (string, error) {

	if !helper.IsValidTypeGenerate(cmd.Type) {
		return "", domaingeneraterenstra.InvalidType(cmd.Type)
	}

	uuidTarget, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaingeneraterenstra.InvalidUuid()
	}

	renstraUUID, err := uuid.Parse(cmd.UuidRenstra)
	if err != nil {
		return "", domaingeneraterenstra.InvalidRenstra()
	}

	existingRenstra, err := h.RepoRenstra.GetByUuid(ctx, renstraUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domainrenstra.NotFound(cmd.UuidRenstra)
		}
		return "", err
	}

	if cmd.Type=="renstra" {
		if err := h.Repo.ForceDeleteRenstraNilai(ctx, uuidTarget, existingRenstra.ID); err != nil {
			return "", err
		}
	} else{
		if err := h.Repo.ForceDeleteDokumenTambahan(ctx, uuidTarget, existingRenstra.ID); err != nil {
			return "", err
		}
	}

	return cmd.Uuid, nil
}
