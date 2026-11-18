package helper

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	UploadDirPerm  os.FileMode = 0775
	UploadFilePerm os.FileMode = 0664
)

var (
	uploadOwnerOnce sync.Once
	uploadOwnerUID  = -1
	uploadOwnerGID  = -1
)

func resolveUploadOwner() (int, int) {
	uploadOwnerOnce.Do(func() {
		ownerUser := strings.TrimSpace(os.Getenv("UPLOAD_OWNER_USER"))
		if ownerUser == "" {
			ownerUser = "www-data"
		}
		ownerGroup := strings.TrimSpace(os.Getenv("UPLOAD_OWNER_GROUP"))
		if ownerGroup == "" {
			ownerGroup = "www-data"
		}

		if u, err := user.Lookup(ownerUser); err == nil {
			if uid, convErr := strconv.Atoi(u.Uid); convErr == nil {
				uploadOwnerUID = uid
			}
			if gid, convErr := strconv.Atoi(u.Gid); convErr == nil {
				uploadOwnerGID = gid
			}
		}
		if g, err := user.LookupGroup(ownerGroup); err == nil {
			if gid, convErr := strconv.Atoi(g.Gid); convErr == nil {
				uploadOwnerGID = gid
			}
		}
	})
	return uploadOwnerUID, uploadOwnerGID
}

func chownUploadPath(path string) {
	uid, gid := resolveUploadOwner()
	if uid < 0 && gid < 0 {
		return
	}
	// Abaikan error jika proses bukan root (mis. dev lokal non-root)
	_ = os.Chown(path, uid, gid)
}

// EnsureDir775 membuat direktori secara rekursif dan memastikan setiap direktori baru
// memiliki permission 0775 dan kepemilikan www-data:www-data (bila proses berhak chown).
func EnsureDir775(dir string) error {
	clean := filepath.Clean(dir)
	var missing []string
	cur := clean
	for {
		if _, err := os.Stat(cur); err == nil {
			break
		}
		missing = append(missing, cur)
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}

	if err := os.MkdirAll(clean, UploadDirPerm); err != nil {
		return err
	}
	for i := len(missing) - 1; i >= 0; i-- {
		_ = os.Chmod(missing[i], UploadDirPerm)
		chownUploadPath(missing[i])
	}
	return nil
}

// EnsureFile664 menetapkan permission 0664 dan kepemilikan www-data:www-data pada file hasil upload.
func EnsureFile664(path string) error {
	chownUploadPath(path)
	return os.Chmod(path, UploadFilePerm)
}
