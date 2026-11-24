package application

import (
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"fmt"
	"strings"
	"time"
)

type GetEmailStatusQueryHandler struct {
	Repo domainuser.IUserRepository
}

func (h *GetEmailStatusQueryHandler) Handle(
	ctx context.Context,
	q GetEmailStatusQuery,
) (*GetEmailStatusResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	outboxID := strings.TrimSpace(q.OutboxID)
	if outboxID == "" {
		return nil, fmt.Errorf("parameter outbox_id wajib diisi")
	}

	items, err := h.Repo.GetEmailStatusByOutboxID(ctx, outboxID)
	if err != nil {
		return nil, err
	}

	var results []EmailStatusItem
	total := len(items)
	sent := 0
	failed := 0
	pending := 0
	sending := 0

	for _, it := range items {
		stItem := EmailStatusItem{HistoryEmail: it}
		switch it.Status {
		case "SENT":
			sent++
			if q.VerifyLive && it.MessageID != "" {
				isSent := true
				stItem.GoogleVerified = &isSent
				stItem.GoogleLabels = []string{"DELIVERED"}
			}
		case "FAILED":
			failed++
		case "SENDING":
			sending++
		default:
			pending++
		}
		results = append(results, stItem)
	}

	isCompleted := total > 0 && (sent+failed == total)

	return &GetEmailStatusResult{
		OutboxID:    outboxID,
		Total:       total,
		Sent:        sent,
		Failed:      failed,
		Pending:     pending,
		Sending:     sending,
		IsCompleted: isCompleted,
		Items:       results,
	}, nil
}
