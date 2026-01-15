package event

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"context"
)

type KtsUpdatedEventHandler struct {
	Telegram commoninfra.TelegramSender
}

func NewKtsUpdatedEventHandler(
	tg commoninfra.TelegramSender,
) *KtsUpdatedEventHandler {
	return &KtsUpdatedEventHandler{
		Telegram: tg,
	}
}

func (h KtsUpdatedEventHandler) Handle(
	ctx context.Context,
	event KtsUpdatedEvent,
) error {
	// godump.Dump(event)
	msg := RenderKtsUpdatedTemplate(event)

	if err := h.Telegram.SendHTML(msg); err != nil {
		return err
	}

	return nil
}
