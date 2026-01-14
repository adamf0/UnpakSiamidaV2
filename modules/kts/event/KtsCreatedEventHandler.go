package event

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"context"
)

type KtsCreatedEventHandler struct {
	Telegram *commoninfra.TelegramClient
}

func NewKtsCreatedEventHandler(
	tg *commoninfra.TelegramClient,
) *KtsCreatedEventHandler {
	return &KtsCreatedEventHandler{
		Telegram: tg,
	}
}

func (h KtsCreatedEventHandler) Handle(
	ctx context.Context,
	event KtsCreatedEvent,
) error {
	// godump.Dump(event)
	msg := RenderKtsCreatedTemplate(event)

	if err := h.Telegram.SendHTML(msg); err != nil {
		return err
	}

	return nil
}
