package domain

import "github.com/google/uuid"

type PreviewTemplate struct {
	UUID            	uuid.UUID
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