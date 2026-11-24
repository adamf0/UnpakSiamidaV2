package application

import (
	domainuser "UnpakSiamida/modules/user/domain"
)

type GetEmailStatusQuery struct {
	OutboxID   string
	VerifyLive bool
}

type EmailStatusItem struct {
	domainuser.HistoryEmail
	GoogleVerified *bool    `json:"google_verified,omitempty"`
	GoogleLabels   []string `json:"google_labels,omitempty"`
}

type GetEmailStatusResult struct {
	OutboxID    string            `json:"outbox_id"`
	Total       int               `json:"total"`
	Sent        int               `json:"sent"`
	Failed      int               `json:"failed"`
	Pending     int               `json:"pending"`
	Sending     int               `json:"sending"`
	IsCompleted bool              `json:"is_completed"`
	Items       []EmailStatusItem `json:"items"`
}
