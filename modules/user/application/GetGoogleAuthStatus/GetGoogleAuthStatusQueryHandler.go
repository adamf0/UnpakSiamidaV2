package application

import (
	"context"
	"os"
	"strings"
)

type GetGoogleAuthStatusQueryHandler struct{}

func (h *GetGoogleAuthStatusQueryHandler) Handle(
	ctx context.Context,
	q GetGoogleAuthStatusQuery,
) (*GetGoogleAuthStatusResult, error) {
	mailUser := strings.TrimSpace(os.Getenv("MAIL_USERNAME"))
	mailPass := strings.TrimSpace(os.Getenv("MAIL_PASSWORD"))
	mailHost := strings.TrimSpace(os.Getenv("MAIL_HOST"))

	hasSMTP := mailUser != "" && mailPass != "" && mailHost != ""
	return &GetGoogleAuthStatusResult{
		HasSMTP:           hasSMTP,
		SMTPEmail:         mailUser,
		SMTPHost:          mailHost,
		PreferredProvider: "GOOGLE_SMTP",
		IsReady:           hasSMTP,
	}, nil
}
