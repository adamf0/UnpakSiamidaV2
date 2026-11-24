package application

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	"UnpakSiamida/common/infrastructure/emailhub"
	domainuser "UnpakSiamida/modules/user/domain"
	eventUser "UnpakSiamida/modules/user/event"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SendAccountEmailCommandHandler struct {
	Repo domainuser.IUserRepository
}

func (h *SendAccountEmailCommandHandler) Handle(
	ctx context.Context,
	cmd SendAccountEmailCommand,
) (*SendAccountEmailResult, error) {
	if len(cmd.UserUUIDs) == 0 {
		return nil, domainuser.EmptyData()
	}

	var parsedUUIDs []uuid.UUID
	for _, raw := range cmd.UserUUIDs {
		if u, err := uuid.Parse(strings.TrimSpace(raw)); err == nil {
			parsedUUIDs = append(parsedUUIDs, u)
		}
	}

	if len(parsedUUIDs) == 0 {
		return nil, domainuser.InvalidData()
	}

	users, err := h.Repo.GetByUuids(ctx, parsedUUIDs)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, domainuser.EmptyData()
	}

	var recipients []eventUser.UserAccountEmailRecipient
	var queuedUUIDs []string

	for _, u := range users {
		email := strings.TrimSpace(u.Email)
		if email == "" {
			continue
		}

		// Kumpulkan daftar role unik dari penugasan ABAC
		var roleList []string
		roleSeen := make(map[string]bool)
		for _, a := range u.UserAbac {
			rName := strings.TrimSpace(a.Role)
			if rName == "" {
				rName = strings.TrimSpace(a.ComputedRole)
			}
			if rName != "" && !roleSeen[strings.ToLower(rName)] {
				roleSeen[strings.ToLower(rName)] = true
				if a.Tahun != "" {
					roleList = append(roleList, fmt.Sprintf("%s (%s)", rName, a.Tahun))
				} else {
					roleList = append(roleList, rName)
				}
			}
		}

		abacRoles := strings.Join(roleList, ", ")
		if abacRoles == "" {
			abacRoles = u.Level
		}
		if abacRoles == "" {
			abacRoles = "user"
		}

		recipients = append(recipients, eventUser.UserAccountEmailRecipient{
			UserUUID:  u.UUID,
			Username:  u.Username,
			Password:  u.Password,
			Name:      u.Name,
			Email:     email,
			Level:     u.Level,
			AbacRoles: abacRoles,
		})
		queuedUUIDs = append(queuedUUIDs, u.UUID.String())
	}

	if len(recipients) == 0 {
		return nil, fmt.Errorf("tidak ada pengguna terpilih yang memiliki alamat email valid")
	}

	eventID := uuid.New()
	evt := eventUser.UserSendAccountEmailEvent{
		EventID:     eventID,
		OccurredOn:  time.Now().UTC(),
		Users:       recipients,
		Attachments: cmd.Attachments,
	}

	payload, err := json.Marshal(evt)
	if err != nil {
		return nil, fmt.Errorf("gagal serialisasi outbox payload: %w", err)
	}

	outbox := commoninfra.OutboxMessage{
		ID:            eventID.String(),
		Type:          commoninfra.CanonicalTypeName(&evt),
		Payload:       string(payload),
		OccurredOnUTC: time.Now().UTC(),
	}

	now := time.Now().UTC()
	var histories []domainuser.HistoryEmail
	for _, rec := range recipients {
		histories = append(histories, domainuser.HistoryEmail{
			OutboxID:       outbox.ID,
			UserUUID:       rec.UserUUID.String(),
			Username:       rec.Username,
			RecipientEmail: rec.Email,
			Subject:        "Informasi Akun SIAMIDA V3 - Universitas Pakuan",
			Status:         "PENDING",
			Provider:       "GOOGLE_SMTP",
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}

	err = h.Repo.WithTx(ctx, func(txRepo domainuser.IUserRepositoryTx) error {
		if err := txRepo.InsertOutbox(ctx, &outbox); err != nil {
			return err
		}
		if len(histories) > 0 {
			if err := txRepo.InsertHistoryEmails(ctx, histories); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan outbox & history email: %w", err)
	}

	// Broadcast live status ke WebSocket & SSE bahwa ada email antrean baru
	for _, hItem := range histories {
		emailhub.GetHub().Broadcast(emailhub.EmailStatusEvent{
			Type:           "EMAIL_INSERTED",
			OutboxID:       hItem.OutboxID,
			UserUUID:       hItem.UserUUID,
			RecipientEmail: hItem.RecipientEmail,
			Username:       hItem.Username,
			Subject:        hItem.Subject,
			Status:         "PENDING",
			Provider:       hItem.Provider,
			UpdatedAt:      hItem.UpdatedAt,
		})
	}

	return &SendAccountEmailResult{
		Message:     fmt.Sprintf("Berhasil memasukkan %d pengiriman email akun ke antrean background job (Outbox).", len(queuedUUIDs)),
		TotalQueued: len(queuedUUIDs),
		UserUUIDs:   queuedUUIDs,
		OutboxID:    outbox.ID,
	}, nil
}
