package presentation

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const (
	// MaxTimelineUploadBytes matches main.go BodyLimit (150 MB).
	MaxTimelineUploadBytes   int64 = 150 * 1024 * 1024
	MaxTimelineLinkURLLen          = 500
	MaxTimelineNamaFileRunes       = 255
)

var allowedTimelineUploadExtensions = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".csv":  true,
	".txt":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".zip":  true,
	".rar":  true,
}

func ClampRunes(s string, maxRunes int) string {
	trimmed := strings.TrimSpace(s)
	if maxRunes <= 0 || utf8.RuneCountInString(trimmed) <= maxRunes {
		return trimmed
	}
	runes := []rune(trimmed)
	return strings.TrimSpace(string(runes[:maxRunes]))
}

func IsAllowedTimelineUploadExt(fileName string) bool {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(strings.TrimSpace(fileName))))
	return allowedTimelineUploadExtensions[ext]
}

func ValidateTimelineUploadFile(fileName string, fileSize int64) error {
	cleanName := strings.TrimSpace(fileName)
	if cleanName == "" {
		return errors.New("Nama file bukti tidak boleh kosong")
	}
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(cleanName)))
	if !allowedTimelineUploadExtensions[ext] {
		if ext == "" {
			return errors.New("Format file tidak didukung (tanpa ekstensi). Gunakan PDF, Word, Excel, PowerPoint, gambar (JPG/PNG/WEBP), CSV, TXT, atau ZIP/RAR")
		}
		return fmt.Errorf("Format file %q tidak didukung. Gunakan PDF, Word, Excel, PowerPoint, gambar (JPG/PNG/WEBP), CSV, TXT, atau ZIP/RAR", ext)
	}
	if fileSize < 0 {
		return errors.New("Ukuran file tidak valid")
	}
	if fileSize > MaxTimelineUploadBytes {
		return fmt.Errorf("Ukuran file %q melebihi batas maksimal 150 MB", cleanName)
	}
	return nil
}

func ValidateTimelineLinkURL(rawURL string) (normalizedURL string, fallbackName string, err error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", "", errors.New("URL tautan bukti wajib diisi")
	}
	if strings.ContainsAny(trimmed, " \t\r\n") {
		return "", "", errors.New("URL tautan bukti tidak boleh mengandung spasi")
	}
	if len(trimmed) > MaxTimelineLinkURLLen {
		return "", "", fmt.Errorf("URL tautan bukti maksimal %d karakter", MaxTimelineLinkURLLen)
	}

	parsed, parseErr := url.ParseRequestURI(trimmed)
	if parseErr != nil {
		return "", "", errors.New("Format URL tautan bukti tidak valid")
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", "", errors.New("URL tautan bukti harus diawali http:// atau https://")
	}

	host := parsed.Hostname()
	if host == "" || (!strings.Contains(host, ".") && !strings.EqualFold(host, "localhost")) {
		return "", "", errors.New("Host URL tautan bukti tidak valid")
	}

	parsed.Scheme = scheme
	normalizedURL = parsed.String()
	if len(normalizedURL) > MaxTimelineLinkURLLen {
		return "", "", fmt.Errorf("URL tautan bukti maksimal %d karakter", MaxTimelineLinkURLLen)
	}

	cleanPath := strings.TrimRight(parsed.Path, "/")
	if cleanPath != "" {
		fallbackName = ClampRunes(parsed.Host+cleanPath, MaxTimelineNamaFileRunes)
	} else {
		fallbackName = ClampRunes(parsed.Host, MaxTimelineNamaFileRunes)
	}

	return normalizedURL, fallbackName, nil
}

func IsTimelineLinkRecord(fileObj *domain.ProkerTimelineFileV3) bool {
	if fileObj == nil {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(fileObj.FileType), "link") {
		return true
	}
	fileURL := strings.ToLower(strings.TrimSpace(fileObj.FileURL))
	filePath := strings.ToLower(strings.TrimSpace(fileObj.File))
	return strings.HasPrefix(fileURL, "http://") ||
		strings.HasPrefix(fileURL, "https://") ||
		strings.HasPrefix(filePath, "http://") ||
		strings.HasPrefix(filePath, "https://")
}
