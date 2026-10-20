package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	createAktivitas "UnpakSiamida/modules/prokerv3/application/CreateAktivitasV3"
	create "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"
	createTimeline "UnpakSiamida/modules/prokerv3/application/CreateTimelineV3"
	deleteAktivitas "UnpakSiamida/modules/prokerv3/application/DeleteAktivitasV3"
	deleteDokumen "UnpakSiamida/modules/prokerv3/application/DeleteDokumenV3"
	delete "UnpakSiamida/modules/prokerv3/application/DeleteProkerV3"
	deleteTimelineFile "UnpakSiamida/modules/prokerv3/application/DeleteTimelineFileV3"
	deleteTimeline "UnpakSiamida/modules/prokerv3/application/DeleteTimelineV3"
	getAktivitasByMatrix "UnpakSiamida/modules/prokerv3/application/GetAktivitasByMatrix"
	getAktivitasDetail "UnpakSiamida/modules/prokerv3/application/GetAktivitasDetail"
	getDokumen "UnpakSiamida/modules/prokerv3/application/GetDokumenV3"
	getHistory "UnpakSiamida/modules/prokerv3/application/GetHistoryStateV3"
	getAll "UnpakSiamida/modules/prokerv3/application/GetAllProkerV3"
	get "UnpakSiamida/modules/prokerv3/application/GetProkerV3"
	getOptions "UnpakSiamida/modules/prokerv3/application/GetProkerV3Options"
	getTargetUnpak "UnpakSiamida/modules/prokerv3/application/GetTargetUnpak"
	reorderTimeline "UnpakSiamida/modules/prokerv3/application/ReorderTimelineV3"
	updateAktivitas "UnpakSiamida/modules/prokerv3/application/UpdateAktivitasV3"
	update "UnpakSiamida/modules/prokerv3/application/UpdateProkerV3"
	updateStatus "UnpakSiamida/modules/prokerv3/application/UpdateStatusV3"
	updateTimeline "UnpakSiamida/modules/prokerv3/application/UpdateTimelineV3"
	updateVerifikasiDokumen "UnpakSiamida/modules/prokerv3/application/UpdateVerifikasiDokumenV3"
	domain "UnpakSiamida/modules/prokerv3/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleProkerV3(db *gorm.DB) error {
	repoProker := NewProkerV3Repository(db)

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
		getAktivitasByMatrix.GetAktivitasByMatrixQuery,
		*domain.AktivitasByMatrixResult,
	](&getAktivitasByMatrix.GetAktivitasByMatrixQueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getAktivitasDetail.GetAktivitasDetailQuery,
		*domain.AktivitasDetailResult,
	](&getAktivitasDetail.GetAktivitasDetailQueryHandler{
		Repo: repoProker,
		DB:   db,
	})

	mediatr.RegisterRequestHandler[
		createAktivitas.CreateAktivitasV3Command,
		*domain.AktivitasV2,
	](&createAktivitas.CreateAktivitasV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		updateAktivitas.UpdateAktivitasV3Command,
		*domain.AktivitasV2,
	](&updateAktivitas.UpdateAktivitasV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		deleteAktivitas.DeleteAktivitasV3Command,
		string,
	](&deleteAktivitas.DeleteAktivitasV3CommandHandler{
		Repo: repoProker,
	})

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

	mediatr.RegisterRequestHandler[
		createTimeline.CreateTimelineV3Command,
		*domain.ProkerTimelineV3,
	](&createTimeline.CreateTimelineV3CommandHandler{
		Repo: repoProker,
		DB:   db,
	})

	mediatr.RegisterRequestHandler[
		updateTimeline.UpdateTimelineV3Command,
		*domain.ProkerTimelineV3,
	](&updateTimeline.UpdateTimelineV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		deleteTimeline.DeleteTimelineV3Command,
		string,
	](&deleteTimeline.DeleteTimelineV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		reorderTimeline.ReorderTimelineV3Command,
		int,
	](&reorderTimeline.ReorderTimelineV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		deleteTimelineFile.DeleteTimelineFileV3Command,
		string,
	](&deleteTimelineFile.DeleteTimelineFileV3CommandHandler{
		Repo: repoProker,
	})

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

	return nil
}
