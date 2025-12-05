package helper

import (
	"fmt"
	"regexp"
	"strings"
)

func IsValidUnpakEmail(email string) bool {

	// 1. Base pattern
	reg := regexp.MustCompile(
		`^[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?@unpak\.ac\.id$`,
	)

	if !reg.MatchString(email) {
		return false
	}

	// 2. No plus (+)
	if regexp.MustCompile(`\+`).MatchString(email) {
		return false
	}

	// 3. Double separator
	if regexp.MustCompile(`(\.\.|__|--)`).MatchString(email) {
		return false
	}

	// 4. No whitespace
	if regexp.MustCompile(`\s`).MatchString(email) {
		return false
	}

	// 5. No URL-encoded chars
	if regexp.MustCompile(`%[0-9A-Fa-f]{2}`).MatchString(email) {
		return false
	}
	if regexp.MustCompile(`%25[0-9A-Fa-f]{2}`).MatchString(email) {
		return false
	}

	// 6. No non-ASCII
	if regexp.MustCompile(`[^\x20-\x7F]`).MatchString(email) {
		return false
	}

	return true
}

func ValidateUnpakEmail(value interface{}) error {
	if value == nil {
		return fmt.Errorf("Email cannot be blank")
	}

	email, ok := value.(string)
	if !ok {
		return fmt.Errorf("Email invalid type")
	}

	if !IsValidUnpakEmail(email) {
		return fmt.Errorf("Email is not valid unpak.ac.id")
	}

	return nil
}


func ValidateUUIDv4(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("UUID invalid type")
	}

	s = strings.TrimSpace(s)

	// Cek null padding ASCII ( \x00 )
	if strings.Contains(s, "\x00") {
		return fmt.Errorf("UUID contains invalid null padding")
	}

	// format regex UUID v4
	regexV4 := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`

	matched := regexp.MustCompile(regexV4).MatchString(s)
	if !matched {
		return fmt.Errorf("UUID must be a valid UUIDv4 format")
	}

	return nil
}