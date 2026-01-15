package event

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"context"
)

type UserUpdatedEventHandler struct {
	Telegram commoninfra.TelegramSender
}

func NewUserUpdatedEventHandler(
	tg commoninfra.TelegramSender,
) *UserUpdatedEventHandler {
	return &UserUpdatedEventHandler{
		Telegram: tg,
	}
}

func (h UserUpdatedEventHandler) Handle(
	ctx context.Context,
	event UserUpdatedEvent,
) error {
	// godump.Dump(event)
	msg := RenderUserUpdatedTemplate(event)

	if err := h.Telegram.SendHTML(msg); err != nil {
		return err
	}

	return nil
}
