package GetAllUserAbac

import (
	"context"

	domainuser "UnpakSiamida/modules/user/domain"
)

type GetAllUserAbacQueryHandler struct {
	Repo domainuser.IUserRepository
}

func (h *GetAllUserAbacQueryHandler) Handle(ctx context.Context, q GetAllUserAbacQuery) ([]domainuser.UserAbac, error) {
	return h.Repo.GetAllUserAbac(ctx, q.Tahun, q.Status, q.Role, q.Search, q.UserUuid, q.FuUuid)
}
