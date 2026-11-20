package event

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type UserSendAccountEmailEventHandler struct {
	Telegram commoninfra.TelegramSender
}

func NewUserSendAccountEmailEventHandler(
	tg commoninfra.TelegramSender,
) *UserSendAccountEmailEventHandler {
	return &UserSendAccountEmailEventHandler{
		Telegram: tg,
	}
}

func (h UserSendAccountEmailEventHandler) Handle(
	ctx context.Context,
	event UserSendAccountEmailEvent,
) error {
	successCount := 0
	var failedList []string

	for _, user := range event.Users {
		email := strings.TrimSpace(user.Email)
		if email == "" {
			log.Printf("[Outbox:UserAccountEmail] Skip user %s: email is empty\n", user.Username)
			continue
		}

		if err := sendAccountCredentialEmail(user); err != nil {
			log.Printf("[Outbox:UserAccountEmail] Error sending to %s (%s): %v\n", user.Username, email, err)
			failedList = append(failedList, fmt.Sprintf("%s (%v)", user.Username, err))
		} else {
			successCount++
			log.Printf("[Outbox:UserAccountEmail] Email sent successfully to %s <%s>\n", user.Username, email)
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

func sendAccountCredentialEmail(recipient UserAccountEmailRecipient) error {
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
	plainText := RenderUserAccountEmailPlainText(recipient)
	htmlText := RenderUserAccountEmailHTML(recipient, portalURL)

	boundary := fmt.Sprintf("----=_Part_%d_%s", time.Now().UnixNano(), recipient.Username)

	var msgBuilder strings.Builder
	msgBuilder.WriteString(fmt.Sprintf("From: %s <%s>\r\n", mailFromName, mailFrom))
	msgBuilder.WriteString(fmt.Sprintf("To: %s\r\n", recipient.Email))
	msgBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msgBuilder.WriteString("MIME-Version: 1.0\r\n")
	msgBuilder.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
	msgBuilder.WriteString("\r\n")

	// Part 1: Text Plain
	msgBuilder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msgBuilder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	msgBuilder.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(plainText)
	msgBuilder.WriteString("\r\n\r\n")

	// Part 2: HTML
	msgBuilder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msgBuilder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msgBuilder.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(htmlText)
	msgBuilder.WriteString("\r\n\r\n")

	msgBuilder.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	rawMsg := []byte(msgBuilder.String())

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
