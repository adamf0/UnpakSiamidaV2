package application

type GetMonitoringProkerByTargetTahunQuery struct {
	TahunUuid  string
	TargetUuid string
	Page       *int
	Limit      *int
}
