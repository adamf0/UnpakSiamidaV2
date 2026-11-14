package infrastructure

import (
	createAktivitas "UnpakSiamida/modules/aktivitasv2/application/CreateAktivitasV3"
	createTimeline "UnpakSiamida/modules/aktivitasv2/application/CreateTimelineV3"
	deleteAktivitas "UnpakSiamida/modules/aktivitasv2/application/DeleteAktivitasV3"
	deleteTimelineFile "UnpakSiamida/modules/aktivitasv2/application/DeleteTimelineFileV3"
	deleteTimeline "UnpakSiamida/modules/aktivitasv2/application/DeleteTimelineV3"
	getAktivitasByMatrix "UnpakSiamida/modules/aktivitasv2/application/GetAktivitasByMatrix"
	getAktivitasDetail "UnpakSiamida/modules/aktivitasv2/application/GetAktivitasDetail"
	reorderTimeline "UnpakSiamida/modules/aktivitasv2/application/ReorderTimelineV3"
	updateAktivitas "UnpakSiamida/modules/aktivitasv2/application/UpdateAktivitasV3"
	updateTimeline "UnpakSiamida/modules/aktivitasv2/application/UpdateTimelineV3"
	domain "UnpakSiamida/modules/aktivitasv2/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleAktivitasV2(db *gorm.DB) error {
	repoAktivitas := NewAktivitasV2Repository(db)

	mediatr.RegisterRequestHandler[
		getAktivitasByMatrix.GetAktivitasByMatrixQuery,
		*domain.AktivitasByMatrixResult,
	](&getAktivitasByMatrix.GetAktivitasByMatrixQueryHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		getAktivitasDetail.GetAktivitasDetailQuery,
		*domain.AktivitasDetailResult,
	](&getAktivitasDetail.GetAktivitasDetailQueryHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		createAktivitas.CreateAktivitasV3Command,
		*domain.AktivitasV2,
	](&createAktivitas.CreateAktivitasV3CommandHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		updateAktivitas.UpdateAktivitasV3Command,
		*domain.AktivitasV2,
	](&updateAktivitas.UpdateAktivitasV3CommandHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		deleteAktivitas.DeleteAktivitasV3Command,
		string,
	](&deleteAktivitas.DeleteAktivitasV3CommandHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		createTimeline.CreateTimelineV3Command,
		*domain.ProkerTimelineV3,
	](&createTimeline.CreateTimelineV3CommandHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		updateTimeline.UpdateTimelineV3Command,
		*domain.ProkerTimelineV3,
	](&updateTimeline.UpdateTimelineV3CommandHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		deleteTimeline.DeleteTimelineV3Command,
		string,
	](&deleteTimeline.DeleteTimelineV3CommandHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		reorderTimeline.ReorderTimelineV3Command,
		string,
	](&reorderTimeline.ReorderTimelineV3CommandHandler{
		Repo: repoAktivitas,
	})

	mediatr.RegisterRequestHandler[
		deleteTimelineFile.DeleteTimelineFileV3Command,
		string,
	](&deleteTimelineFile.DeleteTimelineFileV3CommandHandler{
		Repo: repoAktivitas,
	})

	return nil
}
