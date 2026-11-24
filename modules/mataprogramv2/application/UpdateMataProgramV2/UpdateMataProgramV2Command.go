package UpdateMataProgramV2

type UpdateMataProgramV2Command struct {
	UUID               string
	MataProgram        string
	Tahun              string
	TemplateRenstraIDs []uint
	HasTemplates       bool
}

type UpdateMataProgramV2Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UUID    string `json:"uuid"`
}
