package event

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"testing"
)

func TestBuildAccountEmailMIME_WithoutAttachments(t *testing.T) {
	raw := buildAccountEmailMIME(
		"LPMI Universitas Pakuan",
		"noreply@unpak.ac.id",
		"user@gmail.com",
		"Informasi Akun SIAMIDA V3",
		"plain body",
		"<p>html body</p>",
		nil,
	)

	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("failed to parse MIME message: %v", err)
	}
	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("failed to parse Content-Type: %v", err)
	}
	if mediaType != "multipart/alternative" {
		t.Fatalf("expected multipart/alternative, got %s", mediaType)
	}

	mr := multipart.NewReader(msg.Body, params["boundary"])
	partCount := 0
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed reading part: %v", err)
		}
		partCount++
		_, _ = io.ReadAll(p)
	}
	if partCount != 2 {
		t.Fatalf("expected 2 alternative parts, got %d", partCount)
	}
}

func TestBuildAccountEmailMIME_WithAttachmentsAndHeaderInjectionGuard(t *testing.T) {
	pdfBytes := []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\n%%EOF\n")
	txtBytes := []byte("Panduan singkat login SIAMIDA V3")

	raw := buildAccountEmailMIME(
		"LPMI Universitas Pakuan",
		"noreply@unpak.ac.id",
		"upmf_ilkom@gmail.com",
		"Informasi Akun SIAMIDA V3",
		"plain body",
		"<p>html body</p>",
		[]loadedEmailAttachment{
			{
				FileName:    "Panduan.pdf\r\nBcc: evil@attacker.test",
				ContentType: "application/pdf",
				Data:        pdfBytes,
			},
			{
				FileName:    "catatan.txt",
				ContentType: "text/plain",
				Data:        txtBytes,
			},
		},
	)

	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("failed to parse mixed MIME message: %v", err)
	}
	if bcc := msg.Header.Get("Bcc"); bcc != "" {
		t.Fatalf("unexpected injected Bcc header: %q", bcc)
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("failed to parse Content-Type: %v", err)
	}
	if mediaType != "multipart/mixed" {
		t.Fatalf("expected multipart/mixed, got %s", mediaType)
	}

	mr := multipart.NewReader(msg.Body, params["boundary"])
	var decodedPayloads [][]byte
	var filenames []string
	hasAlternative := false

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed reading mixed part: %v", err)
		}
		partMT, _, _ := mime.ParseMediaType(p.Header.Get("Content-Type"))
		if partMT == "multipart/alternative" {
			hasAlternative = true
			_, _ = io.ReadAll(p)
			continue
		}

		disp, dispParams, err := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
		if err != nil {
			t.Fatalf("failed parsing Content-Disposition: %v", err)
		}
		if disp != "attachment" {
			t.Fatalf("expected attachment disposition, got %q", disp)
		}
		filenames = append(filenames, dispParams["filename"])

		bodyBytes, err := io.ReadAll(p)
		if err != nil {
			t.Fatalf("failed reading attachment body: %v", err)
		}
		compact := strings.ReplaceAll(strings.ReplaceAll(string(bodyBytes), "\r", ""), "\n", "")
		decoded, err := base64.StdEncoding.DecodeString(compact)
		if err != nil {
			t.Fatalf("failed decoding base64 attachment: %v", err)
		}
		decodedPayloads = append(decodedPayloads, decoded)
	}

	if !hasAlternative {
		t.Fatal("expected nested multipart/alternative part")
	}
	if len(decodedPayloads) != 2 {
		t.Fatalf("expected 2 attachments, got %d", len(decodedPayloads))
	}
	if !bytes.Equal(decodedPayloads[0], pdfBytes) {
		t.Fatalf("pdf attachment mismatch: got %q", string(decodedPayloads[0]))
	}
	if !bytes.Equal(decodedPayloads[1], txtBytes) {
		t.Fatalf("txt attachment mismatch: got %q", string(decodedPayloads[1]))
	}
	if strings.Contains(filenames[0], "\r") || strings.Contains(filenames[0], "\n") {
		t.Fatalf("filename was not sanitized: %q", filenames[0])
	}
}

func TestLoadEmailAttachments_RejectsPathTraversal(t *testing.T) {
	_, _, err := loadEmailAttachments([]UserAccountEmailAttachment{
		{
			FileName: "passwd",
			FilePath: "/etc/passwd",
		},
	})
	if err == nil {
		t.Fatal("expected error for attachment path outside base dir")
	}
}
