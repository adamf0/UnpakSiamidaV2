package application

import (
	"context"
	"errors"
	"strings"

	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateTahunProkerCommandHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *UpdateTahunProkerCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateTahunProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	tahunprokerUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintahunproker.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING tahunproker
	// -------------------------
	existingTahunProker, err := h.Repo.GetByUuid(ctx, tahunprokerUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domaintahunproker.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaintahunproker.UpdateTahunProker(
		existingTahunProker,
		tahunprokerUUID,
		cmd.Tahun,
		cmd.Status,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedTahunProker := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedTahunProker); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return "", domaintahunproker.DuplicateData()
			}
		}
		if strings.Contains(err.Error(), "Duplicate entry") {
			return "", domaintahunproker.DuplicateData()
		}

		return "", err
	}

	return updatedTahunProker.UUID.String(), nil
}
