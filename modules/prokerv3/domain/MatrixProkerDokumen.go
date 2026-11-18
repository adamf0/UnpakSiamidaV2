package domain

import (
	"strings"
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
	StatusVerifikasi         string         `gorm:"column:status_verifikasi;type:varchar(100);default:'belum_terverifikasi'" json:"StatusVerifikasi"`
	StatusVerifikasiFakultas string         `gorm:"column:status_verifikasi_fakultas;type:varchar(50);default:'belum_terverifikasi'" json:"StatusVerifikasiFakultas"`
	CatatanTolakFakultas     string         `gorm:"column:catatan_tolak_fakultas;type:text" json:"CatatanTolakFakultas"`
	StatusVerifikasiAdmin    string         `gorm:"column:status_verifikasi_admin;type:varchar(50);default:'belum_terverifikasi'" json:"StatusVerifikasiAdmin"`
	CatatanTolakAdmin        string         `gorm:"column:catatan_tolak_admin;type:text" json:"CatatanTolakAdmin"`
	Catatan                  string         `gorm:"column:catatan;type:text" json:"Catatan"`
	CatatanPenolakan         string         `gorm:"column:catatan_penolakan;type:text" json:"CatatanPenolakan"`
	Keterangan       string         `gorm:"column:keterangan;type:text" json:"Keterangan"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"CreatedAt"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"UpdatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (MatrixProkerDokumen) TableName() string {
	return "matrix_proker_dokumen"
}

func (m MatrixProkerDokumen) EffectiveFacultyStatus() string {
	s := strings.ToLower(strings.TrimSpace(m.StatusVerifikasiFakultas))
	if s != "" && s != "belum_terverifikasi" {
		return s
	}
	// Fallback untuk data lama sebelum kolom status_verifikasi_fakultas diisi
	overall := strings.ToLower(strings.TrimSpace(m.StatusVerifikasi))
	if overall == "terverifikasi_fakultas" || overall == "terverifikasi" {
		return "terverifikasi"
	}
	return "belum_terverifikasi"
}

func (m MatrixProkerDokumen) EffectiveAdminStatus() string {
	s := strings.ToLower(strings.TrimSpace(m.StatusVerifikasiAdmin))
	if s != "" && s != "belum_terverifikasi" {
		return s
	}
	overall := strings.ToLower(strings.TrimSpace(m.StatusVerifikasi))
	if overall == "terverifikasi" {
		return "terverifikasi"
	}
	return "belum_terverifikasi"
}

func (m MatrixProkerDokumen) IsVerifiedByFakultas() bool {
	return m.EffectiveFacultyStatus() == "terverifikasi"
}

func (m MatrixProkerDokumen) IsVerifiedByAdmin() bool {
	return m.EffectiveAdminStatus() == "terverifikasi"
}

type DokumenVerifikasiUpdate struct {
	StatusVerifikasi         string
	StatusVerifikasiFakultas *string
	CatatanTolakFakultas     *string
	StatusVerifikasiAdmin    *string
	CatatanTolakAdmin        *string
	CatatanPenolakan         string
}

type DokumenUnitContext struct {
	Dokumen          MatrixProkerDokumen
	IdFakultasUnit   *int   `gorm:"column:id_fakultas_unit"`
	TypeFakultasUnit string `gorm:"column:type_fakultas_unit"`
	ParentFakultas   string `gorm:"column:parent_fakultas"`
	NamaFakultasUnit string `gorm:"column:nama_fak_prod_unit"`
	UuidFakultasUnit string `gorm:"column:uuid_fakultas_unit"`
}

func (d DokumenUnitContext) IsProdi() bool {
	return strings.EqualFold(strings.TrimSpace(d.TypeFakultasUnit), "prodi") ||
		strings.Contains(strings.ToLower(d.NamaFakultasUnit), "prodi")
}
