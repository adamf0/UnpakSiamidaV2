package domain

import (
	"time"

	"github.com/google/uuid"
)

type SwotDefault struct {
	Id               uint      `json:"Id"`
	UUID             uuid.UUID `json:"UUID"`
	IdFakultasUnit   *int      `json:"IdFakultasUnit"`
	NamaFakultasUnit *string   `json:"NamaFakultasUnit"`
	UuidFakultasUnit *string   `json:"UuidFakultasUnit"`
	IdMasterTahun    *int      `json:"IdMasterTahun"`
	Tahun            string    `json:"Tahun"`
	TahunUuid        *string   `json:"TahunUuid"`
	Kategori         string    `json:"Kategori"`
	Faktor           string    `json:"Faktor"`
	Deskripsi        string    `json:"Deskripsi"`
	Strategi         string    `json:"Strategi"`
	Bobot            *float64  `json:"Bobot"`
	Rating           *int      `json:"Rating"`
	Skor             *float64  `json:"Skor"`
	CreatedAt        time.Time `json:"CreatedAt"`
	UpdatedAt        time.Time `json:"UpdatedAt"`
}
