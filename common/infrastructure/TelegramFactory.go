package infrastructure

import (
	"errors"
	"os"
)

type DefaultTelegramFactory struct {
	UseFake bool
}

func (f *DefaultTelegramFactory) Create() (TelegramSender, error) {
	if f.UseFake {
		return NewFakeTelegramClient(), nil
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if token == "" || chatID == "" {
		return nil, errors.New("telegram env not configured")
	}

	return NewTelegramClient(token, chatID), nil
}
