package application

import (
	"context"
	
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"github.com/google/uuid"
	"time"
)

type CreateRenstraCommandHandler struct{
	Repo domainrenstra.IRenstraRepository
	FakultasUnitRepo domainfakultasunit.IFakultasUnitRepository
	UserRepo domainuser.IUserRepository
}

func (h *CreateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd CreateRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uuidFakultasUnit, err := uuid.Parse(cmd.FakultasUnit)
	if err != nil {
		return "", domainrenstra.InvalidParsing("Fakultas Unit")
	}
	uuidAuditee, err := uuid.Parse(cmd.Auditee)
	if err != nil {
		return "", domainrenstra.InvalidParsing("auditee")
	}
	uuidAuditor1, err := uuid.Parse(cmd.Auditor1)
	if err != nil {
		return "", domainrenstra.InvalidParsing("auditor1")
	}
	uuidAuditor2, err := uuid.Parse(cmd.Auditor2)
	if err != nil {
		return "", domainrenstra.InvalidParsing("auditor2")
	}


	fakultasUnit, err := h.FakultasUnitRepo.GetDefaultByUuid(ctx, uuidFakultasUnit)
	if err != nil {
		return "", domainrenstra.InvalidFakultasUnit()
	}
	auditee, err := h.UserRepo.GetByUuid(ctx, uuidAuditee)
	if err != nil {
		return "", domainrenstra.MissingAuditee()
	}
	auditor1, err := h.UserRepo.GetByUuid(ctx, uuidAuditor1)
	if err != nil {
		return "", domainrenstra.MissingAuditor1()
	}
	auditor2, err := h.UserRepo.GetByUuid(ctx, uuidAuditor2)
	if err != nil {
		return "", domainrenstra.MissingAuditor2()
	}


	isUnique, err := h.Repo.IsUnique(ctx, fakultasUnit.ID, cmd.Tahun)
	if err != nil {
		return "", err
	}

	result := domainrenstra.NewRenstra(
		cmd.Tahun,
		fakultasUnit.ID,
		cmd.PeriodeUploadMulai, 
		cmd.PeriodeUploadAkhir,
		cmd.PeriodeAssesmentDokumenMulai, 
		cmd.PeriodeAssesmentDokumenAkhir,
		cmd.PeriodeAssesmentLapanganMulai, 
		cmd.PeriodeAssesmentLapanganAkhir,
		auditee.ID, 
		auditor1.ID, 
		auditor2.ID,
		isUnique,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	renstra := result.Value

	if err := h.Repo.Create(ctx, renstra); err != nil {
		return "", err
	}

	return renstra.UUID.String(), nil
}
