package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MatrixProkerDokumen struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"Id"`
	UUID           uuid.UUID      `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	IdMatrikProker uint           `gorm:"column:id_matrik_proker;index" json:"IdMatrikProker"`
	JenisDokumen   string         `gorm:"column:jenis_dokumen;type:varchar(100)" json:"JenisDokumen"`
	NamaFile       string         `gorm:"column:nama_file;type:varchar(255)" json:"NamaFile"`
	File           string         `gorm:"column:file;type:varchar(500)" json:"File"`
	FileURL        string         `gorm:"column:file_url;type:varchar(500)" json:"FileURL"`
	FileSize         int64          `gorm:"column:file_size" json:"FileSize"`
	FileType         string         `gorm:"column:file_type;type:varchar(100)" json:"FileType"`
	StatusVerifikasi string         `gorm:"column:status_verifikasi;type:varchar(100);default:'belum_terverifikasi'" json:"StatusVerifikasi"`
	Catatan          string         `gorm:"column:catatan;type:text" json:"Catatan"`
	CatatanPenolakan string         `gorm:"column:catatan_penolakan;type:text" json:"CatatanPenolakan"`
	Keterangan       string         `gorm:"column:keterangan;type:text" json:"Keterangan"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"CreatedAt"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"UpdatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (MatrixProkerDokumen) TableName() string {
	return "matrix_proker_dokumen"
}
