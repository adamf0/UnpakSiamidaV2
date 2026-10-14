package application

type UpdateSwotCommand struct {
	UUID           string   `json:"uuid"`
	FakultasUuid   string   `json:"fakultas_unit"`
	IdFakultasUnit *int     `json:"id_fakultas_unit"`
	TahunUuid      string   `json:"tahun_uuid"`
	IdMasterTahun  *int     `json:"id_master_tahun"`
	Tahun          string   `json:"tahun"`
	Kategori       string   `json:"kategori"`
	Faktor         string   `json:"faktor"`
	Deskripsi      string   `json:"deskripsi"`
	Strategi       string   `json:"strategi"`
	Bobot          *float64 `json:"bobot"`
	Rating         *int     `json:"rating"`
	Skor           *float64 `json:"skor"`
}
