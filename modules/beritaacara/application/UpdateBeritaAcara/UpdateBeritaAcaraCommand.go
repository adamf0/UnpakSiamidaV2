package application

type UpdateBeritaAcaraCommand struct {
	Uuid         string
	Tahun        string
	FakultasUnit int
	Tanggal      string
	Auditee      *int
	Auditor1     *int
	Auditor2     *int
}
