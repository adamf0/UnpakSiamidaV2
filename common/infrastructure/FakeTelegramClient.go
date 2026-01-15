package infrastructure

type FakeTelegramClient struct {
	Messages []string
}

func NewFakeTelegramClient() TelegramSender {
	return &FakeTelegramClient{}
}

func (f *FakeTelegramClient) SendHTML(message string) error {
	f.Messages = append(f.Messages, message)
	return nil
}
