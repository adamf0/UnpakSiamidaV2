package event

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"UnpakSiamida/common/infrastructure/emailhub"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

const EmailAttachmentBaseDir = "./upload/email_attachments"

type loadedEmailAttachment struct {
	FileName    string
	ContentType string
	Data        []byte
}

type UserSendAccountEmailEventHandler struct {
	DB       *gorm.DB
	Telegram commoninfra.TelegramSender
}

func NewUserSendAccountEmailEventHandler(
	db *gorm.DB,
	tg commoninfra.TelegramSender,
) *UserSendAccountEmailEventHandler {
	return &UserSendAccountEmailEventHandler{
		DB:       db,
		Telegram: tg,
	}
}

func (h UserSendAccountEmailEventHandler) Handle(
	ctx context.Context,
	event UserSendAccountEmailEvent,
) error {
	loadedAttachments, batchDirs, loadErr := loadEmailAttachments(event.Attachments)
	if loadErr != nil {
		log.Printf("[Outbox:UserAccountEmail] Failed loading attachments: %v\n", loadErr)
		return loadErr
	}

	successCount := 0
	var failedList []string

	mailHost := strings.TrimSpace(os.Getenv("MAIL_HOST"))
	mailPort := strings.TrimSpace(os.Getenv("MAIL_PORT"))
	mailUser := strings.TrimSpace(os.Getenv("MAIL_USERNAME"))
	mailPass := strings.TrimSpace(os.Getenv("MAIL_PASSWORD"))
	mailFrom := strings.TrimSpace(os.Getenv("MAIL_FROM_ADDRESS"))
	mailFromName := strings.TrimSpace(os.Getenv("MAIL_FROM_NAME"))
	portalURL := strings.TrimSpace(os.Getenv("SIAMIDA_V3_URL"))

	if mailFrom == "" {
		mailFrom = mailUser
	}
	if mailFromName == "" {
		mailFromName = "LPMI Universitas Pakuan"
	}
	if mailPort == "" {
		mailPort = "587"
	}
	if portalURL == "" {
		portalURL = "https://siamida.unpak.ac.id"
	}

	subject := "Informasi Akun SIAMIDA V3 - Universitas Pakuan"

	for _, user := range event.Users {
		email := strings.TrimSpace(user.Email)
		if email == "" {
			log.Printf("[Outbox:UserAccountEmail] Skip user %s: email is empty\n", user.Username)
			continue
		}

		plainText := RenderUserAccountEmailPlainText(user, event.Attachments...)
		htmlText := RenderUserAccountEmailHTML(user, portalURL, event.Attachments...)
		rawMsg := buildAccountEmailMIME(mailFromName, mailFrom, email, subject, plainText, htmlText, loadedAttachments)

		outboxID := event.EventID.String()
		userUUIDStr := user.UserUUID.String()

		// Tandai status history_email menjadi SENDING
		if h.DB != nil {
			h.updateHistoryEmailStatus(ctx, outboxID, userUUIDStr, "SENDING", "GOOGLE_SMTP", "", "", nil)
		}

		var sendErr error
		var providerUsed string
		var googleMessageID string
		var googleThreadID string

		// Google SMTP via App Password (Jalur Utama - 100% Headless tanpa OAuth/Token)
		if mailHost != "" && mailUser != "" && mailPass != "" {
			smtpErr := sendViaSMTPFallback(mailHost, mailPort, mailUser, mailPass, mailFrom, email, rawMsg)
			if smtpErr == nil {
				providerUsed = "GOOGLE_SMTP"
				googleMessageID = fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), userUUIDStr[:8], mailHost)
				log.Printf("[Outbox:UserAccountEmail] Email sent successfully via Google SMTP to %s <%s>\n", user.Username, email)
			} else {
				log.Printf("[Outbox:UserAccountEmail] Google SMTP gagal: %v\n", smtpErr)
				sendErr = fmt.Errorf("Google SMTP error: %v", smtpErr)
			}
		} else {
			sendErr = fmt.Errorf("kredensial Google SMTP belum dikonfigurasi di .env")
		}

		now := time.Now().UTC()
		if sendErr != nil {
			log.Printf("[Outbox:UserAccountEmail] Error sending to %s (%s): %v\n", user.Username, email, sendErr)
			failedList = append(failedList, fmt.Sprintf("%s (%v)", user.Username, sendErr))

			errMsg := sendErr.Error()
			if h.DB != nil {
				h.updateHistoryEmailStatus(ctx, outboxID, userUUIDStr, "FAILED", "", "", "", &errMsg)
			}
		} else {
			successCount++
			if h.DB != nil {
				h.updateHistoryEmailStatus(ctx, outboxID, userUUIDStr, "SENT", providerUsed, googleMessageID, googleThreadID, nil, &now)
			}
		}
	}

	if len(failedList) == 0 && len(batchDirs) > 0 {
		for dir := range batchDirs {
			if err := os.RemoveAll(dir); err != nil {
				log.Printf("[Outbox:UserAccountEmail] Cleanup attachment dir %s failed: %v\n", dir, err)
			}
		}
	}

	if h.Telegram != nil && len(event.Users) > 0 {
		msg := RenderUserAccountEmailTelegram(event)
		if len(failedList) > 0 {
			msg += fmt.Sprintf("\n⚠️ Gagal kirim (%d): %s", len(failedList), strings.Join(failedList, ", "))
		}
		if err := h.Telegram.SendHTML(msg); err != nil {
			fmt.Printf("[Outbox:UserAccountEmail] Telegram notify error: %v\n", err)
		}
	}

	fmt.Printf("[Outbox:UserAccountEmail] Selesai memproses email akun: %d berhasil, %d gagal dari %d total\n",
		successCount, len(failedList), len(event.Users))

	return nil
}

func (h *UserSendAccountEmailEventHandler) updateHistoryEmailStatus(
	ctx context.Context,
	outboxID, userUUID, status, provider, messageID, threadID string,
	errMsg *string,
	sentAt ...*time.Time,
) {
	if h.DB == nil {
		return
	}
	now := time.Now().UTC()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}
	if provider != "" {
		updates["provider"] = provider
	}
	if messageID != "" {
		updates["message_id"] = messageID
	}
	if threadID != "" {
		updates["thread_id"] = threadID
	}
	if errMsg != nil {
		updates["error_message"] = *errMsg
	}
	var actualSentAt *time.Time
	if len(sentAt) > 0 && sentAt[0] != nil {
		actualSentAt = sentAt[0]
		updates["sent_at"] = actualSentAt
	}

	_ = h.DB.WithContext(ctx).Table("history_email").
		Where("outbox_id = ? AND user_uuid = ?", outboxID, userUUID).
		Updates(updates).Error

	// Ambil detail data untuk siaran live broadcast (WebSocket & SSE)
	var row struct {
		ID             uint   `gorm:"column:id"`
		RecipientEmail string `gorm:"column:recipient_email"`
		Username       string `gorm:"column:username"`
		Subject        string `gorm:"column:subject"`
		Provider       string `gorm:"column:provider"`
		MessageID      string `gorm:"column:message_id"`
	}
	_ = h.DB.WithContext(ctx).Table("history_email").
		Select("id, recipient_email, username, subject, provider, message_id").
		Where("outbox_id = ? AND user_uuid = ?", outboxID, userUUID).
		Limit(1).Scan(&row).Error

	emailhub.GetHub().Broadcast(emailhub.EmailStatusEvent{
		Type:           "EMAIL_STATUS_UPDATE",
		ID:             row.ID,
		OutboxID:       outboxID,
		UserUUID:       userUUID,
		RecipientEmail: row.RecipientEmail,
		Username:       row.Username,
		Subject:        row.Subject,
		Status:         status,
		Provider:       row.Provider,
		MessageID:      row.MessageID,
		ErrorMessage:   errMsg,
		SentAt:         actualSentAt,
		UpdatedAt:      now,
	})
}

func sendViaSMTPFallback(host, port, user, pass, from, to string, rawMsg []byte) error {
	if host == "" || user == "" || pass == "" {
		return fmt.Errorf("SMTP configuration missing in environment")
	}

	if port == "465" {
		err := sendViaSSL(host, "465", user, pass, from, to, rawMsg)
		if err == nil {
			return nil
		}
		log.Printf("[Outbox:UserAccountEmail] SMTP Port 465 failed (%v), fallback to 587...\n", err)
		return sendViaSTARTTLS(host, "587", user, pass, from, to, rawMsg)
	}

	err := sendViaSTARTTLS(host, port, user, pass, from, to, rawMsg)
	if err == nil {
		return nil
	}
	if port == "587" {
		log.Printf("[Outbox:UserAccountEmail] SMTP Port 587 failed (%v), fallback to 465...\n", err)
		return sendViaSSL(host, "465", user, pass, from, to, rawMsg)
	}
	return err
}

func sanitizeMIMEHeaderValue(v string) string {
	v = strings.NewReplacer("\r", "", "\n", "", `"`, "", "\\", "").Replace(strings.TrimSpace(v))
	var b strings.Builder
	for _, r := range v {
		if r >= 32 && r != 127 {
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}

func loadEmailAttachments(attachments []UserAccountEmailAttachment) ([]loadedEmailAttachment, map[string]struct{}, error) {
	if len(attachments) == 0 {
		return nil, nil, nil
	}

	baseAbs, err := filepath.Abs(EmailAttachmentBaseDir)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve attachment base dir: %w", err)
	}
	basePrefix := baseAbs + string(os.PathSeparator)

	loaded := make([]loadedEmailAttachment, 0, len(attachments))
	batchDirs := make(map[string]struct{})

	for _, att := range attachments {
		absPath, err := filepath.Abs(att.FilePath)
		if err != nil || !strings.HasPrefix(absPath, basePrefix) {
			return nil, nil, fmt.Errorf("invalid attachment path: %q", att.FilePath)
		}

		data, err := os.ReadFile(absPath)
		if err != nil {
			return nil, nil, fmt.Errorf("read attachment %q: %w", att.FileName, err)
		}

		safeName := sanitizeMIMEHeaderValue(filepath.Base(att.FileName))
		if safeName == "" {
			safeName = "lampiran"
		}
		ct := sanitizeMIMEHeaderValue(att.ContentType)
		if ct == "" {
			ct = "application/octet-stream"
		}

		loaded = append(loaded, loadedEmailAttachment{
			FileName:    safeName,
			ContentType: ct,
			Data:        data,
		})

		parentDir := filepath.Dir(absPath)
		if strings.HasPrefix(parentDir, basePrefix) {
			batchDirs[parentDir] = struct{}{}
		}
	}

	return loaded, batchDirs, nil
}

func writeWrappedBase64(b *strings.Builder, data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	const lineLen = 76
	for i := 0; i < len(encoded); i += lineLen {
		end := i + lineLen
		if end > len(encoded) {
			end = len(encoded)
		}
		b.WriteString(encoded[i:end])
		b.WriteString("\r\n")
	}
}

func buildAccountEmailMIME(
	mailFromName, mailFrom, toEmail, subject, plainText, htmlText string,
	attachments []loadedEmailAttachment,
) []byte {
	ts := time.Now().UnixNano()
	altBoundary := fmt.Sprintf("----=_Alt_%d", ts)

	var msgBuilder strings.Builder
	msgBuilder.WriteString(fmt.Sprintf("From: %s <%s>\r\n", sanitizeMIMEHeaderValue(mailFromName), sanitizeMIMEHeaderValue(mailFrom)))
	msgBuilder.WriteString(fmt.Sprintf("To: %s\r\n", sanitizeMIMEHeaderValue(toEmail)))
	msgBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", sanitizeMIMEHeaderValue(subject)))
	msgBuilder.WriteString("MIME-Version: 1.0\r\n")

	if len(attachments) == 0 {
		msgBuilder.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", altBoundary))
		msgBuilder.WriteString("\r\n")
		writeAlternativeParts(&msgBuilder, altBoundary, plainText, htmlText)
		return []byte(msgBuilder.String())
	}

	mixedBoundary := fmt.Sprintf("----=_Mixed_%d", ts)
	msgBuilder.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", mixedBoundary))
	msgBuilder.WriteString("\r\n")

	// Body part (multipart/alternative)
	msgBuilder.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
	msgBuilder.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", altBoundary))
	msgBuilder.WriteString("\r\n")
	writeAlternativeParts(&msgBuilder, altBoundary, plainText, htmlText)
	msgBuilder.WriteString("\r\n")

	// Attachment parts
	for _, att := range attachments {
		safeName := sanitizeMIMEHeaderValue(att.FileName)
		if safeName == "" {
			safeName = "lampiran"
		}
		safeCT := sanitizeMIMEHeaderValue(att.ContentType)
		if safeCT == "" {
			safeCT = "application/octet-stream"
		}

		msgBuilder.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
		msgBuilder.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", safeCT, safeName))
		msgBuilder.WriteString("Content-Transfer-Encoding: base64\r\n")
		msgBuilder.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", safeName))
		msgBuilder.WriteString("\r\n")
		writeWrappedBase64(&msgBuilder, att.Data)
		msgBuilder.WriteString("\r\n")
	}

	msgBuilder.WriteString(fmt.Sprintf("--%s--\r\n", mixedBoundary))
	return []byte(msgBuilder.String())
}

func writeAlternativeParts(b *strings.Builder, boundary, plainText, htmlText string) {
	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	b.WriteString("\r\n")
	b.WriteString(plainText)
	b.WriteString("\r\n\r\n")

	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	b.WriteString("\r\n")
	b.WriteString(htmlText)
	b.WriteString("\r\n\r\n")

	b.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
}

func sendAccountCredentialEmail(
	recipient UserAccountEmailRecipient,
	attachmentMeta []UserAccountEmailAttachment,
	loadedAttachments []loadedEmailAttachment,
) error {
	mailHost := strings.TrimSpace(os.Getenv("MAIL_HOST"))
	mailPort := strings.TrimSpace(os.Getenv("MAIL_PORT"))
	mailUser := strings.TrimSpace(os.Getenv("MAIL_USERNAME"))
	mailPass := strings.TrimSpace(os.Getenv("MAIL_PASSWORD"))
	mailFrom := strings.TrimSpace(os.Getenv("MAIL_FROM_ADDRESS"))
	mailFromName := strings.TrimSpace(os.Getenv("MAIL_FROM_NAME"))
	portalURL := strings.TrimSpace(os.Getenv("SIAMIDA_V3_URL"))

	if mailFrom == "" {
		mailFrom = mailUser
	}
	if mailFromName == "" {
		mailFromName = "LPMI Universitas Pakuan"
	}
	if mailPort == "" {
		mailPort = "587"
	}
	if portalURL == "" {
		portalURL = "https://siamida.unpak.ac.id"
	}

	if mailHost == "" || mailUser == "" || mailPass == "" {
		log.Printf("[Outbox:UserAccountEmail] SMTP configuration missing, skipping mail send\n")
		return nil
	}

	subject := "Informasi Akun SIAMIDA V3 - Universitas Pakuan"
	plainText := RenderUserAccountEmailPlainText(recipient, attachmentMeta...)
	htmlText := RenderUserAccountEmailHTML(recipient, portalURL, attachmentMeta...)
	rawMsg := buildAccountEmailMIME(mailFromName, mailFrom, recipient.Email, subject, plainText, htmlText, loadedAttachments)

	// Jika port 465, coba via Direct SSL. Jika gagal/reset/timeout, otomatis fallback ke port 587 (STARTTLS).
	if mailPort == "465" {
		err := sendViaSSL(mailHost, "465", mailUser, mailPass, mailFrom, recipient.Email, rawMsg)
		if err == nil {
			return nil
		}
		log.Printf("[Outbox:UserAccountEmail] Port 465 failed (%v), falling back to port 587 (STARTTLS)...\n", err)
		return sendViaSTARTTLS(mailHost, "587", mailUser, mailPass, mailFrom, recipient.Email, rawMsg)
	}

	// Port 587 (atau port lainnya) via STARTTLS. Jika gagal, coba fallback ke port 465 (SSL).
	err := sendViaSTARTTLS(mailHost, mailPort, mailUser, mailPass, mailFrom, recipient.Email, rawMsg)
	if err == nil {
		return nil
	}
	if mailPort == "587" {
		log.Printf("[Outbox:UserAccountEmail] Port 587 failed (%v), falling back to port 465 (SSL)...\n", err)
		return sendViaSSL(mailHost, "465", mailUser, mailPass, mailFrom, recipient.Email, rawMsg)
	}
	return err
}

func sendViaSSL(host, port, user, pass, from, to string, msg []byte) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	tlsConfig := &tls.Config{
		ServerName: host,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls.Dial error to %s: %w", addr, err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp.NewClient error: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", user, pass, host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from error: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt error: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data writer error: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close writer error: %w", err)
	}

	return client.Quit()
}

func sendViaSTARTTLS(host, port, user, pass, from, to string, msg []byte) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("tcp dial error to %s: %w", addr, err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp.NewClient error: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: host}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls error: %w", err)
		}
	}

	auth := smtp.PlainAuth("", user, pass, host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from error: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt error: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data writer error: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close writer error: %w", err)
	}

	return client.Quit()
}
