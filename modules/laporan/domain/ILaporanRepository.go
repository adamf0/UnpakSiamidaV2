package domain

import (
	"context"

	"github.com/google/uuid"
)

type ILaporanRepository interface {
	GetMonitoringByTargetTahun(
		ctx context.Context,
		uuidFakultasUnit uuid.UUID,
		uuidTahun uuid.UUID,
		page, limit *int,
	) ([]MonitoringProker, int64, error)

	GetMonitoringIndikatorByIndikatorTahun(
		ctx context.Context,
		uuidIndikator uuid.UUID,
		tahun string,
		page, limit *int,
	) ([]MonitoringIndikator, int64, error)
}
