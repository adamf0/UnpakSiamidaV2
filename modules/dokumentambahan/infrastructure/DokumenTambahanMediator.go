package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	delete "UnpakSiamida/modules/dokumentambahan/application/DeleteDokumenTambahan"
	getAll "UnpakSiamida/modules/dokumentambahan/application/GetAllDokumenTambahans"
	get "UnpakSiamida/modules/dokumentambahan/application/GetDokumenTambahan"
	setupUuid "UnpakSiamida/modules/dokumentambahan/application/SetupUuidDokumenTambahan"
	update "UnpakSiamida/modules/dokumentambahan/application/UpdateDokumenTambahan"
	domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"

	infraRenstra "UnpakSiamida/modules/renstra/infrastructure"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleDokumenTambahan(db *gorm.DB) error {
	repoDokumenTambahan := NewDokumenTambahanRepository(db)
	repoRenstra := infraRenstra.NewRenstraRepository(db)

	mediatr.RegisterRequestHandler[
		update.UpdateDokumenTambahanCommand,
		string,
	](&update.UpdateDokumenTambahanCommandHandler{
		Repo:        repoDokumenTambahan,
		RepoRenstra: repoRenstra,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteDokumenTambahanCommand,
		string,
	](&delete.DeleteDokumenTambahanCommandHandler{
		Repo: repoDokumenTambahan,
	})

	mediatr.RegisterRequestHandler[
		get.GetDokumenTambahanByUuidQuery,
		*domaindokumentambahan.DokumenTambahan,
	](&get.GetDokumenTambahanByUuidQueryHandler{
		Repo: repoDokumenTambahan,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllDokumenTambahansQuery,
		commondomain.Paged[domaindokumentambahan.DokumenTambahanDefault],
	](&getAll.GetAllDokumenTambahansQueryHandler{
		Repo: repoDokumenTambahan,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidDokumenTambahanCommand,
		string,
	](&setupUuid.SetupUuidDokumenTambahanCommandHandler{
		Repo: repoDokumenTambahan,
	})

	commoninfra.RegisterValidation(update.UpdateDokumenTambahanCommandValidation, "DokumenTambahanUpdate.Validation")
	commoninfra.RegisterValidation(delete.DeleteDokumenTambahanCommandValidation, "DokumenTambahanDelete.Validation")

	return nil
}
