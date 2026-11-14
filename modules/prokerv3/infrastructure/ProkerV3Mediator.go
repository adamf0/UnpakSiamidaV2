package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	"context"

	checkJadwal "UnpakSiamida/modules/prokerv3/application/CheckJadwalProkerV3"
	createDokumen "UnpakSiamida/modules/prokerv3/application/CreateProkerDokumen"
	createJadwal "UnpakSiamida/modules/prokerv3/application/CreateJadwalProkerV3"
	create "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"
	deleteDokumen "UnpakSiamida/modules/prokerv3/application/DeleteDokumenV3"
	deleteJadwal "UnpakSiamida/modules/prokerv3/application/DeleteJadwalProkerV3"
	deletePagu "UnpakSiamida/modules/prokerv3/application/DeletePaguAnggaran"
	delete "UnpakSiamida/modules/prokerv3/application/DeleteProkerV3"
	getAllJadwal "UnpakSiamida/modules/prokerv3/application/GetAllJadwalProkerV3"
	getAllPagu "UnpakSiamida/modules/prokerv3/application/GetAllPaguAnggaran"
	getAll "UnpakSiamida/modules/prokerv3/application/GetAllProkerV3"
	getCoa "UnpakSiamida/modules/prokerv3/application/GetCoa"
	getDashboard "UnpakSiamida/modules/prokerv3/application/GetDashboardProkerMonitoring"
	getDokumen "UnpakSiamida/modules/prokerv3/application/GetDokumenV3"
	getHistory "UnpakSiamida/modules/prokerv3/application/GetHistoryStateV3"
	getJadwal "UnpakSiamida/modules/prokerv3/application/GetJadwalProkerV3"
	getProkerBatalSetting "UnpakSiamida/modules/prokerv3/application/GetProkerBatalSetting"
	get "UnpakSiamida/modules/prokerv3/application/GetProkerV3"
	getOptions "UnpakSiamida/modules/prokerv3/application/GetProkerV3Options"
	getSankey "UnpakSiamida/modules/prokerv3/application/GetSankeyData"
	getTargetUnpak "UnpakSiamida/modules/prokerv3/application/GetTargetUnpak"
	getUnitPagu "UnpakSiamida/modules/prokerv3/application/GetUnitPaguSummary"
	savePagu "UnpakSiamida/modules/prokerv3/application/SavePaguAnggaran"
	updateBatal "UnpakSiamida/modules/prokerv3/application/UpdateBatalProker"
	updateBeban "UnpakSiamida/modules/prokerv3/application/UpdateBebanYayasan"
	updateJadwal "UnpakSiamida/modules/prokerv3/application/UpdateJadwalProkerV3"
	updateProkerBatalSetting "UnpakSiamida/modules/prokerv3/application/UpdateProkerBatalSetting"
	update "UnpakSiamida/modules/prokerv3/application/UpdateProkerV3"
	updateStatus "UnpakSiamida/modules/prokerv3/application/UpdateStatusV3"
	updateVerifikasiDokumen "UnpakSiamida/modules/prokerv3/application/UpdateVerifikasiDokumenV3"
	domain "UnpakSiamida/modules/prokerv3/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleProkerV3(db *gorm.DB, dbs ...*gorm.DB) error {
	repoProker := NewProkerV3Repository(db, dbs...)
	_ = repoProker.EnsureAppSettingsTable(context.Background())
	_ = db.AutoMigrate(&domain.JadwalProkerV3{})

	// 1. Proker Core CRUD
	mediatr.RegisterRequestHandler[
		create.CreateProkerV3Command,
		string,
	](&create.CreateProkerV3CommandHandler{
		Repo: repoProker,
		DB:   db,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateProkerV3Command,
		string,
	](&update.UpdateProkerV3CommandHandler{
		Repo: repoProker,
		DB:   db,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteProkerV3Command,
		string,
	](&delete.DeleteProkerV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		get.GetProkerV3ByUuidQuery,
		*domain.ProkerV3Default,
	](&get.GetProkerV3ByUuidQueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllProkerV3Query,
		commondomain.Paged[domain.ProkerV3Default],
	](&getAll.GetAllProkerV3QueryHandler{
		Repo: repoProker,
	})

	// 2. Options, Target, COA
	mediatr.RegisterRequestHandler[
		getOptions.GetProkerV3OptionsQuery,
		*domain.ProkerOptionsResult,
	](&getOptions.GetProkerV3OptionsQueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getTargetUnpak.GetTargetUnpakQuery,
		*domain.TargetUnpakResult,
	](&getTargetUnpak.GetTargetUnpakQueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getCoa.GetCoaQuery,
		[]domain.CoaItem,
	](&getCoa.GetCoaQueryHandler{
		Repo: repoProker,
	})

	// 3. Status, Verification, History
	mediatr.RegisterRequestHandler[
		updateStatus.UpdateStatusV3Command,
		*domain.UpdateStatusResult,
	](&updateStatus.UpdateStatusV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getHistory.GetHistoryStateV3Query,
		[]domain.HistoryState,
	](&getHistory.GetHistoryStateV3QueryHandler{
		Repo: repoProker,
	})

	// 4. Dokumen
	mediatr.RegisterRequestHandler[
		getDokumen.GetDokumenV3Query,
		*domain.ProkerDokumenResult,
	](&getDokumen.GetDokumenV3QueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		deleteDokumen.DeleteDokumenV3Command,
		string,
	](&deleteDokumen.DeleteDokumenV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		updateVerifikasiDokumen.UpdateVerifikasiDokumenV3Command,
		*domain.MatrixProkerDokumen,
	](&updateVerifikasiDokumen.UpdateVerifikasiDokumenV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		createDokumen.CreateProkerDokumenCommand,
		*domain.MatrixProkerDokumen,
	](&createDokumen.CreateProkerDokumenCommandHandler{
		Repo: repoProker,
	})

	// 5. Beban Yayasan & Batal Proker
	mediatr.RegisterRequestHandler[
		updateBeban.UpdateBebanYayasanCommand,
		*updateBeban.UpdateBebanYayasanResult,
	](&updateBeban.UpdateBebanYayasanCommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		updateBatal.UpdateBatalProkerCommand,
		*updateBatal.UpdateBatalProkerResult,
	](&updateBatal.UpdateBatalProkerCommandHandler{
		Repo: repoProker,
	})

	// 6. Settings Batal
	mediatr.RegisterRequestHandler[
		getProkerBatalSetting.GetProkerBatalSettingQuery,
		*getProkerBatalSetting.GetProkerBatalSettingResult,
	](&getProkerBatalSetting.GetProkerBatalSettingQueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		updateProkerBatalSetting.UpdateProkerBatalSettingCommand,
		*updateProkerBatalSetting.UpdateProkerBatalSettingResult,
	](&updateProkerBatalSetting.UpdateProkerBatalSettingCommandHandler{
		Repo: repoProker,
	})

	// 7. Jadwal Proker V3
	mediatr.RegisterRequestHandler[
		getAllJadwal.GetAllJadwalProkerV3Query,
		*getAllJadwal.GetAllJadwalProkerV3Result,
	](&getAllJadwal.GetAllJadwalProkerV3QueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getJadwal.GetJadwalProkerV3Query,
		*domain.JadwalProkerV3DTO,
	](&getJadwal.GetJadwalProkerV3QueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		checkJadwal.CheckJadwalProkerV3Query,
		*domain.CheckJadwalResult,
	](&checkJadwal.CheckJadwalProkerV3QueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		createJadwal.CreateJadwalProkerV3Command,
		*createJadwal.CreateJadwalProkerV3Result,
	](&createJadwal.CreateJadwalProkerV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		updateJadwal.UpdateJadwalProkerV3Command,
		*updateJadwal.UpdateJadwalProkerV3Result,
	](&updateJadwal.UpdateJadwalProkerV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		deleteJadwal.DeleteJadwalProkerV3Command,
		*deleteJadwal.DeleteJadwalProkerV3Result,
	](&deleteJadwal.DeleteJadwalProkerV3CommandHandler{
		Repo: repoProker,
	})

	// 8. Pagu Anggaran
	mediatr.RegisterRequestHandler[
		getAllPagu.GetAllPaguAnggaranQuery,
		*getAllPagu.GetAllPaguAnggaranResult,
	](&getAllPagu.GetAllPaguAnggaranQueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		savePagu.SavePaguAnggaranCommand,
		*savePagu.SavePaguAnggaranResult,
	](&savePagu.SavePaguAnggaranCommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		deletePagu.DeletePaguAnggaranCommand,
		*deletePagu.DeletePaguAnggaranResult,
	](&deletePagu.DeletePaguAnggaranCommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getUnitPagu.GetUnitPaguSummaryQuery,
		*domain.UnitPaguSummaryResult,
	](&getUnitPagu.GetUnitPaguSummaryQueryHandler{
		Repo: repoProker,
	})

	// 9. Dashboard Monitoring & Sankey
	mediatr.RegisterRequestHandler[
		getDashboard.GetDashboardProkerMonitoringQuery,
		*domain.DashboardMonitoringResult,
	](&getDashboard.GetDashboardProkerMonitoringQueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getSankey.GetSankeyDataQuery,
		*domain.SankeyDataResult,
	](&getSankey.GetSankeyDataQueryHandler{
		Repo: repoProker,
	})

	return nil
}
