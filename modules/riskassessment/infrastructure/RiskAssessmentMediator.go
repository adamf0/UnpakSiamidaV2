package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
	create "UnpakSiamida/modules/riskassessment/application/CreateRiskAssessment"
	delete "UnpakSiamida/modules/riskassessment/application/DeleteRiskAssessment"
	getAll "UnpakSiamida/modules/riskassessment/application/GetAllRiskAssessments"
	get "UnpakSiamida/modules/riskassessment/application/GetRiskAssessment"
	getOptions "UnpakSiamida/modules/riskassessment/application/GetRiskOptions"
	update "UnpakSiamida/modules/riskassessment/application/UpdateRiskAssessment"
	domainRisk "UnpakSiamida/modules/riskassessment/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleRiskAssessment(db *gorm.DB) error {
	// if err := db.AutoMigrate(&domainRisk.RiskAssessment{}); err != nil {
	// 	return err
	// }

	repoRisk := NewRiskAssessmentRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)

	mediatr.RegisterRequestHandler[
		create.CreateRiskAssessmentCommand,
		string,
	](&create.CreateRiskAssessmentCommandHandler{
		Repo:             repoRisk,
		RepoFakultasUnit: repoFakultasUnit,
		DB:               db,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateRiskAssessmentCommand,
		string,
	](&update.UpdateRiskAssessmentCommandHandler{
		Repo:             repoRisk,
		RepoFakultasUnit: repoFakultasUnit,
		DB:               db,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteRiskAssessmentCommand,
		string,
	](&delete.DeleteRiskAssessmentCommandHandler{
		Repo: repoRisk,
	})

	mediatr.RegisterRequestHandler[
		get.GetRiskAssessmentByUuidQuery,
		*domainRisk.RiskAssessmentDefault,
	](&get.GetRiskAssessmentByUuidQueryHandler{
		Repo: repoRisk,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllRiskAssessmentsQuery,
		commondomain.Paged[domainRisk.RiskAssessmentDefault],
	](&getAll.GetAllRiskAssessmentsQueryHandler{
		Repo: repoRisk,
	})

	// Options Handlers
	mediatr.RegisterRequestHandler[
		getOptions.GetRiskOptionsQuery,
		*domainRisk.RiskOptionsResult,
	](&getOptions.GetRiskOptionsQueryHandler{
		Repo: repoRisk,
	})

	mediatr.RegisterRequestHandler[
		getOptions.GetRiskCriteriaQuery,
		[]domainRisk.RiskCriteriaItem,
	](&getOptions.GetRiskCriteriaQueryHandler{
		Repo: repoRisk,
	})

	mediatr.RegisterRequestHandler[
		getOptions.GetRiskLevelsQuery,
		[]domainRisk.RiskLevelItem,
	](&getOptions.GetRiskLevelsQueryHandler{
		Repo: repoRisk,
	})

	mediatr.RegisterRequestHandler[
		getOptions.GetRiskCategoriesQuery,
		[]domainRisk.RiskCategoryItem,
	](&getOptions.GetRiskCategoriesQueryHandler{
		Repo: repoRisk,
	})

	mediatr.RegisterRequestHandler[
		getOptions.GetRiskStatusesQuery,
		[]domainRisk.RiskStatusItem,
	](&getOptions.GetRiskStatusesQueryHandler{
		Repo: repoRisk,
	})

	mediatr.RegisterRequestHandler[
		getOptions.GetRiskProgramsQuery,
		[]domainRisk.RiskMataProgramItem,
	](&getOptions.GetRiskProgramsQueryHandler{
		Repo: repoRisk,
	})

	commoninfra.RegisterValidation(create.CreateRiskAssessmentCommandValidation, "RiskAssessmentCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateRiskAssessmentCommandValidation, "RiskAssessmentUpdate.Validation")
	commoninfra.RegisterValidation(delete.DeleteRiskAssessmentCommandValidation, "RiskAssessmentDelete.Validation")

	return nil
}
