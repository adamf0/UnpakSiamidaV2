package application

import (
	"context"
	"encoding/json"
	"reflect"

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

	fakultasUnitUUID := uuid.New()
	if cmd.UuidFakultasUnit != nil {
		parsed, err := uuid.Parse(*cmd.UuidFakultasUnit)
		if err != nil {
			return "", domainuser.InvalidParsing("Fakultas Unit")
		}
		fakultasUnitUUID = parsed
	}

	var fu *int
	fuId := 0
	target, err := h.RepoFakultasUnit.GetDefaultByUuid(ctx, fakultasUnitUUID)
	if target != nil {
		fuId = int(target.ID)
	}
	fu = &fuId

	// if cmd.FakultasUnit != nil {
	// 	raw := strings.TrimSpace(*cmd.FakultasUnit)

	// 	if raw != "" {
	// 		parsed, err := strconv.ParseInt(raw, 10, 64)
	// 		if err != nil {
	// 			return "", err
	// 		}

	// 		v := int(parsed)
	// 		fu = &v
	// 	}
	// }

	result := domainuser.NewUser(
		cmd.Username,
		cmd.Password,
		cmd.Name,
		cmd.Email,
		cmd.Level,
		fu,
		&target.Nama,
		&target.Type,
	)

	//[pr] event dispatch

	if !result.IsSuccess {
		return "", result.Error
	}

	createUser := result.Value

	err = h.Repo.WithTx(ctx, func(txRepo domainuser.IUserRepositoryTx) error {

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
