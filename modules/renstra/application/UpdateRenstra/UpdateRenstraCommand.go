package application

type UpdateRenstraCommand struct {
    Uuid     	 string
	Tahun string
	FakultasUnit string
	PeriodeUploadMulai string
	PeriodeUploadAkhir string
	PeriodeAssesmentDokumenMulai string
	PeriodeAssesmentDokumenAkhir string
	PeriodeAssesmentLapanganMulai string
	PeriodeAssesmentLapanganAkhir string
	Auditee string
	Auditor1 string
	Auditor2 string
}
