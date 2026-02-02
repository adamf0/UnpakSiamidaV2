package infrastructure

import (
	get2 "UnpakSiamida/modules/laporan/application/GetMonitoringIndikator"
	get "UnpakSiamida/modules/laporan/application/GetMonitoringProker"
	domainlaporan "UnpakSiamida/modules/laporan/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleLaporan(db *gorm.DB) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoLaporan := NewLaporanRepository(db)
	// if err := db.AutoMigrate(&domainlaporan.Laporan{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorLaporan())

	// Register request handler
	mediatr.RegisterRequestHandler[
		get.GetMonitoringProkerByTargetTahunQuery,
		domainlaporan.Paged[domainlaporan.MonitoringProker],
	](&get.GetMonitoringProkerByTargetTahunQueryHandler{
		Repo: repoLaporan,
	})

	mediatr.RegisterRequestHandler[
		get2.GetMonitoringIndikatorByIndikatorTahunQuery,
		domainlaporan.Paged[domainlaporan.MonitoringIndikator],
	](&get2.GetMonitoringIndikatorByIndikatorTahunQueryHandler{
		Repo: repoLaporan,
	})

	return nil
}
