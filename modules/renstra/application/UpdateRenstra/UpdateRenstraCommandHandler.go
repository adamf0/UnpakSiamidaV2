package application

import (
	"context"
	
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"github.com/google/uuid"
)

type UpdateRenstraCommandHandler struct {
	Repo domainrenstra.IRenstraRepository
	FakultasUnitRepo domainfakultasunit.IFakultasUnitRepository
	UserRepo domainuser.IUserRepository
}

func (h *UpdateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateRenstraCommand,
) (string, error) {

	renstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainrenstra.InvalidUuid()
	}
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



	existingRenstra, err := h.Repo.GetByUuid(ctx, renstraUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		return "", err
	}
	if existingRenstra == nil {
		return "", domainrenstra.NotFound(cmd.Uuid)
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

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainrenstra.UpdateRenstra(
		existingRenstra,
		renstraUUID,
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
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedRenstra := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedRenstra); err != nil {
		return "", err
	}

	return updatedRenstra.UUID.String(), nil
}
