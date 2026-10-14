package application

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	commoninfra "UnpakSiamida/common/infrastructure"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainuser "UnpakSiamida/modules/user/domain"
	"time"

	"github.com/google/uuid"
)

type CreateUserCommandHandler struct {
	Repo             domainuser.IUserRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *CreateUserCommandHandler) Handle(
	ctx context.Context,
	cmd CreateUserCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var fu *int
	var targetNama *string
	var targetType *string

	if cmd.UuidFakultasUnit != nil && strings.TrimSpace(*cmd.UuidFakultasUnit) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(*cmd.UuidFakultasUnit))
		if err != nil {
			return "", domainuser.InvalidParsing("Fakultas Unit")
		}
		target, err := h.RepoFakultasUnit.GetDefaultByUuid(ctx, parsed)
		if err == nil && target != nil {
			fuId := int(target.ID)
			fu = &fuId
			targetNama = &target.Nama
			targetType = &target.Type
		}
	}

	result := domainuser.NewUser(
		cmd.Username,
		cmd.Password,
		cmd.Name,
		cmd.Email,
		cmd.Level,
		fu,
		targetNama,
		targetType,
	)

	//[pr] event dispatch

	if !result.IsSuccess {
		return "", result.Error
	}

	createUser := result.Value

	err := h.Repo.WithTx(ctx, func(txRepo domainuser.IUserRepositoryTx) error {

		if err := txRepo.Create(ctx, createUser); err != nil {
			return err
		}

		for _, event := range createUser.DomainEvents() {
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
	createUser.ClearDomainEvents()

	// if err := h.Repo.Create(ctx, createUser); err != nil {
	// 	return "", err
	// }

	return result.Value.UUID.String(), nil
}
