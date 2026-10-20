package UpdateUserAbac

type UpdateUserAbacCommand struct {
	UUID             string `json:"uuid" form:"uuid"`
	Tahun            string `json:"tahun" form:"tahun"`
	IdFakultasUnit   *uint  `json:"id_fakultas_unit" form:"id_fakultas_unit"`
	UuidFakultasUnit string `json:"uuid_fakultas_unit" form:"uuid_fakultas_unit"`
	Role             string `json:"role" form:"role"`
	Status           string `json:"status" form:"status"`
}
