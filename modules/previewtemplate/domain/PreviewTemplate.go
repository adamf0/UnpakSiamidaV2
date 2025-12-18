package domain

type PreviewTemplate struct {
	Tahun               int
	IndikatorId         int
	Indikator           string
	IndikatorTahun      int
	IsPertanyaan        int

	ParentIndikatorId   *int

	FakultasUnitId      uint
	FakultasUnit        string
	FakultasUnitType    string
	Fakultas            string
	Klasifikasi         string
	Satuan              *string
	Target              string
	Kategori            string

	Pointing            string
}