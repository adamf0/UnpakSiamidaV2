package DeleteTimelineFileV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type DeleteTimelineFileV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *DeleteTimelineFileV3CommandHandler) Handle(ctx context.Context, cmd DeleteTimelineFileV3Command) (string, error) {
	fileUUID, err := uuid.Parse(cmd.FileUUID)
	if err != nil {
		return "", errors.New("Invalid file UUID")
	}

	fileObj, err := h.Repo.GetTimelineFileByUuid(ctx, fileUUID)
	if err != nil {
		return "", err
	}
	if fileObj == nil {
		return "", errors.New("File tidak ditemukan")
	}

	if err := h.Repo.DeleteTimelineFile(ctx, fileUUID); err != nil {
		return "", err
	}

	if fileObj.File != "" {
		_ = os.Remove(fileObj.File)
		if fileObj.FileURL != "" {
			_ = os.Remove("." + fileObj.FileURL)
		}
		_ = os.Remove(filepath.Join("./upload/proker_v3", fileObj.File))
	}

	return cmd.FileUUID, nil
}
