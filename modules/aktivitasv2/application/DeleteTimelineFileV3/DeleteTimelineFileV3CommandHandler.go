package DeleteTimelineFileV3

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type DeleteTimelineFileV3CommandHandler struct {
	Repo domain.IAktivitasV2Repository
}

func (h *DeleteTimelineFileV3CommandHandler) Handle(ctx context.Context, cmd DeleteTimelineFileV3Command) (string, error) {
	fileUUID, err := uuid.Parse(cmd.FileUUID)
	if err != nil {
		return "", domain.InvalidFileUuid()
	}

	fileObj, err := h.Repo.GetTimelineFileByUuid(ctx, fileUUID)
	if err != nil {
		return "", err
	}
	if fileObj == nil {
		return "", domain.FileNotFound()
	}

	if err := h.Repo.DeleteTimelineFile(ctx, fileUUID); err != nil {
		return "", err
	}

	isLink := strings.EqualFold(strings.TrimSpace(fileObj.FileType), "link") ||
		strings.HasPrefix(strings.ToLower(strings.TrimSpace(fileObj.FileURL)), "http://") ||
		strings.HasPrefix(strings.ToLower(strings.TrimSpace(fileObj.FileURL)), "https://") ||
		strings.HasPrefix(strings.ToLower(strings.TrimSpace(fileObj.File)), "http://") ||
		strings.HasPrefix(strings.ToLower(strings.TrimSpace(fileObj.File)), "https://")

	if !isLink && fileObj.File != "" {
		_ = os.Remove(fileObj.File)
		if fileObj.FileURL != "" {
			_ = os.Remove("." + fileObj.FileURL)
		}
		_ = os.Remove(filepath.Join("./upload/proker_v3", fileObj.File))
	}

	return cmd.FileUUID, nil
}
