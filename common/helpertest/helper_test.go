package commontest

import (
	"UnpakSiamida/common/helper"
	"strings"
	"testing"
	"time"
)

func TestIsValidUnpakEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid_simple", "adam@unpak.ac.id", true},
		{"valid_dot", "adam.f@unpak.ac.id", true},

		{"invalid_domain", "adam@gmail.com", false},
		{"invalid_plus", "adam+test@unpak.ac.id", false},
		{"double_dot", "adam..f@unpak.ac.id", false},
		{"double_dash", "adam--f@unpak.ac.id", false},
		{"whitespace", "adam f@unpak.ac.id", false},
		{"leading_space", " adam@unpak.ac.id", false},
		{"url_encoded", "adam%40unpak.ac.id", false},
		{"double_url_encoded", "%2540@unpak.ac.id", false},
		{"unicode", "ádam@unpak.ac.id", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.IsValidUnpakEmail(tt.input); got != tt.want {
				t.Fatalf("expected %v got %v", tt.want, got)
			}
		})
	}
}

func TestValidateFakultasUnit(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name  string
		val   interface{}
		level interface{}
		want  bool
	}{
		{"fakultas_required", "", "fakultas", false},
		{"fakultas_valid", validUUID, "fakultas", true},
		{"admin_must_blank", validUUID, "admin", false},
		{"user_blank_ok", "", "user", true},
		{"nil_ok", nil, "user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ValidateFakultasUnit(tt.val, tt.level)
			if (err == nil) != tt.want {
				t.Fatalf("expected %v err=%v", tt.want, err)
			}
		})
	}
}

func TestValidateLevel(t *testing.T) {
	tests := []struct {
		level string
		want  bool
	}{
		{"admin", true},
		{"user", true},
		{"fakultas", true},
		{"superadmin", false},
		{"", false},
	}

	for _, tt := range tests {
		err := helper.ValidateLevel(tt.level)
		if (err == nil) != tt.want {
			t.Fatalf("level=%s expected %v", tt.level, tt.want)
		}
	}
}

func TestValidateParent(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name  string
		input interface{}
		want  bool
	}{
		{"nil_ok", nil, true},
		{"empty_string_ok", "", true},
		{"valid_uuid", validUUID, true},
		{"invalid_uuid", "not-a-uuid", false},
		{"wrong_type", 123, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ValidateParent(tt.input)
			if (err == nil) != tt.want {
				t.Fatalf("expected valid=%v err=%v", tt.want, err)
			}
		})
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int64
		ok    bool
	}{
		{"valid", "123", 123, true},
		{"negative", "-10", -10, true},
		{"zero", "0", 0, true},
		{"syntax_error", "abc", 0, false},
		{"range_error", "999999999999999999999", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := helper.ParseInt64(tt.input)
			if (err == nil) != tt.ok {
				t.Fatalf("expected ok=%v err=%v", tt.ok, err)
			}
			if tt.ok && val != tt.want {
				t.Fatalf("expected %d got %d", tt.want, val)
			}
		})
	}
}

func TestIsValidTugas(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"auditor1", true},
		{"auditor2", true},
		{"auditor3", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := helper.IsValidTugas(tt.input); got != tt.want {
			t.Fatalf("tugas=%s expected %v got %v", tt.input, tt.want, got)
		}
	}
}

func TestIsValidTypeGenerate(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"renstra", true},
		{"dokumen_tambahan", true},
		{"dokumen", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := helper.IsValidTypeGenerate(tt.input); got != tt.want {
			t.Fatalf("type=%s expected %v got %v", tt.input, tt.want, got)
		}
	}
}

func TestFormatWIB(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	tm := time.Date(2024, 1, 2, 3, 4, 5, 0, loc)

	ctx := helper.DateContext{}
	ctx.SetStrategy(helper.IndonesianDateFormatter{})
	got := ctx.FormatWithTime(tm)

	if !strings.Contains(got, "WIB") {
		t.Fatalf("expected WIB timezone, got %s", got)
	}
}

//[pr] ganti ke chain
// func TestFormatDateTimeID(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		input *string
// 		want  string
// 	}{
// 		{"nil", nil, "-"},
// 		{"empty", helper.StrPtr(""), "-"},
// 		{"date_only", helper.StrPtr("2024-01-02"), "02 Januari 2024"},
// 		{"datetime", helper.StrPtr("2024-01-02 10:20:30"), "02 Januari 2024"},
// 		{"rfc3339", helper.StrPtr("2024-01-02T10:20:30Z"), "02 Januari 2024"},
// 		{"invalid", helper.StrPtr("abc"), "abc"},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := helper.FormatDateTimeID(tt.input)
// 			if tt.want != "-" && !strings.Contains(got, tt.want) {
// 				t.Fatalf("expected %s got %s", tt.want, got)
// 			}
// 		})
// 	}
// }

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
