package UpdateMataProgramV2

type UpdateMataProgramV2Command struct {
	UUID               string
	Scope              string // "group" (default) or "single"
	MataProgram        string
	Tahun              string
	Targets            []string
	TemplateRenstraIDs []uint
	IdFakultasUnit     string
}

type UpdateMataProgramV2Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UUID    string `json:"uuid"`
}
