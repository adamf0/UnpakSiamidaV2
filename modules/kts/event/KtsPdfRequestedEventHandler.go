package event

import (
	"context"

	commondomain "UnpakSiamida/common/domain"
)

type KtsPdfRequestedEventHandler struct {
	Redis commondomain.IRedisStore
}

func NewKtsPdfRequestedEventHandler(
	redis commondomain.IRedisStore,
) *KtsPdfRequestedEventHandler {
	return &KtsPdfRequestedEventHandler{
		Redis: redis,
	}
}

func (h *KtsPdfRequestedEventHandler) Handle(
	ctx context.Context,
	event KtsPdfRequestedEvent,
) error {
	key := "pdf_kts:" + event.KtsUUID.String() + ":" + event.Token

	return h.Redis.Set(ctx, key)
}
