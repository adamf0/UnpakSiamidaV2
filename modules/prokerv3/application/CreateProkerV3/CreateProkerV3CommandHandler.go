package CreateProkerV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateProkerV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func parseDateFlexible(val string) *time.Time {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", val); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return &t
	}
	return nil
}

func (h *CreateProkerV3CommandHandler) Handle(
	ctx context.Context,
	cmd CreateProkerV3Command,
) (string, error) {
	if err := cmd.Validate(); err != nil {
		return "", err
	}

	// 1. Resolve Fakultas / Unit ID
	var fuID *int = cmd.IdFakultasUnit
	if fuID == nil && strings.TrimSpace(cmd.FakultasUuid) != "" {
		trimmed := strings.TrimSpace(cmd.FakultasUuid)
		if id, err := strconv.Atoi(trimmed); err == nil && id > 0 {
			fuID = &id
		} else if idPtr, err := h.Repo.GetFakultasUnitIDByUUIDOrName(ctx, trimmed); err == nil && idPtr != nil {
			fuID = idPtr
		}
	}

	// 2. Resolve or Create Mata Program
	var mpID *int = cmd.IdMataProgram
	tahunStr := strings.TrimSpace(cmd.Tahun)
	if tahunStr == "" {
		tahunStr = "2026"
	}
	masterTahunID, _ := h.Repo.ResolveMasterTahun(ctx, tahunStr)

	if mpID == nil || *mpID <= 0 {
		namaProg := strings.TrimSpace(cmd.NamaProgram)
		if namaProg != "" {
			if createdID, err := h.Repo.ResolveOrCreateMataProgram(ctx, namaProg, masterTahunID, fuID); err == nil && createdID != nil {
				mpID = createdID
			}
		}
	} else if mpID != nil && *mpID > 0 && strings.TrimSpace(cmd.Tahun) != "" {
		_ = h.Repo.UpdateMataProgramTahun(ctx, *mpID, masterTahunID)
	}

	// 3. Resolve Indikator ID
	var indID *int = cmd.IdIndikator
	if indID == nil || *indID <= 0 {
		indKinerja := strings.TrimSpace(cmd.IndikatorKinerja)
		if indKinerja != "" {
			if resolvedInd, err := h.Repo.ResolveIndikatorID(ctx, indKinerja); err == nil && resolvedInd != nil {
				indID = resolvedInd
			}
		}
	}

	// 4. Create matrix_proker record
	matrixProker := &domain.ProkerV3{
		UUID:           uuid.New(),
		IdFakultasUnit: fuID,
		IdMataProgram:  mpID,
		IdIndikator:    indID,
		PIC:            strings.TrimSpace(cmd.PIC),
		BebanYayasan:   cmd.BebanYayasan,
	}

	if len(cmd.Scopes) > 0 {
		firstScope := cmd.Scopes[0]
		if firstScope.IdIndikator != nil && *firstScope.IdIndikator > 0 {
			matrixProker.IdIndikator = firstScope.IdIndikator
		}
	}

	if err := h.Repo.Create(ctx, matrixProker); err != nil {
		return "", err
	}

	// 5. Save dynamic scopes into matrix_proker_scope
	if len(cmd.Scopes) > 0 {
		scopesToSave := make([]domain.MatrixProkerScope, 0, len(cmd.Scopes))
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

			scopesToSave = append(scopesToSave, domain.MatrixProkerScope{
				UUID:             scUUID,
				IdMatrixProker:   matrixProker.ID,
				IdMataProgram:    mpID,
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
		_ = h.Repo.SaveMatrixProkerScopes(ctx, scopesToSave)
	} else {
		fallbackScope := []domain.MatrixProkerScope{
			{
				UUID:             uuid.New(),
				IdMatrixProker:   matrixProker.ID,
				IdMataProgram:    mpID,
				IdIndikator:      indID,
				TargetUnpak:      strings.TrimSpace(cmd.TargetUnpak),
				CapaianTerdahulu: strings.TrimSpace(cmd.CapaianUnitLalu),
				TargetUnit:       strings.TrimSpace(cmd.TargetUnit),
				Output:           strings.TrimSpace(cmd.Output),
				Outcome:          strings.TrimSpace(cmd.Outcome),
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}
		_ = h.Repo.SaveMatrixProkerScopes(ctx, fallbackScope)
	}

	return matrixProker.UUID.String(), nil
}
