package application

type GetMonitoringIndikatorByIndikatorTahunQuery struct {
	Tahun         string
	IndikatorUuid string
	Page          *int
	Limit         *int
}
