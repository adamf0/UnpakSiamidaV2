package application

type CreateBeritaAcaraCommand struct {
	Tahun        string
	FakultasUnit int
	Tanggal      string
	Auditee      *int
	Auditor1     *int
	Auditor2     *int
}
