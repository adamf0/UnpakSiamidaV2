package application

import (
	"context"

	domainberitaacara "UnpakSiamida/modules/beritaacara/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type UpdateBeritaAcaraCommandHandler struct {
	Repo             domainberitaacara.IBeritaAcaraRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
	RepoUser         domainuser.IUserRepository
}

func (h *UpdateBeritaAcaraCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateBeritaAcaraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	BeritaAcaraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainberitaacara.InvalidUuid()
	}

	// -------------------------
	// GET EXISTING BeritaAcara
	// -------------------------
	existingBeritaAcara, err := h.Repo.GetByUuid(ctx, BeritaAcaraUUID) // ← memastikan pakai nama interface yg benar
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", domainberitaacara.NotFound(cmd.Uuid)
		}
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------

	tanggal, err := time.Parse("2006-01-02", cmd.Tanggal)
	if err != nil {
		return "", domainberitaacara.InvalidTanggal()
	}

	uuidFakultasUnit, err := uuid.Parse(cmd.FakultasUnitUuid)
	if err != nil {
		return "", domainberitaacara.InvalidFakultasUnit()
	}

	uuidAuditee, err := uuid.Parse(cmd.AuditeeUuid)
	if err != nil {
		return "", domainberitaacara.InvalidAuditee()
	}

	var (
		idAuditor1 *uint
		idAuditor2 *uint
		fakultas   *domainfakultasunit.FakultasUnit
		auditee    *domainuser.User
	)

	g, ctxg := errgroup.WithContext(ctx)

	if cmd.Auditor1Uuid != nil {
		auditor1UUID := *cmd.Auditor1Uuid
		g.Go(func() error {
			uuidAuditor1, err := uuid.Parse(auditor1UUID)
			if err != nil {
				return nil // optional field → ignore invalid
			}
			auditor1, err := h.RepoUser.GetByUuid(ctxg, uuidAuditor1)
			if err != nil {
				return nil
			}
			idAuditor1 = &auditor1.ID
			return nil
		})
	}

	if cmd.Auditor2Uuid != nil {
		auditor2UUID := *cmd.Auditor2Uuid
		g.Go(func() error {
			uuidAuditor2, err := uuid.Parse(auditor2UUID)
			if err != nil {
				return nil
			}
			auditor2, err := h.RepoUser.GetByUuid(ctxg, uuidAuditor2)
			if err != nil {
				return nil
			}
			idAuditor2 = &auditor2.ID
			return nil
		})
	}

	g.Go(func() error {
		f, err := h.RepoFakultasUnit.GetDefaultByUuid(ctxg, uuidFakultasUnit)
		if err != nil {
			return domainberitaacara.NotFoundFakultas()
		}
		fakultas = f
		return nil
	})

	g.Go(func() error {
		u, err := h.RepoUser.GetByUuid(ctxg, uuidAuditee)
		if err != nil {
			return domainberitaacara.NotFoundAuditee()
		}
		auditee = u
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	result := domainberitaacara.UpdateBeritaAcara(
		existingBeritaAcara,
		BeritaAcaraUUID,
		cmd.Tahun,
		fakultas.ID,
		tanggal,
		auditee.ID,
		idAuditor1,
		idAuditor2,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedBeritaAcara := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedBeritaAcara); err != nil {
		return "", err
	}

	return updatedBeritaAcara.UUID.String(), nil
}
