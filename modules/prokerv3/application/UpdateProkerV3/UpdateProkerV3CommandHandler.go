package UpdateProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UpdateProkerV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateProkerV3CommandHandler) Handle(
	ctx context.Context,
	cmd UpdateProkerV3Command,
) (string, error) {
	if err := cmd.Validate(); err != nil {
		return "", err
	}

	prokerUUID, err := uuid.Parse(strings.TrimSpace(cmd.UUID))
	if err != nil {
		return "", domain.ProkerNotFound()
	}

	existing, err := h.Repo.GetByUuid(ctx, prokerUUID)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return "", domain.ProkerNotFound()
	}

	// 1. Resolve Fakultas Unit
	var fuID *int = cmd.IdFakultasUnit
	if fuID == nil && strings.TrimSpace(cmd.FakultasUuid) != "" {
		trimmed := strings.TrimSpace(cmd.FakultasUuid)
		if id, err := strconv.Atoi(trimmed); err == nil && id > 0 {
			fuID = &id
		} else if idPtr, err := h.Repo.GetFakultasUnitIDByUUIDOrName(ctx, trimmed); err == nil && idPtr != nil {
			fuID = idPtr
		}
	}
	if fuID != nil {
		existing.IdFakultasUnit = fuID
	}

	// 2. Resolve Mata Program & Update NamaProgram
	namaProg := strings.TrimSpace(cmd.NamaProgram)

	// If explicit id_mata_program provided and different from existing, switch to it
	if cmd.IdMataProgram != nil && *cmd.IdMataProgram > 0 {
		if existing.IdMataProgram == nil || *existing.IdMataProgram != *cmd.IdMataProgram {
			existing.IdMataProgram = cmd.IdMataProgram
		}
	}

	// Block duplicate pengajuan untuk kombinasi (id_mata_program / nama_program, fakultas_unit, tahun)
	if existing.IdFakultasUnit != nil && *existing.IdFakultasUnit > 0 {
		exists, err := h.Repo.ExistsActiveProker(ctx, *existing.IdFakultasUnit, strings.TrimSpace(cmd.Tahun), existing.IdMataProgram, namaProg, existing.ID)
		if err != nil {
			return "", err
		}
		if exists {
			return "", domain.DuplicateProkerUnitTahun()
		}
	}

	updatedMpID, err := h.Repo.SyncMataProgramForUpdate(ctx, existing.ID, existing.IdMataProgram, namaProg, strings.TrimSpace(cmd.Tahun), fuID)
	if err == nil && updatedMpID != nil {
		existing.IdMataProgram = updatedMpID
	}

	// 3. Resolve Indikator ID
	if cmd.IdIndikator != nil && *cmd.IdIndikator > 0 {
		existing.IdIndikator = cmd.IdIndikator
	} else if strings.TrimSpace(cmd.IndikatorKinerja) != "" {
		if resolvedInd, err := h.Repo.ResolveIndikatorID(ctx, strings.TrimSpace(cmd.IndikatorKinerja)); err == nil && resolvedInd != nil {
			existing.IdIndikator = resolvedInd
		}
	}

	// 4. Update PIC in matrix_proker & sync to aktivitasv2
	existing.PIC = strings.TrimSpace(cmd.PIC)
	_ = h.Repo.SyncAktivitasPIC(ctx, existing.ID, cmd.PIC)

	if cmd.BebanYayasan != nil {
		existing.BebanYayasan = *cmd.BebanYayasan
	}

	// 5. Save and sync dynamic scopes in matrix_proker_scope
	if len(cmd.Scopes) > 0 {
		firstScope := cmd.Scopes[0]
		if firstScope.IdIndikator != nil && *firstScope.IdIndikator > 0 {
			existing.IdIndikator = firstScope.IdIndikator
		}

		scopesToSync := make([]domain.MatrixProkerScope, 0, len(cmd.Scopes))
		for _, sc := range cmd.Scopes {
			var scIndID *int = sc.IdIndikator
			if scIndID == nil || *scIndID <= 0 {
				scIndKinerja := strings.TrimSpace(sc.IndikatorKinerja)
				if scIndKinerja != "" {
					if resolvedInd, err := h.Repo.ResolveIndikatorID(ctx, scIndKinerja); err == nil && resolvedInd != nil {
						scIndID = resolvedInd
					}
				}
			}

			scUUID := uuid.New()
			if strings.TrimSpace(sc.UUID) != "" {
				if parsed, err := uuid.Parse(strings.TrimSpace(sc.UUID)); err == nil {
					scUUID = parsed
				}
			}

			scopesToSync = append(scopesToSync, domain.MatrixProkerScope{
				UUID:             scUUID,
				IdMatrixProker:   existing.ID,
				IdMataProgram:    existing.IdMataProgram,
				IdIndikator:      scIndID,
				TargetUnpak:      strings.TrimSpace(sc.TargetUnpak),
				CapaianTerdahulu: strings.TrimSpace(sc.CapaianUnitLalu),
				TargetUnit:       strings.TrimSpace(sc.TargetUnit),
				Output:           strings.TrimSpace(sc.Output),
				Outcome:          strings.TrimSpace(sc.Outcome),
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			})
		}
		_ = h.Repo.SyncScopesForUpdate(ctx, existing.ID, scopesToSync, nil)
	} else {
		fallbackScope := &domain.MatrixProkerScope{
			UUID:             uuid.New(),
			IdMatrixProker:   existing.ID,
			IdMataProgram:    existing.IdMataProgram,
			IdIndikator:      existing.IdIndikator,
			TargetUnpak:      strings.TrimSpace(cmd.TargetUnpak),
			CapaianTerdahulu: strings.TrimSpace(cmd.CapaianUnitLalu),
			TargetUnit:       strings.TrimSpace(cmd.TargetUnit),
			Output:           strings.TrimSpace(cmd.Output),
			Outcome:          strings.TrimSpace(cmd.Outcome),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		_ = h.Repo.SyncScopesForUpdate(ctx, existing.ID, nil, fallbackScope)
	}

	if err := h.Repo.Update(ctx, existing); err != nil {
		return "", err
	}

	return existing.UUID.String(), nil
}
