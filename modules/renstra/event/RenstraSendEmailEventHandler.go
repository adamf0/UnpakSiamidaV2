package event

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"context"
	"fmt"
	"html"
	"net/smtp"
	"os"
	"strings"
)

type RenstraSendEmailEventHandler struct {
	Telegram commoninfra.TelegramSender
}

func NewRenstraSendEmailEventHandler(
	tg commoninfra.TelegramSender,
) *RenstraSendEmailEventHandler {
	return &RenstraSendEmailEventHandler{
		Telegram: tg,
	}
}

func (h RenstraSendEmailEventHandler) Handle(
	ctx context.Context,
	event RenstraSendEmailEvent,
) error {
	for _, r := range event.Recipients {
		if strings.TrimSpace(r.Email) == "" {
			continue
		}
		_ = sendRenstraAccessCodeEmail(r, event)
	}

	if h.Telegram != nil {
		msg := RenderRenstraSendEmailTemplate(event)
		if err := h.Telegram.SendHTML(msg); err != nil {
			// Do not block outbox completion if telegram fails in dev
			fmt.Printf("[Outbox:RenstraSendEmailEvent] telegram error: %v\n", err)
		}
	}

	fmt.Printf("[Outbox:RenstraSendEmailEvent] Processed kode akses email for %s (%s) -> %d recipient(s)\n",
		event.TargetAudit, event.Tahun, len(event.Recipients))
	return nil
}

func sendRenstraAccessCodeEmail(recipient RenstraEmailRecipient, event RenstraSendEmailEvent) error {
	mailHost := strings.TrimSpace(os.Getenv("MAIL_HOST"))
	mailPort := strings.TrimSpace(os.Getenv("MAIL_PORT"))
	mailUser := strings.TrimSpace(os.Getenv("MAIL_USERNAME"))
	mailPass := strings.TrimSpace(os.Getenv("MAIL_PASSWORD"))
	mailFrom := strings.TrimSpace(os.Getenv("MAIL_FROM_ADDRESS"))
	if mailFrom == "" {
		mailFrom = mailUser
	}
	if mailPort == "" {
		mailPort = "587"
	}

	if mailHost == "" || mailUser == "" || mailPass == "" {
		return nil
	}

	subject := "Pemberitahuan Kode Akses Audit Mutu Internal - LPMI UNPAK"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Notifikasi Kode Renstra</title>
</head>
<body style="font-family: Arial, sans-serif; color: #1e293b; line-height: 1.6;">
    <p>Halo <strong>%s</strong> (%s),</p>
    <p>Berikut adalah Kode Akses Jadwal Audit Mutu Internal (Renstra) untuk <strong>%s</strong> Tahun <strong>%s</strong>:</p>
    <h2 style="background: #f3e8ff; color: #6b21a8; display: inline-block; padding: 10px 20px; border-radius: 8px; letter-spacing: 2px;">%s</h2>
    <p>Gunakan kode akses ini untuk membuka instrumen audit pada sistem SIAMIDA LPMI Universitas Pakuan.</p>
    <br>
    <p>Terima kasih,<br><strong>LPMI Universitas Pakuan</strong></p>
</body>
</html>`,
		html.EscapeString(recipient.Name),
		html.EscapeString(recipient.Role),
		html.EscapeString(event.TargetAudit),
		html.EscapeString(event.Tahun),
		html.EscapeString(event.KodeAkses),
	)

	msg := strings.Join([]string{
		fmt.Sprintf("From: LPMI UNPAK <%s>", mailFrom),
		fmt.Sprintf("To: %s", recipient.Email),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", mailUser, mailPass, mailHost)
	return smtp.SendMail(fmt.Sprintf("%s:%s", mailHost, mailPort), auth, mailFrom, []string{recipient.Email}, []byte(msg))
}

func RenderRenstraSendEmailTemplate(event RenstraSendEmailEvent) string {
	var recLines []string
	for _, r := range event.Recipients {
		recLines = append(recLines, fmt.Sprintf("• <b>%s</b>: %s (%s)",
			html.EscapeString(r.Role),
			html.EscapeString(r.Name),
			html.EscapeString(r.Email),
		))
	}
	return fmt.Sprintf(
		"<b>[OUTBOX EMAIL KODE AKSES AUDIT]</b>\nUnit/Prodi: <b>%s</b>\nTahun: <b>%s</b>\nKode Akses: <code>%s</code>\nPenerima:\n%s",
		html.EscapeString(event.TargetAudit),
		html.EscapeString(event.Tahun),
		html.EscapeString(event.KodeAkses),
		strings.Join(recLines, "\n"),
	)
}
