package application

import (
	common "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	domainkts "UnpakSiamida/modules/kts/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type UpdateKtsCommandHandler struct {
	Repo     domainkts.IKtsRepository
	RepoUser domainuser.IUserRepository
}

func (h *UpdateKtsCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateKtsCommand,
) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ktsUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainkts.InvalidUuid()
	}

	accUUID, err := uuid.Parse(cmd.Acc)
	if err != nil {
		return "", domainkts.InvalidUuid()
	}

	var (
		existingUser       *domainuser.User
		existingKts        *domainkts.Kts
		existingKtsDefault *domainkts.KtsDefault
	)

	g, ctxg := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		existingUser, err = h.RepoUser.GetByUuid(ctxg, accUUID)
		return err
	})

	g.Go(func() error {
		var err error
		existingKts, err = h.Repo.GetByUuid(ctxg, ktsUUID)
		return err
	})

	g.Go(func() error {
		var err error
		existingKtsDefault, err = h.Repo.GetDefaultByUuid(ctxg, ktsUUID)
		return err
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	if existingUser == nil {
		return "", domainkts.NotFound(cmd.Acc)
	}

	if existingKts == nil {
		return "", domainkts.NotFound(cmd.Uuid)
	}

	if existingKtsDefault == nil {
		return "", domainkts.NotFound(cmd.Uuid)
	}

	var result common.ResultValue[*domainkts.Kts]

	switch cmd.Step {
	case "step1":
		result = domainkts.UpdateKtsStep1(
			existingKts,
			existingKtsDefault,
			ktsUUID,

			*cmd.NomorLaporan,
			*cmd.TanggalLaporan,
			*cmd.UraianKetidaksesuaianP,
			*cmd.UraianKetidaksesuaianL,
			*cmd.UraianKetidaksesuaianO,
			*cmd.UraianKetidaksesuaianR,
			*cmd.AkarMasalah,
			*cmd.TindakanKoreksi,
			existingUser.ID,
			cmd.Tahun,
		)
	case "step2":
		statusAccAuditee, err := StringPtrToUint(cmd.StatusAccAuditee)
		if err != nil {
			return "", err
		}

		result = domainkts.UpdateKtsStep2(
			existingKts,
			existingKtsDefault,
			ktsUUID,

			statusAccAuditee,
			existingUser.ID,
			cmd.KeteranganTolak,
			cmd.TindakanPerbaikan,
			cmd.Tahun,
		)
	case "step2R":
		result = domainkts.UpdateKtsTindakan(
			existingKts,
			existingKtsDefault,
			ktsUUID,

			*cmd.TindakanPerbaikan,
			cmd.Tahun,
		)
	case "step3":
		result = domainkts.UpdateKtsStep3(
			existingKts,
			existingKtsDefault,
			ktsUUID,

			existingUser.ID,
			*cmd.TanggalPenyelesaian,
			cmd.Tahun,
		)
	case "step4":
		result = domainkts.UpdateKtsStep4(
			existingKts,
			existingKtsDefault,
			ktsUUID,

			*cmd.TinjauanTindakanPerbaikan,
			*cmd.TanggalClosing,
			existingUser.ID,
			cmd.Tahun,
		)
	case "step5":
		result = domainkts.UpdateKtsStep5(
			existingKts,
			existingKtsDefault,
			ktsUUID,

			*cmd.TanggalClosingFinal,
			*cmd.WmmUpmfUpmps,
			existingUser.ID,
			cmd.Tahun,
		)
	default:
		return "", domainkts.InvalidStep()
	}

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedKts := result.Value

	// 🔥 TRANSACTIONAL OUTBOX
	err = h.Repo.WithTx(ctx, func(txRepo domainkts.IKtsRepositoryTx) error {

		if err := txRepo.Update(ctx, updatedKts); err != nil {
			return err
		}

		for _, event := range updatedKts.DomainEvents() {
			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			outbox := commoninfra.OutboxMessage{
				ID:            event.ID(),
				Type:          reflect.TypeOf(event).String(),
				Payload:       string(payload),
				OccurredOnUTC: event.OccurredOnUTC(),
			}

			if err := txRepo.InsertOutbox(ctx, &outbox); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// 🧹 CLEAR AFTER COMMIT
	updatedKts.ClearDomainEvents()

	return updatedKts.UUID.String(), nil
}

func StringPtrToUint(ptr *string) (uint, error) {
	if ptr == nil || *ptr == "" {
		return 0, errors.New("value is empty")
	}

	v, err := strconv.ParseUint(*ptr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(v), nil
}
