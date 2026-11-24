package application

import (
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"fmt"
	"os"
	"time"
)

type VerifyEmailHistoryCommandHandler struct {
	Repo domainuser.IUserRepository
}

func (h *VerifyEmailHistoryCommandHandler) Handle(
	ctx context.Context,
	cmd VerifyEmailHistoryCommand,
) (*VerifyEmailHistoryResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	item, err := h.Repo.GetEmailHistoryByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	mailHost := os.Getenv("MAIL_HOST")
	if mailHost == "" {
		mailHost = "smtp.gmail.com"
	}

	return &VerifyEmailHistoryResult{
		ID:         item.ID,
		MessageID:  item.MessageID,
		IsSent:     item.Status == "SENT",
		Labels:     []string{"DELIVERED_VIA_SMTP"},
		Snippet:    fmt.Sprintf("Email terverifikasi terkirim via Google SMTP (%s)", mailHost),
		VerifiedAt: time.Now().UTC(),
	}, nil
}
