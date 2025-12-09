package domain

import (
	"time"

	common "UnpakSiamida/common/domain"

	"github.com/google/uuid"
)

type TemplateDokumenTambahan struct {
	common.Entity
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID      `gorm:"type:char(36);uniqueIndex"`
	Tahun       string
	JenisFileID uint           `gorm:"column:jenis_file"`
	Pertanyaan  string         `gorm:"type:longtext"`
	Klasifikasi string
	Kategori    string         `gorm:"column:fakultas_prodi_unit"`
	Tugas       string
}
func (TemplateDokumenTambahan) TableName() string {
	return "template_dokumen_tambahan"
}


// === CREATE ===
func NewTemplateDokumenTambahan(
	tahun              	string,
	jenisFileID 		uint,
	pertanyaan      	string,
	klasifikasi        	string,
	kategori           	string,
	tugas              	string,
) common.ResultValue[*TemplateDokumenTambahan] {
	if jenisFileID<=0{
		return common.FailureValue[*TemplateDokumenTambahan](JenisFileNotFound())
	}

	templatedokumentambahan := &TemplateDokumenTambahan{
		UUID: uuid.New(),
		Tahun: tahun,
		JenisFileID: jenisFileID,
		Pertanyaan: pertanyaan,
		Klasifikasi: klasifikasi,
		Kategori: kategori,
		Tugas: tugas,
	}

	templatedokumentambahan.Raise(TemplateDokumenTambahanCreatedEvent{
		EventID:    uuid.New(),
		OccurredOn: time.Now().UTC(),
		TemplateDokumenTambahanUUID:   templatedokumentambahan.UUID,
	})

	return common.SuccessValue(templatedokumentambahan)
}

// === UPDATE ===
func UpdateTemplateDokumenTambahan(
	prev *TemplateDokumenTambahan, 
	uid uuid.UUID,
	tahun              	string,
	jenisFileID 		uint,
	pertanyaan      	string,
	klasifikasi        	string,
	kategori           	string,
	tugas              	string,
) common.ResultValue[*TemplateDokumenTambahan] {
	
	if prev == nil {
		return common.FailureValue[*TemplateDokumenTambahan](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*TemplateDokumenTambahan](InvalidData())
	}
	if jenisFileID<=0{
		return common.FailureValue[*TemplateDokumenTambahan](JenisFileNotFound())
	}

	prev.Tahun = tahun
	prev.JenisFileID = jenisFileID
	prev.Pertanyaan = pertanyaan
	prev.Klasifikasi = klasifikasi
	prev.Kategori = kategori
	prev.Tugas = tugas

	prev.Raise(TemplateDokumenTambahanUpdatedEvent{
		EventID:   	uuid.New(),
		OccurredOn: time.Now().UTC(),
		TemplateDokumenTambahanUUID:   prev.UUID,
	})

	return common.SuccessValue(prev)
}
