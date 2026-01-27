package application

import (
	"context"
	"errors"
	"strings"

	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"
	"time"

	"github.com/go-sql-driver/mysql"
)

type CreateTahunProkerCommandHandler struct {
	Repo domaintahunproker.ITahunProkerRepository
}

func (h *CreateTahunProkerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateTahunProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := domaintahunproker.NewTahunProker(
		cmd.Tahun,
		cmd.Status,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createTahunProker := result.Value
	if err := h.Repo.Create(ctx, createTahunProker); err != nil {
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

	return result.Value.UUID.String(), nil
}
