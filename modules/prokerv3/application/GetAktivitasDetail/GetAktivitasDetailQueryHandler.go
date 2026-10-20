package GetAktivitasDetail

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type GetAktivitasDetailQueryHandler struct {
	Repo domain.IProkerV3Repository
	DB   *gorm.DB
}

func (h *GetAktivitasDetailQueryHandler) Handle(ctx context.Context, query GetAktivitasDetailQuery) (*domain.AktivitasDetailResult, error) {
	act, err := h.Repo.GetAktivitasByUuid(ctx, query.UUID)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, errors.New("Aktivitas tidak ditemukan")
	}

	var proker *domain.ProkerV3Default
	if act.IdMatrikProker != nil && *act.IdMatrikProker > 0 && h.DB != nil {
		var mp domain.ProkerV3
		if err := h.DB.WithContext(ctx).Where("id = ?", *act.IdMatrikProker).First(&mp).Error; err == nil {
			proker, _ = h.Repo.GetDefaultByUuid(ctx, mp.UUID)
		}
	}

	return &domain.AktivitasDetailResult{
		Aktivitas: act,
		Proker:    proker,
		Timelines: act.Timelines,
		Data:      act.Timelines,
	}, nil
}
