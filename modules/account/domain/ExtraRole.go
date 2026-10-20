package domain

type ExtraRole struct {
	Tahun            string `gorm:"column:tahun" json:"tahun"`
	Role             string `gorm:"column:role" json:"role"`
	IdFakultasUnit   uint   `gorm:"column:id_fakultas_unit" json:"id_fakultas_unit,omitempty"`
	NamaFakultasUnit string `gorm:"column:nama_fakultas_unit" json:"nama_fakultas_unit,omitempty"`
}