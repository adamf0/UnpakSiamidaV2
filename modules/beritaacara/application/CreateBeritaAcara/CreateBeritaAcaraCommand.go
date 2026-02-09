package application

type CreateBeritaAcaraCommand struct {
	Tahun            string
	FakultasUnitUuid string
	Tanggal          string
	AuditeeUuid      string
	Auditor1Uuid     *string
	Auditor2Uuid     *string
}
