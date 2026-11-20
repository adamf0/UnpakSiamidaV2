package helper

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDir775AndFile664(t *testing.T) {
	base := t.TempDir()
	nested := filepath.Join(base, "prokerv3", "ILMU KOMPUTER (S1)", "2027", "SK")

	if err := EnsureDir775(nested); err != nil {
		t.Fatalf("EnsureDir775 failed: %v", err)
	}

	for _, rel := range []string{
		"prokerv3",
		filepath.Join("prokerv3", "ILMU KOMPUTER (S1)"),
		filepath.Join("prokerv3", "ILMU KOMPUTER (S1)", "2027"),
		filepath.Join("prokerv3", "ILMU KOMPUTER (S1)", "2027", "SK"),
	} {
		info, err := os.Stat(filepath.Join(base, rel))
		if err != nil {
			t.Fatalf("stat %s failed: %v", rel, err)
		}
		if got := info.Mode().Perm(); got != UploadDirPerm {
			t.Fatalf("dir %s perm = %o, want %o", rel, got, UploadDirPerm)
		}
	}

	f := filepath.Join(nested, "test.csv")
	if err := os.WriteFile(f, []byte("a,b\n"), UploadFilePerm); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := EnsureFile664(f); err != nil {
		t.Fatalf("EnsureFile664 failed: %v", err)
	}
	info, err := os.Stat(f)
	if err != nil {
		t.Fatalf("stat file failed: %v", err)
	}
	if got := info.Mode().Perm(); got != UploadFilePerm {
		t.Fatalf("file perm = %o, want %o", got, UploadFilePerm)
	}
}
