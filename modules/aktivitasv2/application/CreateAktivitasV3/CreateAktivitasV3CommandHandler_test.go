package CreateAktivitasV3

import (
	domain "UnpakSiamida/modules/aktivitasv2/domain"
	"context"
	"testing"

	"github.com/google/uuid"
)

type stubCreateAktivitasRepo struct {
	domain.IAktivitasV2Repository
	proker  *domain.MatrixProkerSummary
	created *domain.AktivitasV2
}

func (s *stubCreateAktivitasRepo) GetMatrixProkerByUuid(ctx context.Context, uid uuid.UUID) (*domain.MatrixProkerSummary, error) {
	return s.proker, nil
}

func (s *stubCreateAktivitasRepo) CreateAktivitas(ctx context.Context, data *domain.AktivitasV2, timelines []domain.ProkerTimelineV3) error {
	s.created = data
	return nil
}

func TestCreateAktivitasV3_AlwaysStartsAsDraftEvenWhenProkerAccepted(t *testing.T) {
	for _, prokerStatus := range []string{"draft", "diterima", "diterima_parsial"} {
		t.Run("proker_"+prokerStatus, func(t *testing.T) {
			matUUID := uuid.New()
			repo := &stubCreateAktivitasRepo{
				proker: &domain.MatrixProkerSummary{
					ID:     31,
					UUID:   matUUID,
					Status: prokerStatus,
				},
			}
			handler := &CreateAktivitasV3CommandHandler{Repo: repo}

			act, err := handler.Handle(context.Background(), CreateAktivitasV3Command{
				MatrixUUID: matUUID.String(),
				Aktivitas:  "Aktivitas B2",
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if act.Status != "draft" || act.StatusVerifikator != "draft" || act.StatusEvaluator != "draft" {
				t.Fatalf("expected all initial statuses to be draft, got status=%q verif=%q eval=%q",
					act.Status, act.StatusVerifikator, act.StatusEvaluator)
			}
		})
	}
}
