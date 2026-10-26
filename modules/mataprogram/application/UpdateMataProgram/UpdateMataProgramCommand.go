package application

type UpdateMataProgramCommand struct {
	Uuid             string
	TahunUuid        string
	MataProgram      string
	FakultasUnitUuid string
	IdFakultasUnit   uint
}
