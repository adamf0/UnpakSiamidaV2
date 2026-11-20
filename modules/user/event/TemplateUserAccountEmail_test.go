package event

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRenderUserAccountEmailPlainText(t *testing.T) {
	recipient := UserAccountEmailRecipient{
		UserUUID:  uuid.New(),
		Username:  "evaluator1",
		Password:  "evaluator123",
		Name:      "Dr. John Doe",
		Email:     "john@unpak.ac.id",
		Level:     "user",
		AbacRoles: "evaluator (2027)",
	}

	plain := RenderUserAccountEmailPlainText(recipient)

	if !strings.Contains(plain, "username : evaluator1") {
		t.Errorf("Expected username in plain text, got: %s", plain)
	}
	if !strings.Contains(plain, "password : evaluator123") {
		t.Errorf("Expected password in plain text, got: %s", plain)
	}
	if !strings.Contains(plain, "level : evaluator (2027)") {
		t.Errorf("Expected level in plain text, got: %s", plain)
	}
	if !strings.Contains(plain, "buka siamida.unpak.ac.id dan pilih SIAMIDA V3.") {
		t.Errorf("Expected portal instruction in plain text, got: %s", plain)
	}
	if !strings.Contains(plain, "Terimakasih") {
		t.Errorf("Expected Terimakasih in plain text, got: %s", plain)
	}
}

func TestRenderUserAccountEmailHTML(t *testing.T) {
	recipient := UserAccountEmailRecipient{
		UserUUID:  uuid.New(),
		Username:  "anggaran1",
		Password:  "anggaranPass",
		Name:      "Budi Santoso",
		Email:     "budi@unpak.ac.id",
		Level:     "user",
		AbacRoles: "anggaran (2027)",
	}

	html := RenderUserAccountEmailHTML(recipient, "https://siamida.unpak.ac.id")

	if !strings.Contains(html, "SIAMIDA V3") {
		t.Errorf("Expected SIAMIDA V3 branding in HTML")
	}
	if !strings.Contains(html, "anggaran1") {
		t.Errorf("Expected username in HTML")
	}
	if !strings.Contains(html, "anggaranPass") {
		t.Errorf("Expected password in HTML")
	}
	if !strings.Contains(html, "anggaran (2027)") {
		t.Errorf("Expected abac role in HTML")
	}
	if !strings.Contains(html, "https://siamida.unpak.ac.id") {
		t.Errorf("Expected login URL in HTML")
	}
}
