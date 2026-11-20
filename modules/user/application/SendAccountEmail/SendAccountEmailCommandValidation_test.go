package application

import (
	"testing"
)

func TestSendAccountEmailCommandValidation(t *testing.T) {
	// Empty UUIDs
	cmdEmpty := SendAccountEmailCommand{
		UserUUIDs: []string{},
	}
	if err := SendAccountEmailCommandValidation(cmdEmpty); err == nil {
		t.Errorf("Expected error for empty UserUUIDs, got nil")
	}

	// Valid UUIDs
	cmdValid := SendAccountEmailCommand{
		UserUUIDs: []string{"001ff7bd-4c76-4d43-bff3-1b8a7c81477c"},
	}
	if err := SendAccountEmailCommandValidation(cmdValid); err != nil {
		t.Errorf("Expected nil error for valid UserUUIDs, got %v", err)
	}
}
