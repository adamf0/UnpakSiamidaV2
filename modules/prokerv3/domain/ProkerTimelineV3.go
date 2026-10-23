package domain

import (
	common "UnpakSiamida/common/domain"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProkerTimelineV3 struct {
	common.Entity    `json:"-"`
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID             uuid.UUID  `gorm:"type:varchar(36)" json:"UUID"`
	IdAktivitasV2    uint       `gorm:"column:id_aktivitasv2" json:"IdAktivitasV2"`
	TanggalRkAwal    time.Time  `gorm:"column:tanggal_rk_awal;type:date" json:"TanggalRkAwal"`
	TanggalRkAkhir   time.Time  `gorm:"column:tanggal_rk_akhir;type:date" json:"TanggalRkAkhir"`
	TanggalRAwal     *time.Time `gorm:"column:tanggal_r_awal;type:date" json:"TanggalRAwal"`
	TanggalRAkhir    *time.Time `gorm:"column:tanggal_r_akhir;type:date" json:"TanggalRAkhir"`
	NoAktivitas      int        `gorm:"column:no_aktivitas;default:1" json:"NoAktivitas"`
	NamaAktivitas    string     `gorm:"column:nama_aktivitas;size:500" json:"NamaAktivitas"`
	PIC              string     `gorm:"column:pic;size:255" json:"PIC"`
	TitikKontrol     string     `gorm:"column:titik_kontrol;type:text" json:"TitikKontrol"`
	Anggaran         float64    `gorm:"column:anggaran;default:0" json:"Anggaran"`
	JadwalRk         string     `gorm:"column:jadwal_rk;type:text" json:"JadwalRk"`
	JadwalR          string     `gorm:"column:jadwal_r;type:text" json:"JadwalR"`
	Capaian          string     `gorm:"column:capaian;type:text" json:"Capaian"`
	Masalah          string     `gorm:"column:masalah;type:text" json:"Masalah"`
	PemecahanMasalah string     `gorm:"column:pemecahan_masalah;type:text" json:"PemecahanMasalah"`
	Status           string     `gorm:"column:status;size:50;default:'Belum Terlaksana'" json:"Status"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"CreatedAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"UpdatedAt"`
}

func (ProkerTimelineV3) TableName() string {
	return "aktivitasv2_timeline"
}

// === CREATE ===
func NewProkerTimelineV3(
	act *AktivitasV2,
	uid uuid.UUID,
	noAktivitas int,
	namaAktivitas string,
	pic string,
	titikKontrol string,
	anggaran float64,
	tanggalRkAwal time.Time,
	tanggalRkAkhir time.Time,
	tanggalRAwal *time.Time,
	tanggalRAkhir *time.Time,
	jadwalRk string,
	jadwalR string,
	capaian string,
	masalah string,
	pemecahanMasalah string,
	status string,
) common.ResultValue[*ProkerTimelineV3] {
	if act == nil {
		return common.FailureValue[*ProkerTimelineV3](AktivitasNotFound())
	}

	if anggaran < 0 {
		return common.FailureValue[*ProkerTimelineV3](NegativeAnggaranTimeline())
	}

	maxAnggaran := act.Anggaran
	if act.Realisasi != nil && *act.Realisasi > 0 {
		maxAnggaran = *act.Realisasi
	}
	if maxAnggaran > 0 && anggaran > maxAnggaran {
		return common.FailureValue[*ProkerTimelineV3](AnggaranTimelineExceedsRealisasi(anggaran, maxAnggaran))
	}

	if uid == uuid.Nil {
		uid = uuid.New()
	}
	if strings.TrimSpace(titikKontrol) == "" {
		titikKontrol = "RK"
	}
	if strings.TrimSpace(status) == "" {
		status = "Belum Terlaksana"
	}
	if strings.TrimSpace(namaAktivitas) == "" {
		namaAktivitas = act.Aktivitas
	}
	if strings.TrimSpace(pic) == "" {
		pic = act.PIC
	}

	tl := &ProkerTimelineV3{
		UUID:             uid,
		IdAktivitasV2:    act.ID,
		NoAktivitas:      noAktivitas,
		NamaAktivitas:    strings.TrimSpace(namaAktivitas),
		PIC:              strings.TrimSpace(pic),
		TitikKontrol:     strings.TrimSpace(titikKontrol),
		Anggaran:         anggaran,
		TanggalRkAwal:    tanggalRkAwal,
		TanggalRkAkhir:   tanggalRkAkhir,
		TanggalRAwal:     tanggalRAwal,
		TanggalRAkhir:    tanggalRAkhir,
		JadwalRk:         strings.TrimSpace(jadwalRk),
		JadwalR:          strings.TrimSpace(jadwalR),
		Capaian:          strings.TrimSpace(capaian),
		Masalah:          strings.TrimSpace(masalah),
		PemecahanMasalah: strings.TrimSpace(pemecahanMasalah),
		Status:           strings.TrimSpace(status),
	}

	return common.SuccessValue(tl)
}

// === UPDATE ===
func UpdateProkerTimelineV3(
	prev *ProkerTimelineV3,
	act *AktivitasV2,
	anggaran *float64,
	namaAktivitas string,
	pic string,
	titikKontrol string,
	noAktivitas int,
	tanggalRkAwal *time.Time,
	tanggalRkAkhir *time.Time,
	tanggalRAwal *time.Time,
	tanggalRAkhir *time.Time,
	jadwalRk string,
	jadwalR string,
	capaian string,
	masalah string,
	pemecahanMasalah string,
	status string,
) common.ResultValue[*ProkerTimelineV3] {
	if prev == nil {
		return common.FailureValue[*ProkerTimelineV3](TimelineNotFound())
	}

	if anggaran != nil {
		if *anggaran < 0 {
			return common.FailureValue[*ProkerTimelineV3](NegativeAnggaranTimeline())
		}
		if act != nil {
			maxAnggaran := act.Anggaran
			if act.Realisasi != nil && *act.Realisasi > 0 {
				maxAnggaran = *act.Realisasi
			}
			if maxAnggaran > 0 && *anggaran > maxAnggaran {
				return common.FailureValue[*ProkerTimelineV3](AnggaranTimelineExceedsRealisasi(*anggaran, maxAnggaran))
			}
		}
		prev.Anggaran = *anggaran
	}

	if strings.TrimSpace(namaAktivitas) != "" {
		prev.NamaAktivitas = strings.TrimSpace(namaAktivitas)
	}
	if strings.TrimSpace(pic) != "" {
		prev.PIC = strings.TrimSpace(pic)
	}
	if strings.TrimSpace(titikKontrol) != "" {
		prev.TitikKontrol = strings.TrimSpace(titikKontrol)
	}
	if noAktivitas > 0 {
		prev.NoAktivitas = noAktivitas
	}
	if tanggalRkAwal != nil {
		prev.TanggalRkAwal = *tanggalRkAwal
	}
	if tanggalRkAkhir != nil {
		prev.TanggalRkAkhir = *tanggalRkAkhir
	}
	if tanggalRAwal != nil {
		prev.TanggalRAwal = tanggalRAwal
	}
	if tanggalRAkhir != nil {
		prev.TanggalRAkhir = tanggalRAkhir
	}
	if strings.TrimSpace(jadwalRk) != "" {
		prev.JadwalRk = strings.TrimSpace(jadwalRk)
	}
	if strings.TrimSpace(jadwalR) != "" {
		prev.JadwalR = strings.TrimSpace(jadwalR)
	}
	if strings.TrimSpace(capaian) != "" {
		prev.Capaian = strings.TrimSpace(capaian)
	}
	if strings.TrimSpace(masalah) != "" {
		prev.Masalah = strings.TrimSpace(masalah)
	}
	if strings.TrimSpace(pemecahanMasalah) != "" {
		prev.PemecahanMasalah = strings.TrimSpace(pemecahanMasalah)
	}
	if strings.TrimSpace(status) != "" {
		prev.Status = strings.TrimSpace(status)
	}
	prev.UpdatedAt = time.Now()

	return common.SuccessValue(prev)
}
