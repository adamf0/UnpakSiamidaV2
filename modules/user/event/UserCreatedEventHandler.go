package event

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"context"
)

type UserCreatedEventHandler struct {
	Telegram commoninfra.TelegramSender
}

func NewUserCreatedEventHandler(
	tg commoninfra.TelegramSender,
) *UserCreatedEventHandler {
	return &UserCreatedEventHandler{
		Telegram: tg,
	}
}

func (h UserCreatedEventHandler) Handle(
	ctx context.Context,
	event UserCreatedEvent,
) error {
	if h.Telegram != nil {
		// godump.Dump(event)
		msg := RenderUserCreatedTemplate(event)

		if err := h.Telegram.SendHTML(msg); err != nil {
			return err
		}
	}

	return nil
}
