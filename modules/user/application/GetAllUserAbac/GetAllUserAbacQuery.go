package GetAllUserAbac

import (
	domainuser "UnpakSiamida/modules/user/domain"
)

type GetAllUserAbacQuery struct {
	Tahun    string
	Status   string
	Role     string
	Search   string
	UserUuid string
	FuUuid   string
}

type GetAllUserAbacResult struct {
	Data []domainuser.UserAbac
}
