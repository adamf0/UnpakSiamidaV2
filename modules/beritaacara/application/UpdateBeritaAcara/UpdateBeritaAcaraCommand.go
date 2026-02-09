package application

type UpdateBeritaAcaraCommand struct {
	Uuid             string
	Tahun            string
	FakultasUnitUuid string
	Tanggal          string
	AuditeeUuid      string
	Auditor1Uuid     *string
	Auditor2Uuid     *string
}
