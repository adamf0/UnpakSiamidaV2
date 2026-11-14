package application

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	domainRenstra "UnpakSiamida/modules/renstra/domain"
	eventRenstra "UnpakSiamida/modules/renstra/event"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SendEmailRenstraCommandHandler struct {
	Repo domainRenstra.IRenstraRepository
}

func (h *SendEmailRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd SendEmailRenstraCommand,
) (*SendEmailRenstraResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if len(cmd.Uuids) == 0 {
		return nil, domainRenstra.EmptySelection()
	}

	rows, err := h.Repo.GetRenstraEmailTargets(ctx, cmd.Uuids)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, domainRenstra.NotFound(strings.Join(cmd.Uuids, ", "))
	}

	var queuedUUIDs []string
	var outboxIDs []string
	var skippedTargets []string
	var outboxMessages []commoninfra.OutboxMessage

	for _, row := range rows {
		kode := ""
		if row.KodeAkses != nil {
			kode = strings.TrimSpace(*row.KodeAkses)
		}
		if kode == "" {
			skippedTargets = append(skippedTargets, fmt.Sprintf("%s (%s)", row.TargetAudit, row.Tahun))
			continue
		}

		var recipients []eventRenstra.RenstraEmailRecipient
		addRec := func(role string, namePtr, emailPtr *string) {
			name := ""
			if namePtr != nil {
				name = strings.TrimSpace(*namePtr)
			}
			email := ""
			if emailPtr != nil {
				email = strings.TrimSpace(*emailPtr)
			}
			if name == "" && email == "" {
				return
			}
			recipients = append(recipients, eventRenstra.RenstraEmailRecipient{
				Role:  role,
				Name:  name,
				Email: email,
			})
		}

		addRec("Auditee", row.AuditeeName, row.AuditeeEmail)
		addRec("Auditor 1", row.Auditor1Name, row.Auditor1Mail)
		addRec("Auditor 2", row.Auditor2Name, row.Auditor2Mail)

		parsedRenstraUUID, err := uuid.Parse(row.UUID)
		if err != nil {
			continue
		}

		evt := eventRenstra.RenstraSendEmailEvent{
			EventID:     uuid.New(),
			RenstraUUID: parsedRenstraUUID,
			Tahun:       row.Tahun,
			TargetAudit: row.TargetAudit,
			KodeAkses:   kode,
			Recipients:  recipients,
			OccurredOn:  time.Now().UTC(),
		}

		payload, err := json.Marshal(evt)
		if err != nil {
			return nil, err
		}

		outbox := commoninfra.OutboxMessage{
			ID:            evt.ID(),
			Type:          commoninfra.CanonicalTypeName(evt),
			Payload:       string(payload),
			OccurredOnUTC: evt.OccurredOnUTC(),
		}

		outboxMessages = append(outboxMessages, outbox)
		queuedUUIDs = append(queuedUUIDs, row.UUID)
		outboxIDs = append(outboxIDs, outbox.ID)
	}

	if len(queuedUUIDs) == 0 {
		return nil, domainRenstra.MissingKodeAkses()
	}

	err = h.Repo.WithTx(ctx, func(txRepo domainRenstra.IRenstraRepositoryTx) error {
		for _, msg := range outboxMessages {
			if err := txRepo.InsertOutbox(ctx, &msg); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &SendEmailRenstraResult{
		Message:        fmt.Sprintf("Berhasil memasukkan %d pengiriman email ke antrean background job (Outbox).", len(queuedUUIDs)),
		QueuedCount:    len(queuedUUIDs),
		QueuedUuids:    queuedUUIDs,
		OutboxIds:      outboxIDs,
		SkippedCount:   len(skippedTargets),
		SkippedTargets: skippedTargets,
	}, nil
}
