package event

import (
	"context"

	commondomain "UnpakSiamida/common/domain"
)

type BeritaAcaraPdfRequestedEventHandler struct {
	Redis commondomain.IRedisStore
}

func NewBeritaAcaraPdfRequestedEventHandler(
	redis commondomain.IRedisStore,
) *BeritaAcaraPdfRequestedEventHandler {
	return &BeritaAcaraPdfRequestedEventHandler{
		Redis: redis,
	}
}

func (h *BeritaAcaraPdfRequestedEventHandler) Handle(
	ctx context.Context,
	event BeritaAcaraPdfRequestedEvent,
) error {
	key := "pdf_berita_acara:" + event.BeritaAcaraUUID.String() + ":" + event.Token

	return h.Redis.Set(ctx, key)
}
