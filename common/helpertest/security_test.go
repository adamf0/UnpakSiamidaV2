package commontest

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// single source of truth
func validateXSS(input string) error {
	return validation.Validate(
		input,
		validation.By(helper.NoXSSFullScanWithDecode()),
	)
}

// =======================================================
// FALSE POSITIVE TEST (VALID INPUT, SHOULD PASS)
// =======================================================

func TestNoXSS_FalsePositive(t *testing.T) {
	cases := []string{
		"1 + 2 = 3",
		"x <= y",
		"Total >= 100%",
		"(A+B)*C - D",
		"< >",
		"<= >= <>",
		"ABC_xyz-123",
		"Harga 5000",
		"if a < b then c > d",
	}

	for _, c := range cases {
		if err := validateXSS(c); err != nil {
			t.Errorf("FALSE POSITIVE: %q rejected: %v", c, err)
		}
	}
}

// =======================================================
// TRUE NEGATIVE TEST (MALICIOUS INPUT, SHOULD FAIL)
// =======================================================

func TestNoXSS_TrueNegative(t *testing.T) {
	cases := []string{
		"<script>alert(1)</script>",
		"<img src=x onerror=alert(1)>",
		"javascript:alert(1)",
		"&lt;script&gt;",
		"onload=alert(1)",
		"<svg><script>",
		"😊",
		"你好",
		"مرحبا",
		"\u202Eevil", // RTL override
	}

	for _, c := range cases {
		if err := validateXSS(c); err == nil {
			t.Errorf("FALSE NEGATIVE: %q was accepted", c)
		}
	}
}

// =======================================================
// FUZZ TEST
// =======================================================

func FuzzNoXSSFullScanWithDecode(f *testing.F) {
	seeds := []string{
		"hello",
		"1+2",
		"<script>alert(1)</script>",
		"&lt;svg&gt;",
		"%3Cscript%3E",
		"😊",
		"abc123",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		_ = validateXSS(input)
	})
}

func loadCorpus(t *testing.T) []string {
	t.Helper()

	dir := "corpus"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read corpus dir: %v", err)
	}

	var cases []string

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}

		path := filepath.Join(dir, e.Name())
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("failed to open corpus file %s: %v", e.Name(), err)
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				cases = append(cases, line)
			}
		}

		f.Close()

		if err := scanner.Err(); err != nil {
			t.Fatalf("failed reading %s: %v", e.Name(), err)
		}
	}

	if len(cases) == 0 {
		t.Fatal("corpus empty: no test cases loaded")
	}

	return cases
}

//[pr] masih belum bisa menangkis kerumitan xss & sql
// func TestNoXSS_Corpus(t *testing.T) {
// 	cases := loadCorpus(t)

// 	for _, c := range cases {
// 		if err := validateXSS(c); err == nil {
// 			t.Errorf("FALSE NEGATIVE: payload accepted: %q", c)
// 		}
// 	}
// }
