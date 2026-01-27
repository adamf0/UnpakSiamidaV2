package commontest

import (
	"UnpakSiamida/common/helper"
	"strings"
	"testing"
	"time"
)

func TestReDoS_Email_LongInput(t *testing.T) {
	evil := strings.Repeat("a", 5_000_000) + "@unpak.ac.id"

	start := time.Now()
	ok := helper.IsValidUnpakEmail(evil)
	elapsed := time.Since(start)

	t.Log("result:", ok)
	t.Log("elapsed:", elapsed)

	if elapsed > 200*time.Millisecond {
		t.Fatalf("Potential DoS: took %v", elapsed)
	}
}

func TestReDoS_Email_AlmostMatch(t *testing.T) {
	evil := strings.Repeat("a.", 2_000_000) + "@unpak.ac.id"

	start := time.Now()
	helper.IsValidUnpakEmail(evil)
	elapsed := time.Since(start)

	t.Log("elapsed:", elapsed)
}

func TestReDoS_UUID_Long(t *testing.T) {
	evil := strings.Repeat("a", 3_000_000)

	start := time.Now()
	err := helper.ValidateUUIDv4(evil)
	elapsed := time.Since(start)

	t.Log("err:", err)
	t.Log("elapsed:", elapsed)

	if elapsed > 200*time.Millisecond {
		t.Fatalf("UUID validation too slow")
	}
}

func BenchmarkIsValidUnpakEmail(b *testing.B) {
	email := "adam.f@unpak.ac.id"
	for i := 0; i < b.N; i++ {
		helper.IsValidUnpakEmail(email)
	}
}

func FuzzIsValidUnpakEmail(f *testing.F) {
	// Seed inputs
	seeds := []string{
		"adam@unpak.ac.id",
		"test..test@unpak.ac.id",
		"",
		"@unpak.ac.id",
		strings.Repeat("a", 300),
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, email string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic for input %q: %v", email, r)
			}
		}()

		helper.IsValidUnpakEmail(email)
	})
}

func FuzzValidateUUIDv4(f *testing.F) {
	seeds := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"",
		"not-a-uuid",
		strings.Repeat("a", 100),
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic for input %q: %v", s, r)
			}
		}()

		_ = helper.ValidateUUIDv4(s)
	})
}
