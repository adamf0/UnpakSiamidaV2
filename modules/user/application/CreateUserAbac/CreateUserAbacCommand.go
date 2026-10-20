package CreateUserAbac

type CreateUserAbacCommand struct {
	IdUser           *uint  `json:"id_user" form:"id_user"`
	UserUuid         string `json:"user_uuid" form:"user_uuid"`
	Tahun            string `json:"tahun" form:"tahun"`
	IdFakultasUnit   *uint  `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	UuidFakultasUnit string `json:"uuid_fakultas_unit" form:"uuid_fakultas_unit"`
	Role             string `json:"role" form:"role"`
	Status           string `json:"status" form:"status"`
}

type CreateUserAbacResult struct {
	Message string `json:"message"`
	UUID    string `json:"uuid"`
	ID      uint   `json:"id"`
	Role    string `json:"role"`
}
