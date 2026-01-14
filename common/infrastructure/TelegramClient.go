package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TelegramClient struct {
	Token  string
	ChatID string
	Client *http.Client
}

func NewTelegramClient(token, chatID string) *TelegramClient {
	return &TelegramClient{
		Token:  token,
		ChatID: chatID,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func (c *TelegramClient) SendHTML(message string) error {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		c.Token,
	)

	payload := sendMessageRequest{
		ChatID:    c.ChatID,
		Text:      message,
		ParseMode: "HTML",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram error status: %s", resp.Status)
	}

	return nil
}
