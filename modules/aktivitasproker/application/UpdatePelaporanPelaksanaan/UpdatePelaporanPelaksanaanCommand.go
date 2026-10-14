package application

type UpdatePelaporanPelaksanaanCommand struct {
	Uuid            string
	StatusRealisasi string
	TanggalRAwal    string
	TanggalRAkhir   string
	Detail          string
}
