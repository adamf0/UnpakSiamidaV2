package application

type CreatePelaporanPelaksanaanCommand struct {
	FakultasUuid    string
	MataProgramUuid string
	Aktivitas       string
	PIC             string
	TanggalRKAwal   string
	TanggalRKAkhir  string
	StatusRealisasi string
	TanggalRAwal    string
	TanggalRAkhir   string
	Detail          string
}
