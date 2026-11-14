package presentation

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"strings"
	"testing"
)

func TestValidateTimelineUploadFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		fileSize int64
		wantErr  bool
	}{
		{name: "valid pdf", fileName: "laporan.pdf", fileSize: 1024, wantErr: false},
		{name: "valid uppercase extension", fileName: "FOTO_KEGIATAN.JPG", fileSize: 2048, wantErr: false},
		{name: "valid exact max 150MB", fileName: "archive.zip", fileSize: MaxTimelineUploadBytes, wantErr: false},
		{name: "valid doc/xls/ppt", fileName: "data.xlsx", fileSize: 512, wantErr: false},
		{name: "reject over 150MB", fileName: "huge.pdf", fileSize: MaxTimelineUploadBytes + 1, wantErr: true},
		{name: "reject executable", fileName: "script.exe", fileSize: 100, wantErr: true},
		{name: "reject html upload", fileName: "xss.html", fileSize: 100, wantErr: true},
		{name: "reject missing extension", fileName: "README", fileSize: 100, wantErr: true},
		{name: "reject negative size", fileName: "valid.pdf", fileSize: -1, wantErr: true},
		{name: "reject empty filename", fileName: "   ", fileSize: 100, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateTimelineUploadFile(tc.fileName, tc.fileSize)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateTimelineUploadFile(%q, %d) error = %v, wantErr %v", tc.fileName, tc.fileSize, err, tc.wantErr)
			}
		})
	}
}

func TestValidateTimelineLinkURL(t *testing.T) {
	tooLongURL := "https://drive.google.com/file/d/" + strings.Repeat("a", 500)

	tests := []struct {
		name        string
		rawURL      string
		wantName    string
		wantScheme  string
		wantErr     bool
	}{
		{
			name:       "valid https google drive link",
			rawURL:     "  https://drive.google.com/drive/folders/12345  ",
			wantName:   "drive.google.com/drive/folders/12345",
			wantScheme: "https://",
			wantErr:    false,
		},
		{
			name:       "valid http link with localhost",
			rawURL:     "http://localhost:8080/bukti",
			wantName:   "localhost:8080/bukti",
			wantScheme: "http://",
			wantErr:    false,
		},
		{
			name:       "valid uppercase scheme normalized",
			rawURL:     "HTTPS://unpak.ac.id/",
			wantName:   "unpak.ac.id",
			wantScheme: "https://",
			wantErr:    false,
		},
		{name: "reject empty", rawURL: "   ", wantErr: true},
		{name: "reject whitespace in url", rawURL: "https://unpak.ac.id/my file.pdf", wantErr: true},
		{name: "reject javascript scheme", rawURL: "javascript:alert(1)", wantErr: true},
		{name: "reject ftp scheme", rawURL: "ftp://files.unpak.ac.id/doc.pdf", wantErr: true},
		{name: "reject missing scheme", rawURL: "drive.google.com/file/123", wantErr: true},
		{name: "reject single-label host", rawURL: "https://dokumen", wantErr: true},
		{name: "reject over 500 chars", rawURL: tooLongURL, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			normalized, fallbackName, err := ValidateTimelineLinkURL(tc.rawURL)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateTimelineLinkURL(%q) error = %v, wantErr %v", tc.rawURL, err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if !strings.HasPrefix(normalized, tc.wantScheme) {
				t.Fatalf("normalized URL %q does not start with %q", normalized, tc.wantScheme)
			}
			if fallbackName != tc.wantName {
				t.Fatalf("fallbackName = %q, want %q", fallbackName, tc.wantName)
			}
		})
	}
}

func TestClampRunesAndLinkDetection(t *testing.T) {
	over := strings.Repeat("é", 260)
	clamped := ClampRunes(over, MaxTimelineNamaFileRunes)
	if len([]rune(clamped)) != MaxTimelineNamaFileRunes {
		t.Fatalf("ClampRunes len = %d, want %d", len([]rune(clamped)), MaxTimelineNamaFileRunes)
	}

	if !IsTimelineLinkRecord(&domain.ProkerTimelineFileV3{FileType: "link"}) {
		t.Fatal("expected FileType link to be recognized as link record")
	}
	if !IsTimelineLinkRecord(&domain.ProkerTimelineFileV3{FileURL: "https://unpak.ac.id"}) {
		t.Fatal("expected https FileURL to be recognized as link record")
	}
	if IsTimelineLinkRecord(&domain.ProkerTimelineFileV3{File: "./upload/prokerv2/a.pdf", FileURL: "/upload/prokerv2/a.pdf", FileType: "application/pdf"}) {
		t.Fatal("did not expect local PDF file to be recognized as link record")
	}
}
